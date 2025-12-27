package httpv1

import (
	"github.com/Krokozabra213/sso/internal/platform/business"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	busines *business.Business
}

func NewHandler(business *business.Business) *Handler {
	return &Handler{
		busines: business,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initSchools(v1)
	}
}

func (h *Handler) initSchools(v1 *gin.RouterGroup) {
	schools := v1.Group("/schools")
	{
		schools.GET("")    // достаем все школы
		schools.GET("/id") // достаем одну школу и показываем все курсы

		courses := schools.Group("/courses")
		{
			courses.GET("/:id") // достаем какой то курс

			modules := courses.Group("/modules")
			{
				modules.GET("/:id") // достаем модуль и показываем уроки

				lessons := modules.Group("/lessons")
				{
					lessons.GET("/:id") // показываем какойто урок
				}
			}
		}
	}
}
