package repository

import (
	"cafenetchi-api/internal/model"
	"context"
	"errors"
	"testing"
)

var (
	PHONE = "09123456789"
)

func TestUserRepository_Create(t *testing.T) {
	repo := setupUserRepo(t)

	ctx := context.Background()

	user, err := repo.User.Create(ctx, PHONE)
	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user is nil")
	}

	if user.Phone != PHONE {
		t.Fatalf("unexpected phone %s", user.Phone)
	}

	if user.ID == 0 {
		t.Fatal("unexpected generated id")
	}
}

func TestDuplicateUser(t *testing.T) {
	repo := setupUserRepo(t)

	created, err := repo.User.Create(t.Context(), PHONE)
	if err != nil {
		t.Fatal(err)
	}
	if created == nil {
		t.Fatal(err)
	}
	duplicated, err := repo.User.Create(t.Context(), PHONE)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if duplicated != nil {
		t.Fatalf("expected nil, got %v", duplicated)
	}
	if !errors.Is(err, ErrPhoneExists) {
		t.Fatalf("expected %v, got %v", ErrPhoneExists, err)
	}
}

func TestGetUser(t *testing.T) {
	repo := setupUserRepo(t)

	created, err := repo.User.Create(repo.ctx, PHONE)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.User.Get(repo.ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != created.ID {
		t.Fatalf("expected id %d, got %d", created.ID, got.ID)
	}
}

func TestNotFoundUser(t *testing.T) {
	repo := setupUserRepo(t)

	got, err := repo.User.Get(t.Context(), 1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected %v, got %v", ErrUserNotFound, err)
	}
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestGetUserByPhone(t *testing.T) {
	repo := setupUserRepo(t)

	created, err := repo.User.Create(t.Context(), PHONE)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.User.GetByPhone(t.Context(), created.Phone)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != created.ID {
		t.Fatalf("expected id %d, got %d", created.ID, got.ID)
	}
}

func TestNotFoundUserByPhone(t *testing.T) {
	repo := setupUserRepo(t)

	got, err := repo.User.GetByPhone(t.Context(), "09111111111")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected %v, got %v", ErrUserNotFound, err)
	}
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestUpdateUser(t *testing.T) {
	repo := setupUserRepo(t)

	created, err := repo.User.Create(t.Context(), PHONE)
	if err != nil {
		t.Fatal(err)
	}

	req := model.UpdateUser{
		FullName:  "new name",
		AvatarURL: "new avatar url",
	}

	updated, err := repo.User.Update(t.Context(), created.ID, req)
	if err != nil {
		t.Fatal(err)
	}

	if updated.FullName != req.FullName {
		t.Fatalf("expected %s, got %s", req.FullName, updated.FullName)
	}
	if updated.AvatarURL != req.AvatarURL {
		t.Fatalf("expected %s, got %s", req.AvatarURL, updated.AvatarURL)
	}
}
