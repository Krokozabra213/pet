package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Course struct {
// 	ID          primitive.ObjectID
// 	SchoolID    primitive.ObjectID
// 	Name        string
// 	Description string
// 	Color       string
// 	ImageURL    string
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time
// 	Modules     []ModulesSubCollection
// 	Published   bool
// }

// type ModulesSubCollection struct {
// 	ID          primitive.ObjectID
// 	Name        string
// 	Description string
// 	Position    uint
// }

type Course struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	SchoolID    primitive.ObjectID     `json:"school_id" bson:"school_id"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Color       string                 `json:"color" bson:"color"`
	ImageURL    string                 `json:"image_url" bson:"image_url"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
	Modules     []ModulesSubCollection `json:"modules" bson:"modules"`
	Published   bool                   `json:"published" bson:"published"`
}

type ModulesSubCollection struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Position    uint               `json:"position" bson:"position"`
}