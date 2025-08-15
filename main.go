package main

import (
	"log"
	"os"

	_ "github.com/jhonnydsl/api-restful-gin/docs"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/api-restful-gin/src/controllers"
	"github.com/jhonnydsl/api-restful-gin/src/repositorys"
	"github.com/jhonnydsl/api-restful-gin/src/utils/middlewares"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}
}

// @title API de Autenticação e Tarefas
// @version 1.0
// @description API para gerenciamento de usuários e tarefas usando Gin e MongoDB
// @host localhost:8080
// @BasePath /
func main() {
	url := os.Getenv("URL_DB")
	dbName := os.Getenv("NAME_DB")

	repoUser, errUser := repositorys.NewUserRepository(url, dbName, "users")
	repoTask, errTask := repositorys.NewTaskRepository(url, dbName, "tasks")

	if errUser != nil || errTask != nil {
		log.Fatalf("Erro no repositorio ao iniciar: errUser = %v, errTask = %v", errUser, errTask)
		return
	}

	server := gin.Default()
	server.Use(middlewares.CorsMiddlewares())
	server.Use(middlewares.ErrorMidlewareHandle())
	server.Use(middlewares.JWTAuthMiddleware())

	controllers.NewUserController(server, repoUser)
	controllers.NewTaskController(server, repoTask)

	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization
	// @description Value: Bearer abc... (Bearer+space+token)
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))

	// Iniciar servidor na porta 8080
	server.Run(":8080")
}