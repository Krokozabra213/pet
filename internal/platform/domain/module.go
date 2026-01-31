package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Module struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	Position    uint
	Published   bool
	CourseID    primitive.ObjectID
	Lessons     []LessonSubCollection
}

type LessonSubCollection struct {
	ID       primitive.ObjectID
	Name     string
	Position uint
	Content  string
}
