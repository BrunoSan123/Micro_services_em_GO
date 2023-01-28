package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"order/db"
	"order/queue"
	"os"
	"time"

	"github.com/google/uuid"
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
	in := make(chan []byte)
	connection := queue.Connect()
	queue.StartConsuming(connection, in)

	for payload := range in {
		createOrder(payload)
		fmt.Println(string(payload))
	}
}

func createOrder(payload []byte) {
	var order Order
	json.Unmarshal(payload, &order)

	order.Uuid = string(uuid.NewString())
	order.Status = "pendente"
	order.CreatedAt = time.Now()
	saveOrder(order)

}

func saveOrder(order Order) {
	json, _ := json.Marshal(order)
	connection := db.Connect()

	err := connection.Set(context.Background(), order.Uuid, string(json), 0).Err()

	if err != nil {
		log.Fatal(err)
	}
}

func getProductById(id string) Product {
	response, err := http.Get(productUrl + "/product/" + id)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	var product Product
	json.Unmarshal(data, &product)
	return product

}
