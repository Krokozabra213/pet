package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourseOutput struct {
	ID          primitive.ObjectID
	SchoolID    primitive.ObjectID
	Name        string
	Description string
	Color       string
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Modules     []CourseModuleOutput // добавленное поле
}

type CourseModuleOutput struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	Position    uint
	CourseID    primitive.ObjectID
}

type ModuleOutput struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	Position    uint
	CourseID    primitive.ObjectID
	Lessons     []ModuleLessonOutput
}

type ModuleLessonOutput struct {
	ID       primitive.ObjectID
	Name     string
	Position uint
	Content  string
}

// type LessonOutput struct {
// 	ID       primitive.ObjectID
// 	Name     string
// 	Position uint
// 	Content  string
// 	SchoolID primitive.ObjectID
// }
