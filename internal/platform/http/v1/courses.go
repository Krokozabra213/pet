package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) getCourse(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("course_id")
	if idParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: EmptyIDParam})
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: InvalidIDParam})
	}

	courseOutput, err := h.busines.Courses.GetByID(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: err.Error()})
	}

	response := dataResponse{
		Data: courseOutput,
	}

	userData, err := getUserDataCtx(c)

	if err == nil {
		response.UserData = userData
	}

	c.JSON(http.StatusOK, response)
}
