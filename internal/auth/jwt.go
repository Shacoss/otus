package auth

import (
	"github.com/golang-jwt/jwt/v5"
	models "otus/pkg/model"
	"strconv"
)

func CreateJWT(user models.User, jwtSecret []byte) (string, error) {
	claims := &jwt.RegisteredClaims{
		ID:      strconv.FormatInt(user.ID, 10),
		Subject: user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
