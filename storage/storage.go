package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dleonardoq/meli/models"
)

func NewProduct(filePath string) (*ProductStorage, error) {
	storage := &ProductStorage{
		filePath: filePath,
		products: make(map[string]models.Product),
	}

	if err := storage.loadProducts(); err != nil {
		return nil, fmt.Errorf("error loading products: %w", err)
	}

	return storage, nil
}

func (s *ProductStorage) loadProducts() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return fmt.Errorf("error unmarshalling data: %w", err)
	}

	for _, product := range products {
		s.products[product.ID] = product
	}

	return nil
}

func (s *ProductStorage) GetProductById(id string) (models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, ok := s.products[id]
	if !ok {
		return models.Product{}, fmt.Errorf("product not found")
	}

	return product, nil
}

func (s *ProductStorage) GetProductsByIds(ids []string) ([]models.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []models.Product
	var notFoundIds []string

	for _, id := range ids {
		product, ok := s.products[id]
		if !ok {
			notFoundIds = append(notFoundIds, id)
			continue
		}
		products = append(products, product)
	}

	if len(notFoundIds) > 0 {
		return products, fmt.Errorf("products not found: %v", notFoundIds)
	}

	return products, nil
}

func (s *ProductStorage) GetAllProducts() []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []models.Product

	for _, product := range s.products {
		products = append(products, product)
	}

	return products
}

func (s *ProductStorage) GetProductsByCategory(category string) []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var products []models.Product

	for _, product := range s.products {
		if product.Category == category {
			products = append(products, product)
		}
	}

	return products
}

func (s *ProductStorage) GetProductCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.products)
}

func (s *ProductStorage) SaveProduct(product models.Product) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.products[product.ID] = product

	return "{message: 'Product saved successfully'}"
}

func (s *ProductStorage) DeleteProduct(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.existsProduct(id) {
		return "", fmt.Errorf("product not found")
	}

	delete(s.products, id)

	return "{message: 'Product deleted successfully'}", nil
}

func (s *ProductStorage) UpdateProduct(product models.Product) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.existsProduct(product.ID) {
		return "", fmt.Errorf("product not found")
	}

	s.products[product.ID] = product

	return "{message: 'Product updated successfully'}", nil
}

func (s *ProductStorage) existsProduct(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.products[id]
	return ok
}
