package repositorys

type TaskRepository struct {
	*MongoRepositoryContext
}

func NewTaskRepository(url, dbName, collectionName string) (*TaskRepository, error) {
	// Cria um MongoRepositoryContext, que conecta ao MongoDB e fornece métodos genéricos para CRUD.
	mongoRepo, err := NewMongoRepositoryContext(url, dbName, collectionName)
	if err != nil {
		return nil, err
	}
	return &TaskRepository{MongoRepositoryContext: mongoRepo}, nil
}