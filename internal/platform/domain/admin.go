package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// type Admin struct {
// 	UserID   uint64
// 	Name     string
// 	Email    string
// 	SchoolID primitive.ObjectID
// }

type Admin struct {
	UserID    uint64               `json:"user_id" bson:"user_id"`
	Name      string               `json:"name" bson:"name"`
	Email     string               `json:"email" bson:"email"`
	SchoolIDs []primitive.ObjectID `json:"school_ids" bson:"school_ids"`
}
