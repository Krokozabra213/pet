package httpPlatform

import (
	"net/http"

	"github.com/Krokozabra213/sso/internal/platform/business"
	httpv1 "github.com/Krokozabra213/sso/internal/platform/http/v1"
	platformconfig "github.com/Krokozabra213/sso/newconfigs/platform"
	jwtv1 "github.com/Krokozabra213/sso/pkg/jwt-manager/v1"
	"github.com/Krokozabra213/sso/pkg/limiter"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	business       *business.Business
	tokenValidator *jwtv1.Validator
}

func NewHandler(business *business.Business, validator *jwtv1.Validator) *Handler {
	return &Handler{
		business:       business,
		tokenValidator: validator,
	}
}

func (h *Handler) Init(cfg *platformconfig.Config) *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL),
		corsMiddleware,
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := httpv1.NewHandler(h.business, h.tokenValidator)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
