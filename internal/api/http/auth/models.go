package auth

import "github.com/golang-jwt/jwt"

type Claims struct {
	jwt.StandardClaims

	UserID int `json:"user_id"`
}

type TokenManager interface {
	NewToken(userID int) (string, error)
	Parse(accessToken string) (int, error)
}
