package app

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWT struct {
	key []byte
}

func NewJWT(key []byte) *JWT {
	return &JWT{key: key}
}

func (j *JWT) GenerateAccessToken(uuid *uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), jwt.MapClaims{
		"uuid":      uuid,
		"expiresAt": time.Now().Add(15 * time.Minute)})
	signedString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func (j *JWT) ParseUUID(signedString string) (*uuid.UUID, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(signedString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		return nil, err
	}
	userId, err := uuid.FromString(claims["uuid"].(string))
	if err != nil {
		return nil, err
	}
	return &userId, nil
}
