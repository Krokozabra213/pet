package jwtv1

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (m *JWTManager) GenerateAccess(
	userID uint64, username string, appID int,
) (string, error) {

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)

	m.generateAccessClaims(claims, userID, username, appID)

	tokenString, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", ErrSignedAccessToken
	}

	return tokenString, nil
}

func (m *JWTManager) generateAccessClaims(
	claims jwt.MapClaims, userID uint64, username string, appID int,
) {
	claims[AppID] = appID
	claims[UserID] = userID
	claims[Username] = username
	claims[ExpiredAt] = time.Now().Add(m.accessTTL).Unix()
}

type AccessData struct {
	UserID   int
	Username string
	Exp      time.Time
	AppID    int
}

func (data *AccessData) Validate() error {
	if time.Now().After(data.Exp) {
		return ErrValidExp
	}
	return nil
}

func (m *JWTManager) ParseAccess(token string) (*AccessData, error) {

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
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

func accessClaims(claims jwt.MapClaims) (*AccessData, error) {
	// integers become float64 when decoding JWT
	userID, ok := claims[UserID].(float64)
	if !ok {
		return nil, ErrUserID
	}

	username, ok := claims[Username].(string)
	if !ok {
		return nil, ErrUsername
	}

	exp, ok := claims[ExpiredAt].(float64)
	if !ok {
		return nil, ErrExp
	}

	appID, ok := claims[AppID].(float64)
	if !ok {
		return nil, ErrAppID
	}

	return &AccessData{
		UserID:   int(userID),
		Username: username,
		Exp:      JWTFloatToTime(exp),
		AppID:    int(appID),
	}, nil
}
