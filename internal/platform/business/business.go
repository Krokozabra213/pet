package business

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	"github.com/Krokozabra213/sso/internal/platform/repository"
	platformconfig "github.com/Krokozabra213/sso/newconfigs/platform"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISchools interface {
	GetAllPublished(ctx context.Context) ([]domain.AllSchoolOutput, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.School, error)
}

type ICourses interface {
	GetAllPublished(ctx context.Context) ([]domain.AllCourseOutput, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Course, error)
}

type IModules interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Module, error)
}

type ILessons interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Lesson, error)
}

type IAuth interface {
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type Business struct {
	Schools ISchools
	Courses ICourses
	Modules IModules
	Lessons ILessons
	Auth    IAuth
}

type Deps struct {
	Config *platformconfig.Config
	Repos  *repository.Repositories
}

func New(deps Deps) *Business {
	return &Business{
		Schools: NewSchoolsService(deps.Repos.Schools),
		Courses: NewCoursesService(deps.Repos.Courses),
		Modules: NewModulesService(deps.Repos.Modules),
		Lessons: NewLessonsService(deps.Repos.Lessons),
	}
}
