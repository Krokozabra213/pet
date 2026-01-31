package httpv1

import (
	"errors"

	jwtv1 "github.com/Krokozabra213/sso/pkg/jwt-manager/v1"
	"github.com/gin-gonic/gin"
)

const userData = "user_data"

var (
	ErrTokenExpired = errors.New("err token expired")
)

func (h *Handler) softValidateJWTToken(c *gin.Context) {

	access := getFromCookie(c, "access_token")
	if access == "" {
		c.Next()
		return
	}

	accessData, err := h.validator.ValidateAccess(access)

	if err == nil {
		userDataCtx := UserShortData{
			UserID:   accessData.UserID,
			Username: accessData.Username,
		}
		c.Set(userData, userDataCtx)
		c.Next()
		return
	}

	if errors.Is(err, jwtv1.ErrValidExp) {
		refresh := getFromCookie(c, "refresh_token")
		if refresh == "" {
			c.Next()
			return
		}

		// в generatetokens обращаемся к sso сервису для генерации токнов
		accessToken, refreshToken, err := h.busines.Auth.RefreshTokens(refresh)
		if err != nil {
			// логируем ошибку
			c.Next()
			return
		}

		setTokenCookies(c, accessToken, refreshToken)

		accessData, err = h.validator.ValidateAccess(accessToken)
		if err != nil {
			c.Next()
			return
		}

		userDataCtx := UserShortData{
			UserID:   accessData.UserID,
			Username: accessData.Username,
		}
		c.Set(userData, userDataCtx)

		c.Next()
		return
	}

	c.Next()

}

func setTokenCookies(c *gin.Context, access, refresh string) {
	c.SetCookie(
		"access_token",
		access,
		15*60,
		"/",
		"",
		false,
		true,
	)

	c.SetCookie(
		"refresh_token",
		refresh,
		30*24*60*60,
		"/",
		"",
		false,
		true,
	)
}

func getFromCookie(c *gin.Context, nameCookie string) string {
	if cookie, err := c.Cookie(nameCookie); err == nil && cookie != "" {
		return cookie
	}
	return ""
}
