package business

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	"github.com/Krokozabra213/sso/internal/platform/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoursesService struct {
	repo repository.ICourses
}

func NewCoursesService(repo repository.ICourses) *CoursesService {
	return &CoursesService{repo: repo}
}

func (s *CoursesService) GetAllPublished(ctx context.Context) ([]domain.AllCourseOutput, error) {
	courses, err := s.repo.GetAllPublished(ctx)
	return courses, err
}

func (s *CoursesService) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Course, error) {
	course, err := s.repo.GetByID(ctx, id)
	return course, err
}
