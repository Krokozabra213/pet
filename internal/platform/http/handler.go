package http

import (
	"net/http"

	"github.com/Krokozabra213/sso/internal/platform/business"
	platformconfig "github.com/Krokozabra213/sso/newconfigs/platform"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	business *business.Business
}

func NewHandler(business *business.Business) *Handler {
	return &Handler{
		business: business,
	}
}

func Init(cfg *platformconfig.Config) *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router
}
