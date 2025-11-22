package auth

import (
	"time"

	"CLOB/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("CHANGE_ME_SUPER_SECRET") // ideal vir de env/config

type Claims struct {
	AccountID uint   `json:"account_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(account *models.Account) (string, error) {
	claims := &Claims{
		AccountID: account.ID,
		Email:     account.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24h
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
