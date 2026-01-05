package business

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
)

type CoursesService struct {
	repo repository.Courses
}

func NewCoursesService(repo repository.Courses) *CoursesService {
	return &CoursesService{repo: repo}
}

func (s *CoursesService) GetAllPublished(ctx context.Context) ([]domain.AllSchoolOutput, error) {
	schools, err := s.repo.GetAllPublished(ctx)
	return schools, err
}

func (s *SchoolsService) GetByID(ctx context.Context, id string) (domain.SchoolOutput, error) {
	school, err := s.repo.GetByID(ctx, id)
	return school, err
}
