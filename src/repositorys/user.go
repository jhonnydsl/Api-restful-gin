package repositorys

type UserRepository struct {
	*MongoRepositoryContext
}

func NewUserRepository(url, dbName, collectionName string) (*UserRepository, error) {
	mongoRepo, err := NewMongoRepositoryContext(url, dbName, collectionName)
	if err != nil {
		return nil, err
	}
	return &UserRepository{MongoRepositoryContext: mongoRepo}, nil
}