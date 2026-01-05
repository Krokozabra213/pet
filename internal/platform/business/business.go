package business

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	platformconfig "github.com/Krokozabra213/sso/newconfigs/platform"
)

type ISchools interface {
	GetAllPublished(ctx context.Context) ([]domain.AllSchoolOutput, error)
	GetByID(ctx context.Context, id string) (domain.SchoolOutput, error)
}

type ICourses interface {
	GetByID(ctx context.Context, id string) (domain.CourseOutput, error)
}

type IModules interface {
	GetByCourseID(ctx context.Context, id string) ([]domain.CourseModuleOutput, error)
	GetByID(ctx context.Context, id string) (domain.ModuleOutput, error)
}

type ILessons interface {
	GetByID(ctx context.Context, id string) (domain.LessonOutput, error)
}

type IAuth interface {
	RefreshTokens(refreshToken string) (string, string, error)
}

type Business struct {
	Schools ISchools
	Courses ICourses
	Modules IModules
	Lessons ILessons
	Auth    IAuth
}

func New(cfg *platformconfig.Config) *Business {
	return &Business{
		Schools: NewSchoolsService(nil),
		Courses: NewCoursesService(nil),
		Modules: NewModulesService(nil),
		Lessons: NewLessonsService(nil),
	}
}
