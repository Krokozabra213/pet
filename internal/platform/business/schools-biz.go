package business

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	"github.com/Krokozabra213/sso/internal/platform/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SchoolsService struct {
	repo repository.ISchools
}

func NewSchoolsService(repo repository.ISchools) *SchoolsService {
	return &SchoolsService{repo: repo}
}

func (s *SchoolsService) GetAllPublished(ctx context.Context) ([]domain.AllSchoolOutput, error) {
	schools, err := s.repo.GetAllPublished(ctx)
	return schools, err
}

func (s *SchoolsService) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.School, error) {
	school, err := s.repo.GetByID(ctx, id)
	return school, err
}
