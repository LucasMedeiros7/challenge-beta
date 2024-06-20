package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/LucasMedeiros7/challenge-beta/internal/models"

	"github.com/LucasMedeiros7/challenge-beta/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	var err error

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}
}

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderId int
	err := db.QueryRow("INSERT INTO Pedidos (clienteId, status) VALUES ($1, 'PENDENTE') RETURNING pedidoId", order.ClienteId).Scan(&orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, item := range order.Itens {
		_, err = db.Exec("INSERT INTO ItensPedido (pedidoId, produtoId, quantidade) VALUES ($1, $2, $3)", orderId, item.ProdutoId, item.Quantidade)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	order.PedidoId = orderId
	services.SendOrderToQueue(order)
	c.JSON(http.StatusCreated, order)
}

func GetOrder(c *gin.Context) {
	pedidoId, err := strconv.Atoi(c.Param("pedidoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pedidoId"})
		return
	}

	var order models.Order
	err = db.QueryRow("SELECT pedidoId, clienteId, status FROM Pedidos WHERE pedidoId = $1", pedidoId).Scan(&order.PedidoId, &order.ClienteId, &order.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	rows, err := db.Query("SELECT itemId, produtoId, quantidade FROM ItensPedido WHERE pedidoId = $1", order.PedidoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item models.ItemPedido
		if err := rows.Scan(&item.ItemId, &item.ProdutoId, &item.Quantidade); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		order.Itens = append(order.Itens, item)
	}

	c.JSON(http.StatusOK, order)
}

func ListOrders(c *gin.Context) {
	rows, err := db.Query("SELECT pedidoId, clienteId, status FROM Pedidos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.PedidoId, &order.ClienteId, &order.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}

func StartOrderProcessor() {
	services.ProcessOrders()
}
