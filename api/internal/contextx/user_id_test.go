package contextx

import (
	"context"
	"testing"
)

func TestUserID(t *testing.T) {
	ctx := WithUserID(context.Background(), int64(123))

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
