package newdomain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Lesson struct {
	ID        primitive.ObjectID
	Name      string
	Position  uint
	Published bool
	Content   string
	ModuleID  primitive.ObjectID
}
