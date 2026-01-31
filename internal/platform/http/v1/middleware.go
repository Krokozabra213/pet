package httpv1

import (
	"context"
	"errors"
	"net/http"

	jwtv1 "github.com/Krokozabra213/sso/pkg/jwt-manager/v1"
	"github.com/gin-gonic/gin"
)

const (
	userDataCtx     = "user_data"
	userIDCtx       = "user_id"
	accessTokenTTL  = 15 * 60           // 15 минут
	refreshTokenTTL = 30 * 24 * 60 * 60 // 30 дней
)

var (
	ErrTokenExpired        = errors.New("token expired")
	ErrRefreshEmptyToken   = errors.New("refresh token is empty")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrValidateToken       = errors.New("failed to validate token")
)

func (h *Handler) adminValidateJWTToken(c *gin.Context) {
	ctx := c.Request.Context()

	access := getFromCookie(c, "access_token")
	if access == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response{Message: ErrAuthorization})
		return
	}

	var userID int
	accessData, tokenErr := h.validator.ValidateAccess(access)

	if errors.Is(tokenErr, jwtv1.ErrValidExp) {
		refresh := getFromCookie(c, "refresh_token")
		var newAccess, newRefresh string
		var err error

		accessData, newAccess, newRefresh, err = h.regenerateAccessData(ctx, refresh)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response{Message: err.Error()})
			return
		}

		setTokenCookies(c, newAccess, newRefresh)

	} else if tokenErr != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response{Message: "invalid token"})
		return
	}

	userID = accessData.UserID
	err := h.busines.Auth.IsAdmin(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, response{Message: ErrAccessDenied})
		return
	}

	c.Set(userIDCtx, userID)
	c.Next()
}

func (h *Handler) softValidateJWTToken(c *gin.Context) {
	ctx := c.Request.Context()

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
		c.Set(userDataCtx, userDataCtx)
		c.Next()
		return
	}

	if errors.Is(err, jwtv1.ErrValidExp) {
		refresh := getFromCookie(c, "refresh_token")
		accessData, newAccess, newRefresh, err := h.regenerateAccessData(ctx, refresh)
		if err != nil {
			c.Next()
			return
		}

		userDataCtx := UserShortData{
			UserID:   accessData.UserID,
			Username: accessData.Username,
		}
		c.Set(userDataCtx, userDataCtx)
		setTokenCookies(c, newAccess, newRefresh)
		c.Next()
		return
	}

	c.Next()

}

func (h *Handler) regenerateAccessData(ctx context.Context, refresh string) (*jwtv1.AccessData, string, string, error) {
	if refresh == "" {
		return nil, "", "", ErrRefreshEmptyToken
	}
	accessToken, refreshToken, err := h.busines.Auth.RefreshTokens(ctx, refresh)
	if err != nil {
		return nil, "", "", ErrInvalidRefreshToken
	}
	data, err := h.validator.ValidateAccess(accessToken)
	if err != nil {
		return nil, "", "", ErrValidateToken
	}
	return data, accessToken, refreshToken, nil

}

func setTokenCookies(c *gin.Context, access, refresh string) {
	c.SetCookie(
		"access_token",
		access,
		accessTokenTTL,
		"/",
		"",
		false,
		true,
	)

	c.SetCookie(
		"refresh_token",
		refresh,
		refreshTokenTTL,
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
