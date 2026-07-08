package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Phone  string `json:"phone"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JSON Web Token (JWT) with the provided claims.
//
// Parameters:
//   - userID: The ID of the user.
//   - phone: The phone number of the user.
//   - role: The role of the user.
//   - secretKey: The secret key used for signing the JWT.
//   - expiration: The duration of the JWT's validity.
//
// Returns:
//   - string: The generated JWT.
//   - error: An error if the JWT generation fails.
func GenerateJWT(userID int64, phone, role, secretKey string, expiration time.Duration) (string, error) {

	now := time.Now()

	claims := JWTClaims{
		UserID: userID,
		Phone:  phone,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			Issuer:    "cafenetchi-api",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

// ParseJWT parses a JSON Web Token (JWT) with the given tokenString and secretKey.
//
// Parameters:
//   - tokenString: The string representation of the JWT.
//   - secretKey: The secret key used for parsing the JWT.
//
// Returns:
//   - *JWTClaims: The parsed JWT claims.
//   - error: An error if the JWT parsing fails.
func ParseJWT(tokenString, secretKey string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodES256 {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil

}
