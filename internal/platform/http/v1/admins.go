package httpv1

import (
	"errors"
	"net/http"

	"github.com/Krokozabra213/sso/internal/platform/business"
	"github.com/Krokozabra213/sso/internal/platform/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) adminSignIn(c *gin.Context) {
	ctx := c.Request.Context()

	_, err := getUserIDCtx(c)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: "already authorized"})
		return
	}

	var inp domain.SignInInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: "invalid input body"})
		return
	}

	// проверить что пользователь существует и является админом и вернуть пару токенов при успехе
	res, err := h.busines.Auth.AdminSignIn(ctx, business.SignInInput{
		Username: inp.Username,
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, business.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: ErrInternal.Error()})

		return
	}

	setTokenCookies(c, res.AccessToken, res.RefreshToken)

	c.Status(http.StatusOK)
}

func (h *Handler) adminGetAllSchools(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := getUserIDCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response{Message: "error"})
		return
	}

	schools, err := h.busines.Admin.GetSchools(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: ErrInternal.Error()})
		return
	}

	response := dataResponse{
		UserData: userID,
		Data:     schools,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) adminCreateSchool(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := getUserIDCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response{Message: "error"})
		return
	}

	var inp domain.SchoolCreateInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: "invalid input body"})
		return
	}

	err = h.busines.Admin.CreateSchool(ctx, userID, inp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: ErrInternal.Error()})
		return
	}

	response := dataResponse{
		UserData: userID,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) adminUpdateSchool(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := getUserIDCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response{Message: "error"})
		return
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: EmptyIDParam})
	}

	schoolID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: InvalidIDParam})
	}

	var inp domain.SchoolUpdateInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: "invalid input body"})
		return
	}

	m := inp.ToMap()
	if len(m) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: "no fields to update"})
		return
	}

	err = h.busines.Admin.UpdateSchool(ctx, userID, schoolID, m)
	if err != nil {
		if errors.Is(err, business.ErrSchoolNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, response{Message: "school not found"})
			return
		}
		if errors.Is(err, business.ErrNoAccess) {
			c.AbortWithStatusJSON(http.StatusForbidden, response{Message: "no access to this school"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: ErrInternal.Error()})
		return
	}

	response := dataResponse{
		UserData: userID,
	}

	c.JSON(http.StatusNoContent, response)
}

func (h *Handler) adminDeleteSchool(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := getUserIDCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response{Message: "error"})
		return
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: EmptyIDParam})
	}

	schoolID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: InvalidIDParam})
	}

	err = h.busines.Admin.DeleteSchool(ctx, userID, schoolID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: ErrInternal.Error()})
		return
	}

	response := dataResponse{
		UserData: userID,
	}

	c.JSON(http.StatusNoContent, response)
}

func (h *Handler) adminGetAllCourses(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := getUserIDCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response{Message: "error"})
		return
	}

	courses, err := h.busines.Admin.GetAllCourses(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: ErrInternal.Error()})
		return
	}

	response := dataResponse{
		UserData: userID,
		Data:     courses,
	}

	c.JSON(http.StatusOK, response)
}
