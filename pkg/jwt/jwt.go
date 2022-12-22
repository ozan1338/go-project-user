package jwt

import (
	"errors"
	"net/http"
	resError "project/util/errors_response"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKey = 5

var (
	JWTError = "jwt error"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker) {

	return &JWTMaker{secretKey: secretKey}
}

func (j *JWTMaker) CreateToken(id int, duration time.Duration) (string, *Payload, resError.RespError) {
	payload := NewPayload(id , duration)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(j.secretKey))

	if err != nil {
		return "", nil, resError.NewRespError(err.Error(), http.StatusUnauthorized, JWTError)
	}

	return token, payload, nil
}

func (j *JWTMaker) VerifyToken(token string) (*Payload, resError.RespError) {
	keyFunc := func(token *jwt.Token) (any,error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpired) {
			return nil, resError.NewRespError(ErrExpired.Error(), http.StatusUnauthorized, JWTError)
		}

		return nil, resError.NewRespError(ErrInvalidToken.Error(), http.StatusUnauthorized, JWTError)
	}

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, resError.NewRespError(ErrInvalidToken.Error(), http.StatusUnauthorized, JWTError)
	}

	return payload, nil
}