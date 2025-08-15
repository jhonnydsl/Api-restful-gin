package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/api-restful-gin/src/dtos"
	"github.com/jhonnydsl/api-restful-gin/src/repositorys"
	"github.com/jhonnydsl/api-restful-gin/src/services"
	"github.com/jhonnydsl/api-restful-gin/src/utils"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(server *gin.Engine, repo *repositorys.UserRepository) {
	service := services.NewUserService(repo)
	controller := &UserController{service: service}

	routes := server.Group("/users")
	{
		routes.POST("", controller.CreateUser)
		routes.POST("/login", controller.LoginUser)
	}
}

// @Summary Criar um novo usuário
// @Description Registra um novo usuário na API
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.User true "Dados do usuário"
// @Success 201 {object} dtos.Message "Usuário criado"
// @Failure 400 {object} dtos.APIError "Erro de validação"
// @Failure 409 {object} dtos.APIError "Erro de conflito, dados ja existem"
// @Router /users [post]
func (c *UserController) CreateUser(ginContext *gin.Context) {
	var userDto dtos.User

	err := utils.ValidateRequestBody(ginContext, &userDto)	// <= Valida o corpo da requisição.
	if err != nil {
		ginContext.Error(err)
		return
	}

	err = c.service.ExistsUserByEmail(dtos.ExistsFilter{
		Field: "email",
		Value: userDto.Email,
	})
	if err != nil {
		ginContext.Error(err)
		return
	}

	err = c.service.CreateUser(userDto.Email, userDto.Password)
	if err != nil {
		ginContext.Error(err)
		return
	}

	ginContext.JSON(http.StatusCreated, dtos.Message{
		Message: "Usuário criado com sucesso.",
	})
}

// @Summary Criar o login do usuário
// @Description Fazer o login do usuário
// @Tags users
// @Accept json
// @Produce json
// @Param user body dtos.User true "Dados do usuário"
// @Success 200 {object} dtos.Token "Usuário login"
// @Failure 400 {object} dtos.APIError "Erro de login"
// @Router /users/login [post]
func (c *UserController) LoginUser(ginContext *gin.Context) {
	var userDto dtos.User

	err := utils.ValidateRequestBody(ginContext, &userDto)
	if err != nil {
		ginContext.Error(err)
		return
	}

	token, err := c.service.LoginUser(userDto.Email, userDto.Password)
	if err != nil {
		ginContext.Error(err)
		return
	}

	ginContext.JSON(http.StatusOK, dtos.Token{Token: token})
}