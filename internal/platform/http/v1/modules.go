package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getModule(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("module_id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: EmptyIDParam})
	}

	moduleOutput, err := h.busines.Modules.GetByID(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{Message: err.Error()})
	}

	response := dataResponse{
		Data: moduleOutput,
	}

	userData, err := getUserDataCtx(c)

	if err == nil {
		response.UserData = userData
	}

	c.JSON(http.StatusOK, response)
}
