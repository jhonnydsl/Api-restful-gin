package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jhonnydsl/api-restful-gin/src/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
var SecretKey []byte

func init() {
	_ = godotenv.Load()
	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}

type Claims struct {
	Email  string `json:"email"`
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, userID primitive.ObjectID) (string, error) {
	expitationTime := utils.TimeNowBrazil().Add(24 * time.Hour)	// <= Calcula a expiração do token, 24h depois da hora atual no fuso do Brasil.

	claims := &Claims{
		Email: email,
		UserID: userID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expitationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)	// <= Cria um token com método de assinatura HS256 e as claims definidas.
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", utils.BadRequestError(fmt.Sprintf("Erro ao gerar o token: %v", err))
	}

	return tokenString, nil
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" && (c.Request.URL.Path == "/users" || c.Request.URL.Path == "/users/login") ||		// <= Libera rotas de criação de usuário, login e Swagger sem exigir token.
		c.Request.Method == "GET" && strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return 
		}

		// Pega o token do cabeçalho Authorization e valida se está no formato Bearer <token>.
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido ou mal formado"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Converte o UserID da claim (string) para ObjectID do MongoDB.
		objUserID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID do usuário inválido"})
			c.Abort()
			return
		}

		// Coloca email e userID no contexto do Gin, para qualquer handler poder acessar depois.
		c.Set("email", claims.Email)
		c.Set("userID", objUserID)

		c.Next()
	}
}