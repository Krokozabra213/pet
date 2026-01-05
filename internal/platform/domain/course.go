package newdomain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID          primitive.ObjectID
	SchoolID    primitive.ObjectID
	Name        string
	Description string
	Color       string
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Modules     []ModulesSubCollection
	Published   bool
}

type ModulesSubCollection struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	Position    uint
}
