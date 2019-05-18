package nagase

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
)

var hmacSecret []byte

type Member struct {
	UUID string `gorm:"type:varchar(40);PRIMARY_KEY"`

	LoginID      string `gorm:"type:varchar(40);UNIQUE_INDEX"`
	PasswordHash []byte `gorm:"NOT NULL" json:"-"`
	PasswordSalt []byte `gorm:"NOT NULL" json:"-"`
	Email        string `gorm:"type:varchar(255);UNIQUE_INDEX"`
	PhoneNumber  string `gorm:"type:varchar(20)"`
	Name         string `gorm:"type:varchar(40)"`
	Department   string `gorm:"type:varchar(40)"`
	StudentID    string `gorm:"type:varchar(40);UNIQUE_INDEX"`

	IsActivated bool `gorm:"default:false"`
	IsAdmin     bool `gorm:"default:false"`

	PasswordResetToken           string `gorm:"type:varchar(255)"`
	PasswordResetTokenValidUntil time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type memberClaims struct {
	MemberUUID string `json:"member_uuid"`
	IsAdmin    bool   `json:"is_admin"`

	jwt.StandardClaims
}

func (member *Member) ValidatePassword(password string) bool {
	hash := argon2.IDKey([]byte(password), member.PasswordSalt, 1, 8*1024, 4, 32)
	return bytes.Compare(hash, member.PasswordHash) == 0
}

func (member *Member) GenerateAccessToken() (string, error) {
	claim := memberClaims{
		member.UUID,
		member.IsAdmin,
		jwt.StandardClaims{
			Issuer:    "PoolC/Nagase",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(hmacSecret)
}

func GetMemberUUIDFromToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &memberClaims{}, func(t *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*memberClaims); ok && token.Valid {
		return claims.MemberUUID, nil
	}
	return "", fmt.Errorf("invalid token")
}

func init() {
	hmacSecret = []byte(os.Getenv("NAGASE_SECRET_KEY"))
}
