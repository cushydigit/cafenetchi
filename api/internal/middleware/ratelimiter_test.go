package middleware

import (
	"cafenetchi-api/internal/types"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeLimiter struct {
	allow bool
	retry time.Duration
	err   error
}

func (f *fakeLimiter) Allow(ctx context.Context, key string) (bool, time.Duration, error) {
	return f.allow, f.retry, f.err
}

func TestRateLimit_allowed(t *testing.T) {
	limiter := &fakeLimiter{allow: true}
	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rr := httptest.NewRecorder()

	RateLimit(limiter)(next).ServeHTTP(rr, req)

	if !called {
		t.Errorf("expected called to be true")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRateLimit_Blocked(t *testing.T) {
	limiter := &fakeLimiter{
		allow: false,
		retry: 30 * time.Second,
	}
	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rr := httptest.NewRecorder()

	RateLimit(limiter)(next).ServeHTTP(rr, req)

	if called {
		t.Errorf("expected called to be false")
	}

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("expected status code %d, got %d", http.StatusTooManyRequests, rr.Code)
	}

	if rr.Header().Get("Retry-After") != "30" {
		t.Errorf("expected Retry-After header to be 30, got %s", rr.Header().Get("Retry-After"))
	}
}

func TestRateLimit_Error(t *testing.T) {
	limiter := &fakeLimiter{err: types.ErrInternalServer}
	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rr := httptest.NewRecorder()

	RateLimit(limiter)(next).ServeHTTP(rr, req)

	if called {
		t.Errorf("expected called to be false")
	}

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}
