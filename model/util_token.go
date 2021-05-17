package model

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Ip string `json:"ip,omitempty"`
	jwt.StandardClaims
}

const (
	secretKey = "0AY9n_Qu021MgR7wum8E"
)

var (
	MapAccessToken  map[string]string = make(map[string]string)
	MapRefreshToken map[string]string = make(map[string]string)
)

func GenerateToken(ip string, expiredTime int64) (string, error) {
	claims := &Claims{
		Ip: ip,
		StandardClaims: jwt.StandardClaims{
			// ExpiresAt: expiredTime / 1000,
			ExpiresAt: 0, // remove expire time for token => QRCODE can use forever
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(token string) (*Claims, error) {

	claims := &Claims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		v := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return nil, errors.New(ErrTokenExpiredMsg)
		}
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, err
	}

	return claims, nil
}
