Authentication and Task API

A RESTful API built with Go, using the Gin framework and MongoDB, featuring JWT authentication and Swagger documentation.

🚀 Technologies

Go

Gin

MongoDB

JWT

Swagger

📦 Installation

1. Clone the repository
   git clone https://github.com/jhonnydsl/api-restful-gin.git && cd api-restful-gin

2. Install dependencies
   go mod tidy

3. Set up environment variables
   Create a .env file in the project root with:

MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=your_database_name
MONGO_COLLECTION=your_collection_name
JWT_SECRET=your_jwt_secret
PORT=8080

4. Run the project
   go run main.go

📂 Project structure
.
├── src
│ ├── controllers # Endpoint logic
│ ├── repositorys # MongoDB connection and operations
│ ├── utils
│ │ └── middlewares # Middlewares (CORS, JWT, etc.)
│ └── ...
├── .env # Environment variables (do not commit to GitHub)
├── go.mod
├── go.sum
└── main.go

Swagger available at: http://localhost:8080/swagger/index.html
