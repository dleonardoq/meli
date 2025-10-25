package storage

import (
	"os"
	"testing"

	"github.com/dleonardoq/meli/models"
)

// TestAll runs all storage tests in a single master test function
func TestAll(t *testing.T) {
	t.Run("NewProductStorage", TestNewProductStorage)
	t.Run("GetProductByID", TestGetProductByID)
	t.Run("GetProductsByIDs", TestGetProductsByIDs)
	t.Run("GetProductsByCategory", TestGetProductsByCategory)
	t.Run("GetAllProducts", TestGetAllProducts)
	t.Run("GetProductCount", TestGetProductCount)
	t.Run("SaveProduct", TestSaveProduct)
	t.Run("UpdateProduct", TestUpdateProduct)
	t.Run("DeleteProduct", TestDeleteProduct)
	t.Run("ConcurrentAccess", TestConcurrentAccess)
}

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

func TestGetAllProducts(t *testing.T) {
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

	// Test getting all products
	products := storage.GetAllProducts()
	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}
}

func TestGetProductCount(t *testing.T) {
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

	// Test product count
	count := storage.GetProductCount()
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestSaveProduct(t *testing.T) {
	testData := `[]`

	tmpFile, err := os.CreateTemp("", "test_products_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(testData))
	tmpFile.Close()

	storage, _ := NewProduct(tmpFile.Name())

	// Test saving a new product
	newProduct := models.Product{
		ID:          "TEST001",
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

	err = storage.SaveProduct(newProduct)
	if err != nil {
		t.Errorf("Failed to save product: %v", err)
	}

	// Verify product was saved
	savedProduct, err := storage.GetProductById("TEST001")
	if err != nil {
		t.Errorf("Failed to retrieve saved product: %v", err)
	}

	if savedProduct.Name != "New Test Product" {
		t.Errorf("Expected product name 'New Test Product', got '%s'", savedProduct.Name)
	}

	if storage.GetProductCount() != 1 {
		t.Errorf("Expected count 1 after save, got %d", storage.GetProductCount())
	}
}

func TestUpdateProduct(t *testing.T) {
	testData := `[
		{
			"id": "TEST001",
			"name": "Original Product",
			"image_url": "https://example.com/test.jpg",
			"description": "Original description",
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

	// Test updating existing product
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

	err = storage.UpdateProduct(updatedProduct)
	if err != nil {
		t.Errorf("Failed to update product: %v", err)
	}

	// Verify product was updated
	product, err := storage.GetProductById("TEST001")
	if err != nil {
		t.Errorf("Failed to retrieve updated product: %v", err)
	}

	if product.Name != "Updated Product" {
		t.Errorf("Expected product name 'Updated Product', got '%s'", product.Name)
	}

	if product.Price != 149.99 {
		t.Errorf("Expected price 149.99, got %f", product.Price)
	}

	// Test updating non-existent product
	nonExistentProduct := models.Product{
		ID:   "INVALID",
		Name: "Non-existent",
	}

	err = storage.UpdateProduct(nonExistentProduct)
	if err == nil {
		t.Error("Expected error when updating non-existent product, got nil")
	}
}

func TestDeleteProduct(t *testing.T) {
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

	// Test deleting existing product
	err = storage.DeleteProduct("TEST001")
	if err != nil {
		t.Errorf("Failed to delete product: %v", err)
	}

	// Verify product was deleted
	_, err = storage.GetProductById("TEST001")
	if err == nil {
		t.Error("Expected error when getting deleted product, got nil")
	}

	if storage.GetProductCount() != 1 {
		t.Errorf("Expected count 1 after delete, got %d", storage.GetProductCount())
	}

	// Test deleting non-existent product
	err = storage.DeleteProduct("INVALID")
	if err == nil {
		t.Error("Expected error when deleting non-existent product, got nil")
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
