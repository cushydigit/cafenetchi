package validation

import (
	"cafenetchi-api/internal/types"
	"errors"
	"testing"
)

func TestNormalizeIranianPhone(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:  "valid local",
			input: "09123456789",
			want:  "09123456789",
		},
		{
			name:  "with +98",
			input: "+989123456789",
			want:  "09123456789",
		},
		{
			name:  "with 98",
			input: "989123456789",
			want:  "09123456789",
		},
		{
			name:  "trim spaces",
			input: " 09123456789 ",
			want:  "09123456789",
		},
		{
			name:    "empty",
			input:   "",
			wantErr: types.ErrPhoneRequired,
		},
		{
			name:    "too short",
			input:   "09123",
			wantErr: types.ErrInvalidPhone,
		},
		{
			name:    "contains letters",
			input:   "09123abc789",
			wantErr: types.ErrInvalidPhone,
		},
		{
			name:    "wrong prefix",
			input:   "08123456789",
			wantErr: types.ErrInvalidPhone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeIranianPhone(tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Fatalf("expected %q got %q", tt.want, got)
			}
		})
	}
}

func TestValidateOTP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:  "valid",
			input: "123456",
			want:  "123456",
		},
		{
			name:    "empty",
			input:   "",
			wantErr: types.ErrOTPCodeRequired,
		},
		{
			name:    "too short",
			input:   "12345",
			wantErr: types.ErrInvalidOTP,
		},
		{
			name:    "too long",
			input:   "1234567",
			wantErr: types.ErrInvalidOTP,
		},
		{
			name:    "contains letters",
			input:   "12ab56",
			wantErr: types.ErrInvalidOTP,
		},
		{
			name:    "contains spaces",
			input:   "12 456",
			wantErr: types.ErrInvalidOTP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateOTP(tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if got != tt.want {
				t.Fatalf("expected %q got %q", tt.want, got)
			}
		})
	}
}

func TestValidateSendOTP(t *testing.T) {
	req := types.SendOTPRequest{
		Phone: "+989123456789",
	}

	phone, err := ValidateSendOTP(req)
	if err != nil {
		t.Fatal(err)
	}

	if phone != "09123456789" {
		t.Fatalf("expected normalized phone, got %q", phone)
	}
}

func TestValidateVerifyOTP(t *testing.T) {
	req := types.VerifyOTPRequest{
		Phone: "+989123456789",
		Code:  "123456",
	}

	phone, code, err := ValidateVerifyOTP(req)
	if err != nil {
		t.Fatal(err)
	}

	if phone != "09123456789" {
		t.Fatalf("unexpected phone %q", phone)
	}

	if code != "123456" {
		t.Fatalf("unexpected code %q", code)
	}
}
