package jwtv1

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Validator struct {
	publicKey *rsa.PublicKey
}

func NewValidator(publicKey *rsa.PublicKey) (*Validator, error) {
	if publicKey == nil {
		return nil, ErrEmptyPublicKey
	}

	return &Validator{
		publicKey: publicKey,
	}, nil
}

func (v *Validator) ValidateAccess(tokenString string) (*AccessData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return v.publicKey, nil
	})

	if err != nil {
		return nil, ErrParseJWT
	}

	if !token.Valid {
		return nil, ErrValidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	jwtData, err := accessClaims(claims)
	if err != nil {
		return nil, err
	}

	err = jwtData.Validate()
	if err != nil {
		return nil, err
	}

	return jwtData, nil
}

func (m *JWTManager) ValidateRefresh(tokenString string) (*RefreshData, error) {

	t, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.publicKey, nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, ErrParseJWT
	}

	if !t.Valid {
		return nil, ErrValidToken
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	jwtData, err := refreshClaims(claims)
	if err != nil {
		return nil, err
	}

	err = jwtData.Validate()
	if err != nil {
		return nil, err
	}

	return jwtData, nil
}
