package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/LucasMedeiros7/challenge-beta/internal/models"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
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

func SendOrderToQueue(order models.Order) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"FilaPedidos",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	body, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatal(err)
	}
}

func ProcessOrders() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"FilaPedidos",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var order models.Order
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}

			// Processar o pedido (exemplo: marcar como processado)
			_, err = db.Exec("UPDATE Pedidos SET status = 'PROCESSADO' WHERE pedidoId = $1", order.PedidoId)
			if err != nil {
				log.Printf("Error updating order: %s", err)
			} else {
				log.Printf("Order %d processed", order.PedidoId)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
