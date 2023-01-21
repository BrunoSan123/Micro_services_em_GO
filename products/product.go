package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float32 `json:"price,string"`
}

type Products struct {
	Products []Product
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products", ListProducts)
	r.HandleFunc("/product/{id}", getProductbyId)
	http.ListenAndServe(":8080", r)

}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := loadData()
	w.Write([]byte(products))
}

func getProductbyId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := loadData()

	var products Products
	json.Unmarshal(data, &products)

	for _, v := range products.Products {
		if v.Uuid == vars["id"] {
			product, _ := json.Marshal(v)
			w.Write([]byte(product))
		}
	}

}

func loadData() []byte {
	jsonFile, err := os.Open("products.json")

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	return data

}
