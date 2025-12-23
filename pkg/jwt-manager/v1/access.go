package jwtv1

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (t *JWTManager) GenerateAccess(
	accessTTL time.Duration, userID uint64, username string, appID uint,
) (string, error) {

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)

	t.accessClaims(claims, accessTTL, userID, username, appID)

	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return "", ErrSignedAccessToken
	}

	return tokenString, nil
}

func (t *JWTManager) accessClaims(
	claims jwt.MapClaims, accessTTL time.Duration, userID uint64, username string, appID uint,
) {
	claims[AppID] = appID
	claims[UserID] = userID
	claims[Username] = username
	claims[ExpiredAt] = time.Now().Add(accessTTL).Unix()
}

// type AccessData struct {
// 	UserID   int
// 	Username string
// 	Exp      time.Time
// 	AppID    int
// }

// func (data *AccessData) Validate() error {
// 	if time.Now().After(data.Exp) {
// 		return ErrValidExp
// 	}
// 	return nil
// }

// func ParseAccess(token string, keyManager IPublicKey) (*AccessData, error) {

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

// 	jwtData, err := accessClaims(t.Claims.(jwt.MapClaims))
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = jwtData.Validate()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jwtData, nil
// }

// func accessClaims(claims jwt.MapClaims) (*AccessData, error) {
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

// 	return &AccessData{
// 		UserID:   int(userID),
// 		Username: username,
// 		Exp:      JWTFloatToTime(exp),
// 		AppID:    int(appID),
// 	}, nil
// }
