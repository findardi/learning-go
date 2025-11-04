package token

import (
	"fmt"
	"go-restapi/pkg/common"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(common.GetString("SECRET_KEY", "SUPER_SECRET"))

type Payload struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(id int32, username, role string) (string, error) {
	claims := Payload{
		UserID:   id,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", common.ErrTokenGenerate
	}

	return tokenString, nil
}

func VerifyToken(tokenStr string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Payload{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, common.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Payload)
	if !ok || !token.Valid {
		return nil, common.ErrTokenInvalid
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, common.ErrTokenExpired
	}

	return claims, nil
}
