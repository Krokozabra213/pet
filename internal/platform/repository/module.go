package repository

import (
	"context"
	"sort"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ModuleRepo struct {
	db *mongo.Collection
}

func NewModuleRepo(db *mongo.Database) *ModuleRepo {
	return &ModuleRepo{
		db: db.Collection(moduleCollection),
	}
}

func (r *ModuleRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Module, error) {
	var module domain.Module

	filter := bson.M{
		"_id":       id,
		"published": true,
	}

	err := r.db.FindOne(ctx, filter).Decode(&module)

	sort.Slice(module.Lessons, func(i, j int) bool {
		return module.Lessons[i].Position < module.Lessons[j].Position
	})

	return &module, err
}
