package repository

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LessonRepo struct {
	db *mongo.Collection
}

func NewLessonRepo(db *mongo.Database) *LessonRepo {
	return &LessonRepo{
		db: db.Collection(moduleCollection),
	}
}

func (r *LessonRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Lesson, error) {
	var lesson domain.Lesson

	filter := bson.M{
		"_id":       id,
		"published": true,
	}

	err := r.db.FindOne(ctx, filter).Decode(&lesson)

	return &lesson, err
}
