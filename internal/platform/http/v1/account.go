package httpv1

import (
	"errors"
	"net/http"

	"github.com/Krokozabra213/sso/internal/platform/business"
	"github.com/Krokozabra213/sso/internal/platform/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getStudentProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userData, err := getUserDataCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response{Message: "error"})
	}

	student, err := h.busines.Auth.GetStudentByUserID(ctx, userData.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: "error"})
	}

	response := dataResponse{
		Data: student,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) signUp(c *gin.Context) {
	ctx := c.Request.Context()

	_, err := getUserDataCtx(c)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, response{Message: "already authorized"})
		return
	}

	var inp domain.StudentSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: "invalid input body"})
		return
	}

	if err := h.busines.Auth.SignUp(ctx, business.UserSignUpInput{
		Username: inp.Username,
		Email:    inp.Email,
		Password: inp.Password,
	}); err != nil {
		if errors.Is(err, business.ErrUserAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: err.Error()})

		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) signIn(c *gin.Context) {
	ctx := c.Request.Context()

	_, err := getUserDataCtx(c)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: "already authorized"})
		return
	}

	var inp domain.SignInInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: "invalid input body"})
		return
	}

	res, err := h.busines.Auth.SignIn(ctx, business.SignInInput{
		Username: inp.Username,
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, business.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{ErrInternal.Error()})

		return
	}

	setTokenCookies(c, res.AccessToken, res.RefreshToken)

	c.Status(http.StatusOK)
}

func (h *Handler) logout(c *gin.Context) {
	ctx := c.Request.Context()

	refresh := getFromCookie(c, "refresh_token")
	if refresh == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: "invalid token"})
		return
	}

	res, err := h.busines.Auth.Logout(ctx, refresh)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: err.Error()})
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.Status(http.StatusOK)
}
