package repository

import (
	"context"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SchoolsRepo struct {
	db *mongo.Collection
}

func NewSchoolsRepo(db *mongo.Database) *SchoolsRepo {
	return &SchoolsRepo{
		db: db.Collection(schoolsCollection),
	}
}

func (r *SchoolsRepo) GetAllPublished(ctx context.Context) ([]domain.AllSchoolOutput, error) {
	var schools []domain.AllSchoolOutput

	filter := bson.M{
		"published": true,
	}

	cur, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &schools)

	return schools, err
}

func (r *SchoolsRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.School, error) {
	var school domain.School

	filter := bson.M{
		"_id":       id,
		"published": true,
	}

	err := r.db.FindOne(ctx, filter).Decode(&school)

	return &school, err
}
