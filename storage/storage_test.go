package storage

import (
	"os"
	"testing"
)

func TestNewProductStorage(t *testing.T) {
	// Create temp file with test data
	testData := `[
		{
			"id": "TEST001",
			"name": "Test Product",
			"image_url": "https://example.com/test.jpg",
			"description": "Test description",
			"price": 99.99,
			"currency": "USD",
			"rating": 4.5,
			"review_count": 100,
			"specifications": {"color": "red"},
			"category": "Electronics",
			"brand": "TestBrand",
			"in_stock": true
		}
	]`

	tmpFile, err := os.CreateTemp("", "test_products_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(testData)); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Probar creación de storage
	storage, err := NewProduct(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	if storage.GetProductCount() != 1 {
		t.Errorf("Expected 1 product, got %d", storage.GetProductCount())
	}
}

func TestGetProductByID(t *testing.T) {
	testData := `[
		{
			"id": "TEST001",
			"name": "Test Product",
			"image_url": "https://example.com/test.jpg",
			"description": "Test description",
			"price": 99.99,
			"currency": "USD",
			"rating": 4.5,
			"review_count": 100,
			"specifications": {},
			"category": "Electronics",
			"brand": "TestBrand",
			"in_stock": true
		}
	]`

	tmpFile, err := os.CreateTemp("", "test_products_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(testData))
	tmpFile.Close()

	storage, _ := NewProduct(tmpFile.Name())

	// Probar obtener producto existente
	product, err := storage.GetProductById("TEST001")
	if err != nil {
		t.Errorf("Failed to get product: %v", err)
	}

	if product.Name != "Test Product" {
		t.Errorf("Expected product name 'Test Product', got '%s'", product.Name)
	}

	// Probar obtener producto inexistente
	_, err = storage.GetProductById("INVALID")
	if err == nil {
		t.Error("Expected error for invalid product ID, got nil")
	}
}

func TestGetProductsByIDs(t *testing.T) {
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
			"category": "Electronics",
			"brand": "TestBrand",
			"in_stock": false
		}
	]`

	tmpFile, err := os.CreateTemp("", "test_products_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(testData))
	tmpFile.Close()

	storage, _ := NewProduct(tmpFile.Name())

	// Probar obtener múltiples productos válidos
	products, err := storage.GetProductsByIds([]string{"TEST001", "TEST002"})
	if err != nil {
		t.Errorf("Failed to get products: %v", err)
	}

	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}

	// Probar con algunos IDs inválidos
	products, err = storage.GetProductsByIds([]string{"TEST001", "INVALID"})
	if err == nil {
		t.Error("Expected error for partial invalid IDs, got nil")
	}

	if len(products) != 1 {
		t.Errorf("Expected 1 valid product, got %d", len(products))
	}
}

func TestGetProductsByCategory(t *testing.T) {
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
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(testData))
	tmpFile.Close()

	storage, _ := NewProduct(tmpFile.Name())

	// Probar filtrar por categoría existente
	products := storage.GetProductsByCategory("Electronics")
	if len(products) != 1 {
		t.Errorf("Expected 1 product in Electronics, got %d", len(products))
	}

	// Probar filtrar por categoría inexistente
	products = storage.GetProductsByCategory("NonExistent")
	if len(products) != 0 {
		t.Errorf("Expected 0 products in NonExistent category, got %d", len(products))
	}
}

func TestConcurrentAccess(t *testing.T) {
	testData := `[
		{
			"id": "TEST001",
			"name": "Test Product",
			"image_url": "https://example.com/test.jpg",
			"description": "Test description",
			"price": 99.99,
			"currency": "USD",
			"rating": 4.5,
			"review_count": 100,
			"specifications": {},
			"category": "Electronics",
			"brand": "TestBrand",
			"in_stock": true
		}
	]`

	tmpFile, err := os.CreateTemp("", "test_products_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(testData))
	tmpFile.Close()

	storage, _ := NewProduct(tmpFile.Name())

	// Probar acceso concurrente
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, _ = storage.GetProductById("TEST001")
			_ = storage.GetAllProducts()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
