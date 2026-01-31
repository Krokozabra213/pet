package business

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	"github.com/Krokozabra213/sso/internal/platform/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LessonsService struct {
	repo repository.ILessons
}

func NewLessonsService(repo repository.ILessons) *LessonsService {
	return &LessonsService{repo: repo}
}

func (s *LessonsService) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Lesson, error) {
	lesson, err := s.repo.GetByID(ctx, id)
	return lesson, err
}
