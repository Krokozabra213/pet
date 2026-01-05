package httpv1

import (
	"github.com/Krokozabra213/sso/internal/platform/business"
	jwtv1 "github.com/Krokozabra213/sso/pkg/jwt-manager/v1"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	busines   *business.Business
	validator *jwtv1.Validator
}

func NewHandler(business *business.Business, validator *jwtv1.Validator) *Handler {
	return &Handler{
		busines:   business,
		validator: validator,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initSchoolsRoutes(v1)
		h.initAccountRoutes(v1)
	}
}

func (h *Handler) initAccountRoutes(v1 *gin.RouterGroup) {
	account := v1.Group("account") // middleware валидирующая токен
	{
		account.POST("/sign-up")
		account.POST("/sign-in")
		account.POST("/auth/refresh") // refresh tokena
	}
}

func (h *Handler) initSchoolsRoutes(v1 *gin.RouterGroup) {
	schools := v1.Group("", h.softValidateJWTToken) // middleware валидирующая jwt token
	{
		schools.GET("", h.getAllPublishedSchools)
		schools.GET("/school/:school_id", h.getSchool)
		schools.GET("/course/:course_id", h.getCourse)
		schools.GET("/module/:module_id", h.getModule)
		schools.GET("/lesson/:lesson_id", h.getLesson)
	}
}
