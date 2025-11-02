package jwt

import (
	"log"
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

type ResetClaims struct {
	Sub string `json:"sub"` // User.Id
	Exp int64  `json:"exp"`
}

const resetTTL = 15 * time.Minute

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

func (m *JwtManager) CreateResetToken(userID string) (string, error) {
	log.Println("jwt!")
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(resetTTL).Unix(),
		"prp": "reset",
	}
	log.Println(claims)

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println(tok.SignedString(m.SecretKey))
	return tok.SignedString(m.SecretKey)
}

func (m *JwtManager) VerifyResetToken(tokenStr string) (string, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return m.SecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil || !tok.Valid {
		return "", errs.ErrInvalidToken
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return "", errs.ErrInvalidToken
	}

	if claims["prp"] != "reset" {
		return "", errs.ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errs.ErrInvalidToken
	}

	exp, _ := claims["exp"].(float64)
	if time.Now().Unix() > int64(exp) {
		return "", errs.ErrInvalidToken
	}

	return sub, nil
}
