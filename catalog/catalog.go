package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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

var productsUrl string

func init() {
	productsUrl = os.Getenv("PRODUCT_URL")
}

func loadProducts() []Product {
	res, err := http.Get(productsUrl + "/products")

	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	var products Products
	json.Unmarshal(data, &products)
	fmt.Println(string(data))
	return products.Products
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	products := loadProducts()
	t := template.Must(template.ParseFiles("template/catalog.html"))
	t.Execute(w, products)
}

func showProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err := http.Get(productsUrl + "/product/" + vars["id"])
	if err != nil {
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(res.Body)

	var product Product
	json.Unmarshal(data, &product)
	t := template.Must(template.ParseFiles("template/view.html"))
	t.Execute(w, product)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", listProducts)
	r.HandleFunc("/product/{id}", showProducts)
	http.ListenAndServe(":8081", r)
}
