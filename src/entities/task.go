package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID primitive.ObjectID 	`bson:"_id,omitempty" json:"ID"`
	Title string 			`bson:"title" json:"title"`
	Description string 		`bson:"description" json:"description"`
	UserID primitive.ObjectID `bson:"userID" json:"userID"`
	CreateAt time.Time 		`bson:"createAt" json:"createAt"`
	UpdateAt time.Time 		`bson:"updateAt" json:"updateAt"`
}