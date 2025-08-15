package converts

import (
	"github.com/jhonnydsl/api-restful-gin/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringToObject(value string) (primitive.ObjectID, error) {
	valueObj, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return primitive.NilObjectID, utils.BadRequestError("Erro ao converter string para objeto")
	}
	return valueObj, err
}