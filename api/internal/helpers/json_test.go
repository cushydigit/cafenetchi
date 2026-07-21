package helpers_test

import (
	"cafenetchi-api/internal/helpers"
	"cafenetchi-api/internal/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadJSON(t *testing.T) {
	type testCase struct {
		name string
		body string
		err  error
	}

	type userRequest struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	tests := []testCase{
		{
			name: "valid json",
			body: `{"name": "John Doe", "age": 30}`,
			err:  nil,
		},
		{
			name: "empty body",
			body: "",
			err:  types.ErrInvalidRequest,
		},
		{
			name: "invalid json",
			body: `{"name": "John Doe", "age": 30`,
			err:  types.ErrInvalidRequest,
		},
		{
			name: "multiple json values",
			body: `{}{}`,
			err:  types.ErrInvalidRequest,
		},
		{
			name: "non json",
			body: `hello word`,
			err:  types.ErrInvalidRequest,
		},
		{
			name: "trailing comma",
			body: `{"name": "John Doe", "age": 30,}`,
			err:  types.ErrInvalidRequest,
		},
		{
			name: "unmarshal type error",
			body: `{"name":"Ali","age":"twenty"}`,
			err:  types.ErrInvalidRequest,
		},
		{
			name: "too large body",
			body: `{"name":"` + strings.Repeat("a", 1024*1024+100) + `"}`,
			err:  types.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req userRequest
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			err := helpers.ReadJSON(w, r, &req)
			if tt.err == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != tt.err.Error() {
					t.Fatalf("expected %v, got %v", tt.err, err)
				}
			}
		})
	}
}

func decodeResponse[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()

	var v T

	if err := json.Unmarshal(rr.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}

	return v
}

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	err := helpers.WriteJSON(rr, http.StatusOK, types.Response{
		Data:    "hello",
		Message: "hi",
		Code:    http.StatusOK,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json got %q", got)
	}
	resp := decodeResponse[types.Response](t, rr)
	if resp.Data != "hello" {
		t.Fatalf("expected %s, got %s", "hello", resp.Data)
	}
}

func TestOK(t *testing.T) {
	rr := httptest.NewRecorder()
	helpers.OK(rr, "hello")
	resp := decodeResponse[types.Response](t, rr)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, resp.Code)
	}
	if resp.Data != "hello" {
		t.Fatalf("expected %s, got %s", "hello", resp.Data)
	}
}

func TestCreated(t *testing.T) {
	rr := httptest.NewRecorder()
	helpers.Created(rr, "hello")
	resp := decodeResponse[types.Response](t, rr)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, resp.Code)
	}
	if resp.Data != "hello" {
		t.Fatalf("expected %s, got %s", "hello", resp.Data)
	}
}

func TestError(t *testing.T) {

}
