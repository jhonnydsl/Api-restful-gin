package services

import (
	"fmt"
	"time"

	"github.com/jhonnydsl/api-restful-gin/src/dtos"
	"github.com/jhonnydsl/api-restful-gin/src/entities"
	"github.com/jhonnydsl/api-restful-gin/src/repositorys"
	"github.com/jhonnydsl/api-restful-gin/src/utils"
	"github.com/jhonnydsl/api-restful-gin/src/utils/middlewares"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositorys.UserRepository
}

func NewUserService(repo *repositorys.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) CreateUser(email, password string) error {
	hashPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := entities.User {
		Email: email,
		Password: hashPassword,
		CreatedAt: time.Now(),
		UpdateAt: time.Now(),
	}

	contextServer := utils.CreateContextServerWithTimeout()	// <= Cria um contexto com tempo limite para a operação no banco de dados, evitando travamentos.

	err = service.repo.Create(contextServer, user)	// <= Chama o repositório para inserir o usuário no banco usando o contexto criado.
	if err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)	// <= Converte a senha para []byte e gera o hash de senha.
	if err != nil {
		return "", utils.BadRequestError(fmt.Sprintf("Erro ao gerar a senha: %v", err))
	}
	return string(hashPassword), nil
}

func (service *UserService) ExistsUserByEmail(param dtos.ExistsFilter) error {
	contextServer := utils.CreateContextServerWithTimeout()

	err := service.repo.ExistsByAny(contextServer, param)	// <= Chama o metodo ExistsByAny do repositório para verificar se existe um registro no banco de dados que atenda ao filtro.
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) LoginUser(bodyEmail, bodyPassword string) (string, error) {
	var user entities.User

	contextServer := utils.CreateContextServerWithTimeout()
	params := dtos.GetAnyFilter{
		Field: "email",
		Value: bodyEmail,
		Result: &user,
	}

	err := service.repo.GetByAny(contextServer, params)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(bodyPassword))
	if err != nil {
		return "", utils.BadRequestError("Credenciais inválidas")
	}

	return middlewares.GenerateToken(user.Email, user.ID)
}