package newdomain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AllSchoolOutput struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	CreatedAt   time.Time
	Info        Info
}

type SchoolOutput struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	CreatedAt   time.Time
	Admins      []Admin
	Courses     []SchoolCourseOutput
	Info        Info
}

type SchoolCourseOutput struct {
	ID          primitive.ObjectID
	SchoolID    primitive.ObjectID
	Name        string
	Description string
	Color       string
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
