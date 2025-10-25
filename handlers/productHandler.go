package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/dleonardoq/meli/models"
	"github.com/dleonardoq/meli/storage"

	"github.com/gorilla/mux"
)

var logger = log.New(os.Stdout, "[MELI-API] ", log.LstdFlags|log.Lshortfile)
var store, err = storage.NewProduct("data/products.json")

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating storage"))
		return
	}

	products := store.GetAllProducts()

	response := models.SuccessResponse{
		Products: products,
		Count:    len(products),
	}

	jsonResponse(w, http.StatusOK, response)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error creating storage")
		return
	}

	product, err := store.GetProductById(id)

	if err != nil {
		errorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	jsonResponse(w, http.StatusOK, product)
}

func GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error creating storage")
		return
	}

	products := store.GetProductsByCategory(category)

	response := models.SuccessResponse{
		Products: products,
		Count:    len(products),
	}

	jsonResponse(w, http.StatusOK, response)
}

func CompareProducts(w http.ResponseWriter, r *http.Request) {
	var request models.ComparisonRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	defer r.Body.Close()

	if len(request.ProductIDs) == 0 {
		errorResponse(w, http.StatusBadRequest, "No product IDs provided")
		return
	}

	if len(request.ProductIDs) > 10 {
		errorResponse(w, http.StatusBadRequest, "Cannot compare more than 10 products")
		return
	}
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error creating storage")
		return
	}

	products, err := store.GetProductsByIds(request.ProductIDs)
	if err != nil {
		logger.Printf("Warning: %v", err)
	}

	if len(products) == 0 {
		errorResponse(w, http.StatusNotFound, "None of the requested products were found")
		return
	}

	response := models.SuccessResponse{
		Products: products,
		Count:    len(products),
	}

	statusCode := http.StatusOK
	if len(products) < len(request.ProductIDs) {
		statusCode = http.StatusPartialContent
	}

	jsonResponse(w, statusCode, response)
}

func errorResponse(w http.ResponseWriter, code int, message string) {
	logger.Printf("Error: %s (code: %d)", message, code)

	errorResponse := models.ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
		Code:    code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		logger.Printf("Error encoding JSON response: %v", err)
	}
}

// respondWithJSON envía una respuesta exitosa en formato JSON
func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Printf("Error encoding JSON response: %v", err)
	}
}
