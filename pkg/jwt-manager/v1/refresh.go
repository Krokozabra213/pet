package jwtv1

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (t *JWTManager) GenerateRefresh(
	refreshTTL time.Duration, userID uint64, username string, appID uint,
) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)

	err := t.refreshClaims(claims, refreshTTL, userID, username, appID)
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", ErrSignedRefreshToken
	}

	return tokenString, nil
}

func (gen *JWTManager) refreshClaims(
	claims jwt.MapClaims, refreshTTL time.Duration, userID uint64, username string, appID uint,
) error {

	jwtID, err := generateTokenID()
	if err != nil {
		return ErrGenerateJWTID
	}

	claims[JWTID] = jwtID
	claims[UserID] = userID
	claims[Username] = username
	claims[ExpiredAt] = time.Now().Add(refreshTTL).Unix()
	claims[AppID] = appID

	return nil
}

func generateTokenID() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	tokenID := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
	return tokenID, nil
}

// // Refresh JWT parser
// type RefreshData struct {
// 	JWTID    string
// 	UserID   int
// 	Username string
// 	Exp      time.Time
// 	AppID    int
// }

// func (data *RefreshData) Validate() error {
// 	if time.Now().After(data.Exp) {
// 		return ErrValidExp
// 	}
// 	return nil
// }

// func ParseRefresh(token string, keyManager IPublicKey) (*RefreshData, error) {

// 	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
// 		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
// 		}
// 		return keyManager.GetPublicKey(), nil
// 	}, jwt.WithoutClaimsValidation())

// 	if err != nil {
// 		return nil, ErrParseJWT
// 	}

// 	if !t.Valid {
// 		return nil, ErrValidToken
// 	}

// 	jwtData, err := refreshClaims(t.Claims.(jwt.MapClaims))
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = jwtData.Validate()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jwtData, nil
// }

// func refreshClaims(claims jwt.MapClaims) (*RefreshData, error) {
// 	// integers become float64 when decoding JWT
// 	userID, ok := claims[UserID].(float64)
// 	if !ok {
// 		return nil, ErrUserID
// 	}

// 	username, ok := claims[Username].(string)
// 	if !ok {
// 		return nil, ErrUsername
// 	}

// 	exp, ok := claims[ExpiredAt].(float64)
// 	if !ok {
// 		return nil, ErrExp
// 	}

// 	appID, ok := claims[AppID].(float64)
// 	if !ok {
// 		return nil, ErrAppID
// 	}

// 	jwtID, ok := claims[JWTID].(string)
// 	if !ok {
// 		return nil, ErrAppID
// 	}

// 	return &RefreshData{
// 		JWTID:    jwtID,
// 		UserID:   int(userID),
// 		Username: username,
// 		Exp:      JWTFloatToTime(exp),
// 		AppID:    int(appID),
// 	}, nil
// }

// func JWTFloatToTime(floatTime float64) time.Time {
// 	sec := int64(floatTime)
// 	nsec := int64((floatTime - float64(sec)) * 1e9)
// 	return time.Unix(sec, nsec).UTC()
// }
