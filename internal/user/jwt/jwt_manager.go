package jwt

import (
	"time"

	"github.com/danilobml/user-manager/internal/errs"
	"github.com/danilobml/user-manager/internal/user/model"
	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	SecretKey []byte
}

type Claims struct {
	Email string
	Roles []model.Role
	jwt.RegisteredClaims
}

func NewJwtManager(secretKey []byte) *JwtManager {
	return &JwtManager{
		SecretKey: secretKey,
	}
}

func (j *JwtManager) CreateToken(email string, roles []model.Role) (string, error) {
	claims := Claims{
		Email: email,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.SecretKey)
}

func (j *JwtManager) ParseAndValidateToken(tokenString string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return j.SecretKey, nil
	})
	if err != nil {
		return nil, errs.ErrParsingToken
	}

	if !token.Valid {
		return nil, errs.ErrInvalidToken
	}

	return token.Claims.(*Claims), nil
}
