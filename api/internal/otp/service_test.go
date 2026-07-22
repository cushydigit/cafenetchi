package otp

import (
	"context"
	"errors"
	"testing"
)

func TestGenerate(t *testing.T) {
	store := newFakeStore()
	s := New(store)
	phone := "09123456789"

	code, err := s.Generate(context.Background(), phone)
	if err != nil {
		t.Fatal(err)
	}

	if len(code) != 6 {
		t.Fatalf("expected 6 digit, got %q", code)
	}
}

// generate store failure
func TestGenerateStoreError(t *testing.T) {
	store := newFakeStore()
	store.setErr = errors.New("redis down")

	s := New(store)

	_, err := s.Generate(context.Background(), "09123456789")

	if !errors.Is(err, ErrStoreSet) {
		t.Fatal("expected ErrStoreSet")
	}
}

// validate success
func TestValidateSuccess(t *testing.T) {
	store := newFakeStore()
	phone := "091223456789"
	code := "123456"
	store.values[phone] = code

	s := New(store)
	err := s.Validate(context.Background(), phone, code)

	if err != nil {
		t.Fatal(err)
	}
}

// validate failure
func TestValidateInvalid(t *testing.T) {
	store := newFakeStore()
	phone := "09123456789"
	code := "123456"
	store.values[phone] = code

	s := New(store)
	err := s.Validate(context.Background(), phone, "000000")
	if !errors.Is(err, ErrInvalid) {
		t.Fatal("expected ErrInvalid")
	}
}

// OTP expired
func TestValidateNotFound(t *testing.T) {

	s := New(newFakeStore())

	err := s.Validate(context.Background(), "09113456789", "123456")
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("expected ErrNotFound")
	}
}

// Consume OTP
func TestConsume(t *testing.T) {
	store := newFakeStore()
	phone := "09123456789"
	code := "123456"
	store.values[phone] = code

	s := New(store)
	err := s.Consume(context.Background(), phone)

	if err != nil {
		t.Fatal(err)
	}

	if _, ok := store.values[phone]; ok {
		t.Fatal("otp should be deleted")
	}
}
