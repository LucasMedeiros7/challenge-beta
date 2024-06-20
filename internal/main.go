package main

import (
	"log"
	"os"

	"github.com/LucasMedeiros7/challenge-beta/internal/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getPostgresDSN() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || name == "" {
		log.Fatal("One or more environment variables are missing")
	}

	return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + name + "?sslmode=disable"
}

func main() {
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Verificar e obter a string de conexão do PostgreSQL
	dsn := getPostgresDSN()
	log.Println("Postgres DSN: ", dsn)

	router := gin.Default()
	router.POST("/pedidos", controllers.CreateOrder)
	router.GET("/pedidos/:pedidoId", controllers.GetOrder)
	router.GET("/pedidos", controllers.ListOrders)
	router.GET("/consume", controllers.ConsumeOrdersFromQueue)

	go controllers.StartOrderProcessor()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
