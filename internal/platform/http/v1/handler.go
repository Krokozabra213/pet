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

func (h *Handler) initAdminsRoutes(v1 *gin.RouterGroup) {
	admins := v1.Group("/admins", h.adminValidateJWTToken)
	{
		admins.POST("/sign-in", h.adminSignIn)

		schools := admins.Group("/schools")
		{
			schools.GET("", h.adminGetAllSchools)
			schools.POST("", h.adminCreateSchool)
			schools.PATCH("/:id", h.adminUpdateSchool)
			schools.DELETE("/:id", h.adminDeleteSchool)
		}

		courses := admins.Group("/courses")
		{
			courses.GET("", h.adminGetAllCourses)
			courses.POST("", h.adminCreateCourse--)
			schools.PATCH("/:id", h.adminUpdateSchool)
			schools.DELETE("/:id", h.adminDeleteSchool)
		}
	}

}

func (h *Handler) initAccountRoutes(v1 *gin.RouterGroup) {
	account := v1.Group("account")
	{
		account.GET("/profile", h.getStudentProfile)
		account.POST("/sign-up", h.signUp)
		account.POST("/sign-in", h.signIn)
	}
}

func (h *Handler) initSchoolsRoutes(v1 *gin.RouterGroup) {
	schools := v1.Group("", h.softValidateJWTToken)
	{
		schools.GET("", h.getAllPublishedSchools)
		schools.GET("/school/:school_id", h.getSchool)
		schools.GET("/course/:course_id", h.getCourse)
		schools.POST("/course/:course_id", h.joinCourse)
		schools.GET("/module/:module_id", h.getModule)
		schools.GET("/lesson/:lesson_id", h.getLesson)
	}
}
