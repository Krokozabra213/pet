package business

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
)

type LessonsService struct {
	repo repository.Lessons
}

func NewLessonsService(repo repository.Lessons) *LessonsService {
	return &LessonsService{repo: repo}
}

func (s *LessonsService) GetByID(ctx context.Context, id string) (domain.LessonOutput, error) {
	lesson, err := s.repo.GetByID(ctx, id)
	return lesson, err
}
