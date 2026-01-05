package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) getModule(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("module_id")
	if idParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: EmptyIDParam})
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{Message: InvalidIDParam})
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
