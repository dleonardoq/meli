package storage

import (
	"sync"

	"github.com/dleonardoq/meli/models"
)

type ProductStorage struct {
	filePath string
	products map[string]models.Product
	mu       sync.RWMutex
}
