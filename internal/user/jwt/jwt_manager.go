package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/danilobml/user-manager/internal/errs"
	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	SecretKey []byte
}

func NewJwtManager(secretKey []byte) *JwtManager {
	return &JwtManager{
		SecretKey: secretKey,
	}
}

func (j *JwtManager) CreateToken(email string, roles []string) (string, error) {
	rolesString := fmt.Sprint(strings.Join(roles, ","))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"roles": rolesString,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(j.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtManager) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return j.SecretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errs.ErrInvalidToken
	}

	return nil
}
