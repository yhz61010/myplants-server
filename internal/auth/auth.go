package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
    s := os.Getenv("JWT_SECRET")
    if s == "" {
        s = "dev-secret"
    }
    jwtSecret = []byte(s)
}

// SignToken creates a signed JWT for given user id and username
func SignToken(userID uint, username string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userId":   userID,
        "username": username,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
    })
    return token.SignedString(jwtSecret)
}

// ParseToken parses and validates a token string and returns claims
func ParseToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil || token == nil || !token.Valid {
        return nil, err
    }
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        return claims, nil
    }
    return nil, errors.New("invalid token claims")
}
