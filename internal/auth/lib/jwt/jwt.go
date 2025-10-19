package jwt

import (
	"errors"
	"time"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrJWT = errors.New("failed to generate token")
)

func NewToken(user domain.User, app domain.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Sault))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
