package repository

import (
	"context"
	"sort"

	"github.com/Krokozabra213/sso/internal/platform/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseRepo struct {
	db *mongo.Collection
}

func NewCourseRepo(db *mongo.Database) *CourseRepo {
	return &CourseRepo{
		db: db.Collection(courseCollection),
	}
}

func (r *CourseRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Course, error) {
	var course domain.Course

	filter := bson.M{
		"_id":       id,
		"published": true,
	}

	err := r.db.FindOne(ctx, filter).Decode(&course)

	sort.Slice(course.Modules, func(i, j int) bool {
		return course.Modules[i].Position < course.Modules[j].Position
	})

	return &course, err
}
