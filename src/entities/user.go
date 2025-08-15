package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password" json:"-"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdateAt time.Time `bson:"updateAt" json:"updateAt"`
}