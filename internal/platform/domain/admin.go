package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Admin struct {
	UserID   uint64
	Name     string
	Email    string
	SchoolID primitive.ObjectID
}
