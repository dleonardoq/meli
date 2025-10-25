# Meli Product API

A RESTful API service for managing product information with full CRUD operations, built with Go and Gorilla Mux.

## Features

- **Full CRUD Operations**: Create, Read, Update, and Delete products
- **Category Filtering**: Get products by category
- **Product Comparison**: Compare multiple products at once
- **Thread-Safe**: Concurrent access with RWMutex synchronization
- **In-Memory Storage**: All operations in memory for maximum performance
- **Graceful Shutdown**: Data automatically saved to JSON on server shutdown (Ctrl+C)
- **Persistent Storage**: Data loaded from JSON on startup, saved on shutdown
- **Comprehensive Testing**: Full test coverage for storage and handlers

## Tech Stack

- **Go 1.x**
- **Gorilla Mux** - HTTP router
- **JSON** - Data storage format

## Project Structure

```
meli-test/
├── data/
│   └── products.json          # Product data storage
├── handlers/
│   ├── productHandler.go      # HTTP handlers
│   └── productHandler_test.go # Handler tests
├── models/
│   └── product.go             # Data models
├── orquestator/
│   └── orquestator.go         # Route configuration
├── storage/
│   ├── storage.go             # Storage layer
│   ├── stucts.go              # Storage structs
│   └── storage_test.go        # Storage tests
└── main.go                    # Application entry point
```

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd meli-test
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

The server will start on port `8080` by default.

To specify a custom port:
```bash
go run main.go -http-port=3000
```

### Stopping the Server

To stop the server gracefully and save all data:
```bash
# Press Ctrl+C in the terminal
^C
```

The server will:
1. Catch the shutdown signal
2. Save all in-memory data to `data/products.json`
3. Close all connections gracefully
4. Exit cleanly

**Important:** All changes are kept in memory during execution and only saved to disk when the server stops. This provides maximum performance while ensuring data persistence.

## API Endpoints

### Get All Products
```http
GET /meli/products
```

**Response:**
```json
{
  "products": [...],
  "count": 10
}
```

### Get Product by ID
```http
GET /meli/products/{id}
```

**Response:**
```json
{
  "id": "PROD001",
  "name": "Product Name",
  "price": 99.99,
  ...
}
```

### Get Products by Category
```http
GET /meli/products/category/{category}
```

**Response:**
```json
{
  "products": [...],
  "count": 5
}
```

### Create Product
```http
POST /meli/products
Content-Type: application/json

{
  "id": "PROD001",
  "name": "New Product",
  "image_url": "https://example.com/image.jpg",
  "description": "Product description",
  "price": 99.99,
  "currency": "USD",
  "rating": 4.5,
  "review_count": 100,
  "specifications": {},
  "category": "Electronics",
  "brand": "BrandName",
  "in_stock": true
}
```

**Response:**
```json
{
  "message": "Product created successfully",
  "product": {...}
}
```

### Update Product
```http
PUT /meli/products/{id}
Content-Type: application/json

{
  "name": "Updated Product Name",
  "price": 149.99,
  ...
}
```

**Response:**
```json
{
  "message": "Product updated successfully",
  "product": {...}
}
```

### Delete Product
```http
DELETE /meli/products/{id}
```

**Response:**
```json
{
  "message": "Product deleted successfully"
}
```

### Compare Products
```http
POST /meli/products/compare
Content-Type: application/json

{
  "product_ids": ["PROD001", "PROD002", "PROD003"]
}
```

**Response:**
```json
{
  "products": [...],
  "count": 3
}
```

**Note:** Returns `206 Partial Content` if some products are not found, but at least one is valid.

## Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test ./... -cover
```

Run tests for a specific package:
```bash
go test ./storage -v
go test ./handlers -v
```

## Error Responses

All error responses follow this format:
```json
{
  "error": "Error Type",
  "message": "Detailed error message",
  "code": 400
}
```

### Common HTTP Status Codes

- `200 OK` - Successful GET, PUT, DELETE
- `201 Created` - Successful POST
- `206 Partial Content` - Partial success in comparison
- `400 Bad Request` - Invalid request body or parameters
- `404 Not Found` - Product not found
- `500 Internal Server Error` - Server error

## Data Model

### Product
```go
type Product struct {
    ID             string         `json:"id"`
    Name           string         `json:"name"`
    ImageURL       string         `json:"image_url"`
    Description    string         `json:"description"`
    Price          float64        `json:"price"`
    Currency       string         `json:"currency"`
    Rating         float64        `json:"rating"`
    ReviewCount    int            `json:"review_count"`
    Specifications map[string]any `json:"specifications"`
    Category       string         `json:"category"`
    Brand          string         `json:"brand"`
    InStock        bool           `json:"in_stock"`
}
```

## Development

### Adding New Features

1. Update models in `models/product.go`
2. Add storage methods in `storage/storage.go`
3. Add tests in `storage/storage_test.go`
4. Create handlers in `handlers/productHandler.go`
5. Add handler tests in `handlers/productHandler_test.go`
6. Register routes in `orquestator/orquestator.go`

### Code Quality

- All storage operations are thread-safe using `sync.RWMutex`
- In-memory operations for maximum performance
- Graceful shutdown with data persistence
- Comprehensive error handling
- Full test coverage

### Architecture

**Storage Strategy:**
- **Startup**: Load `data/products.json` into memory
- **Runtime**: All CRUD operations happen in memory (fast)
- **Shutdown**: Save all data back to `data/products.json` (graceful)

**Benefits:**
- ⚡ Ultra-fast operations (no disk I/O during requests)
- 💾 Reduced disk wear (single write on shutdown vs. write per operation)
- 🔒 Thread-safe concurrent access
- 🛡️ Data persistence guaranteed on clean shutdown

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new features
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License
