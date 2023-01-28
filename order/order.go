package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"order/db"
	"order/queue"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float32 `json:"price,string"`
}

type Order struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ProductId string    `json:"product_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,string"`
}

var productUrl string

func init() {
	productUrl = os.Getenv("PRODUCT_URL")

}

func main() {
	var param string
	flag.StringVar(&param, "opt", "", "Usage")
	flag.Parse()
	in := make(chan []byte)
	connection := queue.Connect()

	switch param {
	case "checkout":
		queue.StartConsuming("checkout_queue", connection, in)
		for payload := range in {
			notifyOrderCreated(createOrder(payload), connection)
			fmt.Println(string(payload))
		}
	case "payment":
		queue.StartConsuming("payment_queue", connection, in)
		var order Order
		for payload := range in {
			json.Unmarshal(payload, &order)
			saveOrder(order)
			fmt.Println("Payment", string(payload))
		}

	}

}

func createOrder(payload []byte) Order {
	var order Order
	json.Unmarshal(payload, &order)

	order.Uuid = string(uuid.NewString())
	order.Status = "pendente"
	order.CreatedAt = time.Now()
	saveOrder(order)
	return order
}

func notifyOrderCreated(order Order, ch *amqp.Channel) {
	json, _ := json.Marshal(order)
	queue.Notify(json, "order_ex", "", ch)
}

func saveOrder(order Order) {
	json, _ := json.Marshal(order)
	connection := db.Connect()

	err := connection.Set(context.Background(), order.Uuid, string(json), 0).Err()

	if err != nil {
		log.Fatal(err)
	}
}
