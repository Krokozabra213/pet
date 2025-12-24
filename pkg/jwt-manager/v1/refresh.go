package jwtv1

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (m *JWTManager) GenerateRefresh(
	userID uint64, username string, appID int,
) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)

	err := m.generateRefreshClaims(claims, userID, username, appID)
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", ErrSignedRefreshToken
	}

	return tokenString, nil
}

func (m *JWTManager) generateRefreshClaims(
	claims jwt.MapClaims, userID uint64, username string, appID int,
) error {

	jwtID, err := generateTokenID()
	if err != nil {
		return ErrGenerateJWTID
	}

	claims[JWTID] = jwtID
	claims[UserID] = userID
	claims[Username] = username
	claims[ExpiredAt] = time.Now().Add(m.refreshTTL).Unix()
	claims[AppID] = appID

	return nil
}

type RefreshData struct {
	JWTID    string
	UserID   int
	Username string
	Exp      time.Time
	AppID    int
}

func (data *RefreshData) Validate() error {
	if time.Now().After(data.Exp) {
		return ErrValidExp
	}
	return nil
}

func (m *JWTManager) ParseRefresh(token string) (*RefreshData, error) {

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

	jwtData, err := refreshClaims(t.Claims.(jwt.MapClaims))
	if err != nil {
		return nil, err
	}

	err = jwtData.Validate()
	if err != nil {
		return nil, err
	}

	return jwtData, nil
}

func refreshClaims(claims jwt.MapClaims) (*RefreshData, error) {
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

	jwtID, ok := claims[JWTID].(string)
	if !ok {
		return nil, ErrJWTID
	}

	return &RefreshData{
		JWTID:    jwtID,
		UserID:   int(userID),
		Username: username,
		Exp:      JWTFloatToTime(exp),
		AppID:    int(appID),
	}, nil
}
