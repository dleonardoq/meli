package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dleonardoq/meli/models"
	"github.com/dleonardoq/meli/storage"

	"github.com/gorilla/mux"
)

var logger = log.New(os.Stdout, "[MELI-API] ", log.LstdFlags|log.Lshortfile)
var store *storage.ProductStorage

func init() {
	var err error
	store, err = storage.NewProduct("data/products.json")
	if err != nil {
		// Don't fatal in tests, just log the error
		logger.Printf("Warning: Failed to initialize storage: %v", err)
	}
}

// SetStorage allows tests to inject a custom storage instance
func SetStorage(s *storage.ProductStorage) {
	store = s
}

// SaveDataToDisk persists all in-memory data to the JSON file
func SaveDataToDisk() error {
	if store == nil {
		return fmt.Errorf("storage not initialized")
	}
	return store.SaveToFile()
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {

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

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	if product.ID == "" {
		errorResponse(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	if err := store.SaveProduct(product); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to save product")
		return
	}

	response := models.MutationResponse{
		Message: "Product created successfully",
		Product: &product,
	}

	jsonResponse(w, http.StatusCreated, response)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	product.ID = id

	if err := store.UpdateProduct(product); err != nil {
		if err.Error() == "product not found" {
			errorResponse(w, http.StatusNotFound, "Product not found")
		} else {
			errorResponse(w, http.StatusInternalServerError, "Failed to update product")
		}
		return
	}

	response := models.MutationResponse{
		Message: "Product updated successfully",
		Product: &product,
	}

	jsonResponse(w, http.StatusOK, response)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := store.DeleteProduct(id); err != nil {
		if err.Error() == "product not found" {
			errorResponse(w, http.StatusNotFound, "Product not found")
		} else {
			errorResponse(w, http.StatusInternalServerError, "Failed to delete product")
		}
		return
	}

	response := models.MutationResponse{
		Message: "Product deleted successfully",
	}

	jsonResponse(w, http.StatusOK, response)
}

// respondWithJSON envía una respuesta exitosa en formato JSON
func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Printf("Error encoding JSON response: %v", err)
	}
}
