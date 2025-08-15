# API de Autenticação e Tarefas

API RESTful desenvolvida em **Go** utilizando o framework **Gin** e banco de dados **MongoDB**, com autenticação via **JWT** e documentação **Swagger**.

## 🚀 Tecnologias

- [Go](https://go.dev/)
- [Gin](https://gin-gonic.com/)
- [MongoDB](https://www.mongodb.com/)
- [JWT](https://jwt.io/)
- [Swagger](https://swagger.io/)

## 📦 Instalação

1. **Clonar o repositório**  
   `git clone https://github.com/jhonnydsl/api-restful-gin.git && cd api-restful-gin`

2. **Instalar dependências**  
   `go mod tidy`

3. **Configurar variáveis de ambiente**  
   Criar um arquivo `.env` na raiz do projeto com:  
   `MONGO_URI=mongodb://localhost:27017`  
   `MONGO_DB_NAME=nome_do_banco`  
   `MONGO_COLLECTION=nome_da_colecao`  
   `JWT_SECRET=seu_segredo_jwt`  
   `PORT=8080`

4. **Rodar o projeto**  
   `go run main.go`

## 📂 Estrutura do projeto

.
├── src
│ ├── controllers # Lógica dos endpoints
│ ├── repositorys # Conexão e operações com MongoDB
│ ├── utils
│ │ └── middlewares # Middlewares (CORS, JWT, etc)
│ └── ...
├── .env # Variáveis de ambiente (não subir no GitHub)
├── go.mod
├── go.sum
└── main.go

Swagger disponível em:  
`http://localhost:8080/swagger/index.html`
