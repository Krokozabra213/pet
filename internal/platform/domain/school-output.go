package domain

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
