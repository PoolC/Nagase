package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var hmacSecret []byte

type AuthClaims struct {
	MemberUUID string `json:"member_uuid"`
	IsAdmin    bool   `json:"is_admin"`

	jwt.StandardClaims
}

func GenerateToken(memberUUID string, isAdmin bool) (string, error) {
	claim := AuthClaims{
		memberUUID,
		isAdmin,
		jwt.StandardClaims{
			Issuer:    "PoolC/Nagase",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(hmacSecret)
}

func ValidatedToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.MemberUUID, nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}

func init() {
	hmacSecret = []byte(os.Getenv("NAGASE_SECRET_KEY"))
}
