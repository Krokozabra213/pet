package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// type Module struct {
// 	ID          primitive.ObjectID
// 	Name        string
// 	Description string
// 	Position    uint
// 	Published   bool
// 	CourseID    primitive.ObjectID
// 	Lessons     []LessonSubCollection
// }

// type LessonSubCollection struct {
// 	ID       primitive.ObjectID
// 	Name     string
// 	Position uint
// 	Content  string
// }

type Module struct {
	ID          primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Name        string                `json:"name" bson:"name"`
	Description string                `json:"description" bson:"description"`
	Position    uint                  `json:"position" bson:"position"`
	Published   bool                  `json:"published" bson:"published"`
	CourseID    primitive.ObjectID    `json:"course_id" bson:"course_id"`
	Lessons     []LessonSubCollection `json:"lessons" bson:"lessons"`
}

type LessonSubCollection struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Position uint               `json:"position" bson:"position"`
	Content  string             `json:"content" bson:"content"`
}
