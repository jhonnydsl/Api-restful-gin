package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpdateFilter struct {
	ID primitive.ObjectID
	Dto interface{}
	ForeignKey string
	ForeignKeyValue primitive.ObjectID
}