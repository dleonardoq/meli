package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dleonardoq/meli/models"
	"github.com/dleonardoq/meli/storage"
	"github.com/gorilla/mux"
)

// TestAll runs all handler tests in a single master test function
func TestAll(t *testing.T) {
	t.Run("GetAllProducts", TestGetAllProducts)
	t.Run("GetProductByID", TestGetProductByID)
	t.Run("GetProductsByCategory", TestGetProductsByCategory)
	t.Run("CreateProduct", TestCreateProduct)
	t.Run("UpdateProduct", TestUpdateProduct)
	t.Run("DeleteProduct", TestDeleteProduct)
	t.Run("CompareProducts", TestCompareProducts)
}

func setupTestStorage(t *testing.T) (string, func()) {
	testData := `[
		{
			"id": "TEST001",
			"name": "Test Product 1",
			"image_url": "https://example.com/test1.jpg",
			"description": "Test description 1",
			"price": 99.99,
			"currency": "USD",
			"rating": 4.5,
			"review_count": 100,
			"specifications": {},
			"category": "Electronics",
			"brand": "TestBrand",
			"in_stock": true
		},
		{
			"id": "TEST002",
			"name": "Test Product 2",
			"image_url": "https://example.com/test2.jpg",
			"description": "Test description 2",
			"price": 149.99,
			"currency": "USD",
			"rating": 4.8,
			"review_count": 200,
			"specifications": {},
			"category": "Computers",
			"brand": "TestBrand",
			"in_stock": false
		}
	]`

	tmpFile, err := os.CreateTemp("", "test_products_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Initialize test storage
	testStore, err := storage.NewProduct(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Inject test storage into handlers
	SetStorage(testStore)

	cleanup := func() {
		os.Remove(tmpFile.Name())
	}

	return tmpFile.Name(), cleanup
}

func TestGetAllProducts(t *testing.T) {
	_, cleanup := setupTestStorage(t)
	defer cleanup()

	req, err := http.NewRequest("GET", "/meli/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllProducts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.SuccessResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Count != 2 {
		t.Errorf("Expected count 2, got %d", response.Count)
	}

	if len(response.Products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(response.Products))
	}
}

func TestGetProductByID(t *testing.T) {
	_, cleanup := setupTestStorage(t)
	defer cleanup()

	// Test existing product
	req, err := http.NewRequest("GET", "/meli/products/TEST001", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/meli/products/{id}", GetProductByID)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var product models.Product
	if err := json.NewDecoder(rr.Body).Decode(&product); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if product.ID != "TEST001" {
		t.Errorf("Expected product ID TEST001, got %s", product.ID)
	}

	// Test non-existing product
	req, err = http.NewRequest("GET", "/meli/products/INVALID", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestGetProductsByCategory(t *testing.T) {
	_, cleanup := setupTestStorage(t)
	defer cleanup()

	req, err := http.NewRequest("GET", "/meli/products/category/Electronics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/meli/products/category/{category}", GetProductsByCategory)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.SuccessResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Count != 1 {
		t.Errorf("Expected count 1, got %d", response.Count)
	}
}

func TestCreateProduct(t *testing.T) {
	_, cleanup := setupTestStorage(t)
	defer cleanup()

	newProduct := models.Product{
		ID:          "TEST003",
		Name:        "New Test Product",
		ImageURL:    "https://example.com/new.jpg",
		Description: "New test description",
		Price:       199.99,
		Currency:    "USD",
		Rating:      4.7,
		ReviewCount: 50,
		Category:    "Electronics",
		Brand:       "NewBrand",
		InStock:     true,
	}

	body, _ := json.Marshal(newProduct)
	req, err := http.NewRequest("POST", "/meli/products", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateProduct)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response models.MutationResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Message != "Product created successfully" {
		t.Errorf("Expected success message, got %s", response.Message)
	}

	if response.Product == nil || response.Product.ID != "TEST003" {
		t.Error("Expected product in response")
	}
}

func TestUpdateProduct(t *testing.T) {
	_, cleanup := setupTestStorage(t)
	defer cleanup()

	updatedProduct := models.Product{
		ID:          "TEST001",
		Name:        "Updated Product",
		ImageURL:    "https://example.com/updated.jpg",
		Description: "Updated description",
		Price:       149.99,
		Currency:    "USD",
		Rating:      4.8,
		ReviewCount: 150,
		Category:    "Electronics",
		Brand:       "UpdatedBrand",
		InStock:     false,
	}

	body, _ := json.Marshal(updatedProduct)
	req, err := http.NewRequest("PUT", "/meli/products/TEST001", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/meli/products/{id}", UpdateProduct)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.MutationResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Message != "Product updated successfully" {
		t.Errorf("Expected success message, got %s", response.Message)
	}

	// Test updating non-existing product
	req, err = http.NewRequest("PUT", "/meli/products/INVALID", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestDeleteProduct(t *testing.T) {
	_, cleanup := setupTestStorage(t)
	defer cleanup()

	req, err := http.NewRequest("DELETE", "/meli/products/TEST001", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/meli/products/{id}", DeleteProduct)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.MutationResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Message != "Product deleted successfully" {
		t.Errorf("Expected success message, got %s", response.Message)
	}

	// Test deleting non-existing product
	req, err = http.NewRequest("DELETE", "/meli/products/INVALID", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestCompareProducts(t *testing.T) {
	_, cleanup := setupTestStorage(t)
	defer cleanup()

	compareRequest := models.ComparisonRequest{
		ProductIDs: []string{"TEST001", "TEST002"},
	}

	body, _ := json.Marshal(compareRequest)
	req, err := http.NewRequest("POST", "/meli/products/compare", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CompareProducts)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.SuccessResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Count != 2 {
		t.Errorf("Expected count 2, got %d", response.Count)
	}

	// Test with partial invalid IDs
	compareRequest = models.ComparisonRequest{
		ProductIDs: []string{"TEST001", "INVALID"},
	}

	body, _ = json.Marshal(compareRequest)
	req, err = http.NewRequest("POST", "/meli/products/compare", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusPartialContent {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusPartialContent)
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Count != 1 {
		t.Errorf("Expected count 1, got %d", response.Count)
	}
}
