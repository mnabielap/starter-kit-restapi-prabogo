package jwt

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPayload struct {
	Sub  string `json:"sub"` // User ID
	Type string `json:"type"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT token
func GenerateToken(userID string, expires time.Duration, tokenType string, secret string) (string, time.Time, error) {
	expirationTime := time.Now().Add(expires)

	// Verify UserID is valid UUID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", time.Time{}, err
	}

	claims := &TokenPayload{
		Sub:  uid.String(),
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expirationTime, nil
}

// Helper to generate Auth (Access + Refresh) tokens pair
func GenerateAuthTokens(userID string) (string, string, time.Time, time.Time, error) {
	secret := os.Getenv("JWT_SECRET")
	
	// Access Token
	accessMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRATION_MINUTES"))
	if accessMinutes == 0 { accessMinutes = 30 }
	accessTokenExpires := time.Duration(accessMinutes) * time.Minute
	
	accessToken, accessExp, err := GenerateToken(userID, accessTokenExpires, "access", secret)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	// Refresh Token
	refreshDays, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRATION_DAYS"))
	if refreshDays == 0 { refreshDays = 30 }
	refreshTokenExpires := time.Duration(refreshDays) * 24 * time.Hour
	
	refreshToken, refreshExp, err := GenerateToken(userID, refreshTokenExpires, "refresh", secret)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return accessToken, refreshToken, accessExp, refreshExp, nil
}

// ValidateLocalToken verifies a locally generated HMAC token
func ValidateLocalToken(tokenString string) (*TokenPayload, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenString, &TokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}