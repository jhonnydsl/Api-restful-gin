package repositorys

import (
	"context"
	"fmt"

	"github.com/jhonnydsl/api-restful-gin/src/dtos"
	dtosPage "github.com/jhonnydsl/api-restful-gin/src/dtos/pagination"
	"github.com/jhonnydsl/api-restful-gin/src/utils"
	"github.com/jhonnydsl/api-restful-gin/src/utils/converts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepositoryContext struct {
	*mongo.Collection
	client *mongo.Client
}

func NewMongoRepositoryContext(url, dbname, collectionName string) (*MongoRepositoryContext, error) {
	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, utils.InternalServerError(fmt.Sprintf("Erro ao conectar com o mongoDB: %v", err))
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, utils.InternalServerError(fmt.Sprintf("Erro ao pingar mongoDB: %v", err))
	}

	colection := client.Database(dbname).Collection(collectionName)

	return &MongoRepositoryContext{
		Collection: colection,
		client: client,
	}, nil
}

func (r *MongoRepositoryContext) Create(contextServer context.Context, document interface{}) error {
	_, err := r.Collection.InsertOne(contextServer, document)	// <= Chama o insertOne na coleção do mongoDB e tenta inserir o documento passado pela coleção.
	if err != nil {
		return utils.BadRequestError(fmt.Sprintf("Erro ao inserir documento: %v", err))
	}
	return nil
}

func (r *MongoRepositoryContext) ExistsByAny(ctx context.Context, params dtos.ExistsFilter) error {
	filter := bson.M{params.Field: params.Value}	// <= Cria um filtro mongoDB para buscar documentos onde Field seja igual a Value.

	if params.ForeignKey != "" && params.ForeignKeyValue != primitive.NilObjectID {		// <= Adiciona um filtro composto de houver ForeignKey.
		filter["$and"] = []bson.M{
			{params.Field: params.Value},
			{params.ForeignKey: params.ForeignKeyValue},
		}
	}

	err := r.Collection.FindOne(ctx, filter).Err()	// <= Usa findOne para procurar um documento que combine com o filtro.
	if err == mongo.ErrNoDocuments {
		return nil // Não existe, esta ok
	}
	if err != nil {
		return utils.InternalServerError(fmt.Sprintf("Erro ao verificar existencia: %v", err))
	}

	return utils.ConflictError(fmt.Sprintf("Ja possui cadastro: %v", params.Value))
}

func (r *MongoRepositoryContext) GetByAny(ctx context.Context, params dtos.GetAnyFilter) error {
	filter := bson.M{params.Field: params.Value}

	if params.ForeignKey != "" && params.ForeignKeyValue != primitive.NilObjectID {
		filter["$and"] = []bson.M{
			{params.Field: params.Value},	// <= Filtrar pelo userID.
			{params.ForeignKey: params.ForeignKeyValue},
		}
	}

	err := r.Collection.FindOne(ctx, filter).Decode(params.Result)
	if err == mongo.ErrNoDocuments {
		return utils.NotFoundError("Registro não encontrado")
	} else if err != nil {
		return utils.InternalServerError(fmt.Sprintf("Erro ao buscar registro: %v", err))
	}

	return nil
}

func (r *MongoRepositoryContext) GetPagination(ctx context.Context, params dtosPage.PaginationParams) dtosPage.PaginationResultContext {
	filter := bson.M{params.Field: params.Value}	// <= Filtro inicial pelo userID.

	if params.SearchValue != "" {
		filter["$and"] = []bson.M{
			{params.Field: params.Value},	// <= Filtrar pelo userID.
			{
				params.SearchField: bson.M{
					"$regex": params.SearchValue,
					"$options": "i",	// <= Insensivel a maiusculas/minusculas.
				},
			},
		}
	}

	// Conta o número total de documentos que correspondem ao filtro.
	count, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return dtosPage.PaginationResultContext{
			TotalPages: 0,
			HasNextPage: false,
			TotalItems: 0,
			Err: utils.InternalServerError(fmt.Sprintf("Erro ao contar documentos: %v", err)),
		}
	}

	// Calculando o total de paginas.
	totalPages := int(count) / params.Limit
	if int(count)%params.Limit > 0 {
		totalPages++
	}

	// Garantir que a página seja pelo menos 1.
	if params.Skip < 1 {
		params.Skip = 1
	}

	// Definir valores padrão para os parâmetros opcionais.
	if params.SortField == "" {
		params.SortField = "_id"	// <= Campo padrão para ordenação.
	}

	// Calculando o skip corretamente.
	skipInit := (params.Skip -1) * params.Limit

	// Verificando se há uma próxima página.
	hasNextPage := params.Skip < totalPages

	cursor, err := r.Collection.Find(
		ctx,
		filter,
		options.Find().SetSkip(int64(skipInit)).SetLimit(int64(params.Limit)).SetSort(bson.M{params.SortField: params.SortOrder}),
	)
	if err != nil {
		return dtosPage.PaginationResultContext{
			TotalPages: 0,
			HasNextPage: false,
			TotalItems: skipInit,
			Err: utils.InternalServerError(fmt.Sprintf("Erro ao buscar registros: %v", err)),
		}
	}
	defer cursor.Close(ctx)

	// Decodificando os documentos encontrados no resultado.
	err = cursor.All(ctx, params.Result)
	if err != nil {
		return dtosPage.PaginationResultContext{
			TotalPages: 0,
			HasNextPage: false,
			TotalItems: skipInit,
			Err: utils.InternalServerError(fmt.Sprintf("Erro ao decodificar registros: %v", err)),
		}
	}

	// Retorna os resultados com total de páginas e se há proxima página.
	return dtosPage.PaginationResultContext{
		TotalPages: totalPages,
		HasNextPage: hasNextPage,
		TotalItems: int(count),
		Err: nil,
	}
}

func (r *MongoRepositoryContext) Update(contextServer context.Context, params dtos.UpdateFilter) error {
	filter := bson.M{"_id": params.ID}
	convertedToMap := converts.MapTokeyAndValueUpdate(params.Dto)
	updateDoc := bson.M{"$set": convertedToMap}

	if params.ForeignKey != "" && params.ForeignKeyValue != primitive.NilObjectID {
		filter["$and"] = []bson.M{
			{"_id": params.ID},
			{params.ForeignKey: params.ForeignKeyValue},
		}
	}

	result, err := r.Collection.UpdateOne(contextServer, filter, updateDoc)
	if err != nil {
		return utils.BadRequestError(fmt.Sprintf("Erro ao atualizar documento: %v", err))
	}

	if result.MatchedCount == 0 {
		return utils.NotFoundError("Nenhum documento encontrado para atualizar.")
	}

	return nil
}

func (r *MongoRepositoryContext) Delete(contextServer context.Context, params dtos.DeleteFilter) error {
	filter := bson.M{"_id": params.ID}

	if params.ForeignKey != "" && params.ForeignKeyValue != primitive.NilObjectID {
		filter["$and"] = []bson.M{
			{"_id": params.ID},
			{params.ForeignKey: params.ForeignKeyValue},
		}
	}

	result, err := r.Collection.DeleteOne(contextServer, filter)
	if err != nil {
		return utils.BadRequestError(fmt.Sprintf("Erro ao deletar documento: %v", err))
	}

	if result.DeletedCount == 0 {
		return utils.NotFoundError("Nenhum documento encontrado para deletar.")
	}

	return nil
}