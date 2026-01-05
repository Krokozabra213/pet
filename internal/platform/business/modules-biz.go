package business

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
)

type ModulesService struct {
	repo repository.Modules
}

func NewModulesService(repo repository.Modules) *ModulesService {
	return &ModulesService{repo: repo}
}

func (s *ModulesService) GetByCourseID(ctx context.Context, id string) ([]domain.CourseModuleOutput, error) {
	modules, err := s.repo.GetByCourseID(ctx, id)
	return modules, err
}

func (s *ModulesService) GetByID(ctx context.Context, id string) (domain.ModuleOutput, error) {
	module, err := s.repo.GetByID(ctx, id)
	return module, err
}
