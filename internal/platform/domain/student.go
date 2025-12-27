package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID               uint64
	Name             string
	Email            string
	LastVisitAt      time.Time
	SchoolId         primitive.ObjectID
	AvailableCourses primitive.ObjectID
	Banned           bool
}
