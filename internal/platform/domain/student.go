package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Student struct {
// 	ID          uint64
// 	Name        string
// 	Email       string
// 	LastVisitAt time.Time
// 	SchoolId    primitive.ObjectID
// 	Courses     []CoursesShortInfo
// 	Banned      bool
// }

// type CoursesShortInfo struct {
// 	ID       primitive.ObjectID
// 	Name     string
// 	Color    string
// 	ImageURL string
// }

type Student struct {
	ID          int                `json:"id" bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	Email       string             `json:"email" bson:"email"`
	LastVisitAt time.Time          `json:"last_visit_at" bson:"last_visit_at"`
	Courses     []CoursesShortInfo `json:"courses" bson:"courses"`
	Banned      bool               `json:"banned" bson:"banned"`
}

type CoursesShortInfo struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Color    string             `json:"color" bson:"color"`
	ImageURL string             `json:"image_url" bson:"image_url"`
}

type StudentSignUpInput struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type SignInInput struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}
