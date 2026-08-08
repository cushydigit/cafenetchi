package contextx

import (
	"context"
	"testing"
)

func TestRequestID(t *testing.T) {

	ctx := WithRequestID(context.Background(), "123")

	id := RequestID(ctx)

	if id != "123" {
		t.Errorf("expected id %s, got %s", "123", id)
	}

}

func TestRequestID_NotFound(t *testing.T) {
	id := RequestID(context.Background())
	if id != "" {
		t.Errorf("expected id %s, got %s", "", id)
	}
}
