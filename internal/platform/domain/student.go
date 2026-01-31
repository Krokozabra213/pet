package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID          uint64
	Name        string
	Email       string
	LastVisitAt time.Time
	SchoolId    primitive.ObjectID
	Courses     []CoursesShortInfo
	Banned      bool
}

type CoursesShortInfo struct {
	ID       primitive.ObjectID
	Name     string
	Color    string
	ImageURL string
}
