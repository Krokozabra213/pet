package domain

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
	Published   bool
}

type Module struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	Position    uint
	Published   bool
	CourseID    primitive.ObjectID
	Lessons     []Lesson
}

type Lesson struct {
	ID        primitive.ObjectID
	Name      string
	Position  uint
	Published bool
	Content   string
	SchoolID  primitive.ObjectID
}
