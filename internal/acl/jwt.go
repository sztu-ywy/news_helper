package acl

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/uozi-tech/cosy/settings"
	"time"
)

type JWTClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID uint64) (signedToken string, err error) {
	claims := JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * 24 * time.Hour).Unix(),
			Issuer:    "Store",
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = unsignedToken.SignedString([]byte(settings.AppSettings.JwtSecret))
	return
}

func ValidateJWT(token string) (claims *JWTClaims, err error) {
	if token == "" {
		err = errors.New("token is empty")
		return
	}
	unsignedToken, err := jwt.ParseWithClaims(
		token,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(settings.AppSettings.JwtSecret), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := unsignedToken.Claims.(*JWTClaims)
	if !ok {
		err = errors.New("convert to jwt claims error")
		return
	}
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		err = errors.New("jwt is expired")
	}
	return
}
