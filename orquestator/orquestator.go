package orquestator

import (
	"github.com/dleonardoq/meli/handlers"
	"github.com/gorilla/mux"
)

func Orquestator() *mux.Router {
	router := mux.NewRouter()

	// Products Endpoints
	api := router.PathPrefix("/meli").Subrouter()
	api.HandleFunc("/products", handlers.GetAllProducts).Methods("GET")
	api.HandleFunc("/products/{id}", handlers.GetProductByID).Methods("GET")
	api.HandleFunc("/products/category/{category}", handlers.GetProductsByCategory).Methods("GET")

	// Compare Endpoints
	api.HandleFunc("/products/compare", handlers.CompareProducts).Methods("POST")

	return router
}
