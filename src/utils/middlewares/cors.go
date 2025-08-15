package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CorsMiddlewares() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Em produção, troque por dominios especificos ("http://meusite.com").
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // ou "http://localhost:8080"
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Se for uma requisição OPTIONS, responde direto
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}