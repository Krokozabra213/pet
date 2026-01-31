package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type AllCourseOutput struct {
// 	ID          primitive.ObjectID
// 	SchoolID    primitive.ObjectID
// 	Name        string
// 	Description string
// 	Color       string
// 	ImageURL    string
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time
// }

type AllCourseOutput struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SchoolID    primitive.ObjectID `json:"school_id" bson:"school_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Color       string             `json:"color" bson:"color"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}