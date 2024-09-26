package utils

import (
	// "encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTUtil struct {
	SecretKey string
}

func NewJWTUtil(secretKey string) *JWTUtil {
	return &JWTUtil{SecretKey: secretKey}
}

// Encode a payload into JWT
func (j *JWTUtil) Encode(payload map[string]interface{}, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range payload {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(expiration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// Decode a JWT and validate its signature and expiration
func (j *JWTUtil) Decode(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		payload := map[string]interface{}{}
		for k, v := range claims {
			payload[k] = v
		}
		return payload, nil
	}

	return nil, errors.New("unable to parse claims")
}
