package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Product data structure
type Product struct {
	ID                int      `json:"id"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	Price             float64  `json:"price"`
	DiscountPercentage float64  `json:"discountPercentage"`
	Rating            float64  `json:"rating"`
	Stock             int      `json:"stock"`
	Brand             string   `json:"brand"`
	Category          string   `json:"category"`
	Thumbnail         string   `json:"thumbnail"`
	Images            []string `json:"images"`
}

// In-memory storage for demonstration purposes
var products []Product

func main() {
	// Add the provided product data to the products slice
	products = []Product{
		{
			ID:                1,
			Title:             "iPhone 9",
			Description:       "An apple mobile which is nothing like apple",

			Price:             549,
			DiscountPercentage: 12.96,
			Rating:            4.69,
			Stock:             94,
			Brand:             "PEARRRRRRRRRRRRRR",
			Category:          "smartphones",
			Thumbnail:         "https://i.dummyjson.com/data/products/1/thumbnail.jpg",
			Images: []string{
				"https://i.dummyjson.com/data/products/1/1.jpg",
				"https://i.dummyjson.com/data/products/1/2.jpg",
				"https://i.dummyjson.com/data/products/1/3.jpg",
				"https://i.dummyjson.com/data/products/1/4.jpg",
				"https://i.dummyjson.com/data/products/1/thumbnail.jpg",
			},
		},
		// Add other products similarly...
	}

	// Define your REST endpoints using Gorilla Mux router
	r := mux.NewRouter()
	r.HandleFunc("/products", getProductsHandler).Methods("GET")
	r.HandleFunc("/products", createProductHandler).Methods("POST")
	r.HandleFunc("/products/{id:[0-9]+}", updateProductHandler).Methods("PUT", "PATCH")
	r.HandleFunc("/products/{id:[0-9]+}", deleteProductHandler).Methods("DELETE")

	// Attach the router to the default server
	http.Handle("/", r)

	// Start the server on port 8080
	port := 8080
	fmt.Printf("Server listening on :%d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// Handler for handling GET requests to /products
func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response content type
	w.Header().Set("Content-Type", "application/json")

	// Encode the products slice to JSON and write it to the response
	err := json.NewEncoder(w).Encode(products)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

// Handler for handling POST requests to /products
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a Product struct
	var newProduct Product
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Validate input (simple validation for demonstration purposes)
	if newProduct.Title == "" {
		http.Error(w, "Product title cannot be empty", http.StatusBadRequest)
		return
	}

	// Assign a unique ID (for demonstration purposes)
	newProduct.ID = len(products) + 1

	// Append the new product to the products slice
	products = append(products, newProduct)

	// Set the response content type
	w.Header().Set("Content-Type", "application/json")

	// Encode the newly created product to JSON and write it to the response
	err = json.NewEncoder(w).Encode(newProduct)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

// Handler for handling PUT or PATCH requests to /products/{id}
func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the request URL
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Find the product with the specified ID in the products slice
	var found bool
	var updatedProduct Product
	for i, product := range products {
		if product.ID == productID {
			found = true

			// Parse the JSON request body into a Product struct for update
			err := json.NewDecoder(r.Body).Decode(&updatedProduct)
			if err != nil {
				http.Error(w, "Error decoding JSON", http.StatusBadRequest)
				return
			}

			// Update the existing product with the new data
			products[i] = updatedProduct

			// Set the response content type
			w.Header().Set("Content-Type", "application/json")

			// Encode the updated product to JSON and write it to the response
			err = json.NewEncoder(w).Encode(updatedProduct)
			if err != nil {
				http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
				return
			}

			break
		}
	}

	if !found {
		http.Error(w, "Product not found", http.StatusNotFound)
	}
}

// Handler for handling DELETE requests to /products/{id}
func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the request URL
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Find the index of the product with the specified ID in the products slice
	var found bool
	for i, product := range products {
		if product.ID == productID {
			found = true

			// Remove the product from the products slice
			products = append(products[:i], products[i+1:]...)

			// Respond with a success message
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	if !found {
		http.Error(w, "Product not found", http.StatusNotFound)
	}
}
