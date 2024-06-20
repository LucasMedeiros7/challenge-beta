package main

import (
	"os"

	"github.com/LucasMedeiros7/challenge-beta/internal/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/pedidos", controllers.CreateOrder)
	router.GET("/pedidos/:pedidoId", controllers.GetOrder)
	router.GET("/pedidos", controllers.ListOrders)

	go controllers.StartOrderProcessor()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
