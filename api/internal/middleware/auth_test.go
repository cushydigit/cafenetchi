package middleware

import (
	"cafenetchi-api/internal/utils"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	SECRET = "secret"
)

func TestAuth_ValidToken(t *testing.T) {

	token, _ := utils.GenerateJWT(
		123,
		"09123456789",
		"customer",
		SECRET,
		time.Hour,
	)

	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		id := UserID(r.Context())
		if id != 123 {
			t.Errorf("expected id %d, got %d", 123, id)
		}

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	Auth(SECRET)(next).ServeHTTP(rr, req)

	if !called {
		t.Errorf("expected called to be true")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

}

func TestAuth_MissingToken(t *testing.T) {

	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rr := httptest.NewRecorder()

	Auth(SECRET)(next).ServeHTTP(rr, req)

	if called {
		t.Errorf("expected called to be false")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuth_InvalidToken(t *testing.T) {

	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "invalid_token")

	rr := httptest.NewRecorder()

	Auth(SECRET)(next).ServeHTTP(rr, req)

	if called {
		t.Errorf("expected called to be false")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestUserID(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, int64(123))

	id := UserID(ctx)

	if id != 123 {
		t.Errorf("expected id %d, got %d", 123, id)
	}
}

func TestUserID_NotFound(t *testing.T) {
	id := UserID(context.Background())

	if id != 0 {
		t.Errorf("expected id %d, got %d", 0, id)
	}
}
