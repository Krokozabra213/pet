package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	EmptyIDParam = "empty id param"
)

type dataResponse struct {
	Data     interface{} `json:"data"`
	UserData interface{} `json:"userData"`
}

type response struct {
	Message string `json:"message"`
}

func (h *Handler) getAllPublishedSchools(c *gin.Context) {
	ctx := c.Request.Context()

	schoolsAllOutput, err := h.busines.Schools.GetAllPublished(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: err.Error()})
	}

	response := dataResponse{
		Data: schoolsAllOutput,
	}

	userData, err := getUserDataCtx(c)

	if err == nil {
		response.UserData = userData
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) getSchool(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("school_id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: EmptyIDParam})
	}

	schoolOutput, err := h.busines.Schools.GetByID(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: err.Error()})
	}

	response := dataResponse{
		Data: schoolOutput,
	}

	userData, err := getUserDataCtx(c)

	if err == nil {
		response.UserData = userData
	}

	c.JSON(http.StatusOK, response)
}
