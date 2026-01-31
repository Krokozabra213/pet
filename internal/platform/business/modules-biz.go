package business

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	"github.com/Krokozabra213/sso/internal/platform/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModulesService struct {
	repo repository.IModules
}

func NewModulesService(repo repository.IModules) *ModulesService {
	return &ModulesService{repo: repo}
}

func (s *ModulesService) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Module, error) {
	module, err := s.repo.GetByID(ctx, id)
	return module, err
}
