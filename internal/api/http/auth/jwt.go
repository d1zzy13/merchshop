package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	signingKey string
	tokenTTL   time.Duration
}

func NewJWTManager(signingKey string, tokenTTL time.Duration) (*JWTManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &JWTManager{
		signingKey: signingKey,
		tokenTTL:   tokenTTL,
	}, nil
}

func (m *JWTManager) NewToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.signingKey))
}

func (m *JWTManager) Parse(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return claims.UserID, nil
}
