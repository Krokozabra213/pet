package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getLesson(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("lesson_id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: EmptyIDParam})
	}

	lessonOutput, err := h.busines.Lessons.GetByID(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: err.Error()})
	}

	response := dataResponse{
		Data: lessonOutput,
	}

	userData, err := getUserDataCtx(c)

	if err == nil {
		response.UserData = userData
	}

	c.JSON(http.StatusOK, response)
}
