package validation

import (
	"cafenetchi-api/internal/types"
	"strings"
)

func ValidateSendOTP(req types.SendOTPRequest) (phone string, err error) {

	phone, err = normalizeIranianPhone(phone)
	if err != nil {
		return "", err
	}
	return phone, nil
}

func ValidateVerifyOTP(req types.VerifyOTPRequest) (phone string, code string, err error) {

	phone, err = normalizeIranianPhone(req.Phone)
	if err != nil {
		return "", "", err
	}

	code, err = validateOTP(req.Code)
	if err != nil {
		return "", "", err
	}
	return phone, code, nil

}

func normalizeIranianPhone(phone string) (string, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return "", types.ErrPhoneRequired
	}

	switch {
	case strings.HasPrefix(phone, "+98"):
		phone = "0" + phone[3:]

	case strings.HasPrefix(phone, "98"):
		phone = "0" + phone[2:]
	}

	if !isIranianPhone(phone) {
		return "", types.ErrInvalidPhone
	}

	return phone, nil
}

func isIranianPhone(phone string) bool {
	if len(phone) != 11 {
		return false
	}

	if !strings.HasPrefix(phone, "09") {
		return false
	}

	for _, ch := range phone {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}

func validateOTP(code string) (string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", types.ErrOTPCodeRequired
	}

	if len(code) != 6 {
		return "", types.ErrInvalidOTP
	}

	for _, ch := range code {
		if ch < '0' || ch > '9' {
			return "", types.ErrInvalidOTP
		}
	}
	return code, nil
}
