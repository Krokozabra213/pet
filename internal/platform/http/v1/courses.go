package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getCourse(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("course_id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: EmptyIDParam})
	}

	courseOutput, err := h.busines.Courses.GetByID(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: err.Error()})
	}

	modulesOutput, err := h.busines.Modules.GetByCourseID(ctx, courseOutput.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: err.Error()})
	}
	courseOutput.Modules = modulesOutput

	response := dataResponse{
		Data: courseOutput,
	}

	userData, err := getUserDataCtx(c)

	if err == nil {
		response.UserData = userData
	}

	c.JSON(http.StatusOK, response)
}
