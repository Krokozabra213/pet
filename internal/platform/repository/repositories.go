package repository

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ISchools interface {
	GetAllPublished(ctx context.Context) ([]domain.AllSchoolOutput, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.School, error)
}

type ICourses interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Course, error)
}

type IModules interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Module, error)
}

type ILessons interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Lesson, error)
}

type Repositories struct {
	Schools ISchools
	Courses ICourses
	Modules IModules
	Lessons ILessons
}

func NewRepositories(mongoDB *mongo.Database) *Repositories {
	return &Repositories{
		Schools: NewSchoolsRepo(mongoDB),
		Courses: NewCourseRepo(mongoDB),
		Modules: NewModuleRepo(mongoDB),
		Lessons: NewLessonRepo(mongoDB),
	}
}
