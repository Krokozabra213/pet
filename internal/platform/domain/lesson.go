package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// type Lesson struct {
// 	ID        primitive.ObjectID
// 	Name      string
// 	Position  uint
// 	Published bool
// 	Content   string
// 	ModuleID  primitive.ObjectID
// }

type Lesson struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Position  uint               `json:"position" bson:"position"`
	Published bool               `json:"published" bson:"published"`
	Content   string             `json:"content" bson:"content"`
	ModuleID  primitive.ObjectID `json:"module_id" bson:"module_id"`
}
