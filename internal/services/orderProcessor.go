package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/LucasMedeiros7/challenge-beta/internal/models"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

var db *sql.DB

func init() {
	doterr := godotenv.Load()
	if doterr != nil {
		log.Fatal("Error loading .env file")
	}

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
}
