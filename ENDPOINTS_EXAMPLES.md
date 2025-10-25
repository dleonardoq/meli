# 📚 Ejemplos de Endpoints de la API de Productos

## 📌 Base URL
```
http://localhost:8080/meli
```

## 📦 Productos

### Obtener todos los productos
```http
GET /meli/products
```

**Respuesta Exitosa (200 OK):**
```json
{
  "products": [
    {
      "id": "PROD001",
      "name": "Samsung Galaxy S23 Ultra",
      "image_url": "https://example.com/images/galaxy-s23-ultra.jpg",
      "description": "El último y más avanzado smartphone de Samsung con cámara de 200MP",
      "price": 1199.99,
      "currency": "USD",
      "rating": 4.7,
      "review_count": 1523,
      "specifications": {
        "5g": true,
        "battery": "5000mAh",
        "camera": "108MP + 12MP + 10MP",
        "os": "Android 14",
        "processor": "Snapdragon 8 Gen 2",
        "ram": "12GB",
        "screen_size": "6.7 inches",
        "storage": "256GB",
        "weight": "195g"
      },
      "category": "Electronics",
      "brand": "Samsung",
      "in_stock": true
    }
  ],
  "count": 1
}
```

---

### Obtener un producto por ID
```http
GET /meli/products/PROD001
```

**Respuesta Exitosa (200 OK):**
```json
{
  "id": "PROD001",
  "name": "Samsung Galaxy S23 Ultra",
  "image_url": "https://example.com/images/galaxy-s23-ultra.jpg",
  "description": "El último y más avanzado smartphone de Samsung con cámara de 200MP",
  "price": 1199.99,
  "currency": "USD",
  "rating": 4.7,
  "review_count": 1523,
  "specifications": {
    "5g": true,
    "battery": "5000mAh",
    "camera": "108MP + 12MP + 10MP",
    "os": "Android 14",
    "processor": "Snapdragon 8 Gen 2",
    "ram": "12GB",
    "screen_size": "6.7 inches",
    "storage": "256GB",
    "weight": "195g"
  },
  "category": "Electronics",
  "brand": "Samsung",
  "in_stock": true
}
```

**Producto no encontrado (404 Not Found):**
```json
{
  "error": "Not Found",
  "message": "Product not found",
  "code": 404
}
```

---

### Crear un nuevo producto
```http
POST /meli/products
Content-Type: application/json

{
  "id": "PROD002",
  "name": "iPhone 15 Pro",
  "image_url": "https://example.com/images/iphone15pro.jpg",
  "description": "El último iPhone con chip A17 Pro",
  "price": 999.99,
  "currency": "USD",
  "rating": 4.8,
  "review_count": 0,
  "specifications": {
    "5g": true,
    "battery": "3650mAh",
    "camera": "48MP + 12MP + 12MP",
    "os": "iOS 17",
    "processor": "A17 Bionic",
    "ram": "8GB",
    "screen_size": "6.1 inches",
    "storage": "128GB",
    "weight": "187g"
  },
  "category": "Electronics",
  "brand": "Apple",
  "in_stock": true
}
```

**Respuesta Exitosa (201 Created):**
```json
{
  "message": "Product created successfully",
  "product": {
    "id": "PROD002",
    "name": "iPhone 15 Pro",
    "image_url": "https://example.com/images/iphone15pro.jpg",
    "description": "El último iPhone con chip A17 Pro",
    "price": 999.99,
    "currency": "USD",
    "rating": 4.8,
    "review_count": 0,
    "specifications": {
      "5g": true,
      "battery": "3650mAh",
      "camera": "48MP + 12MP + 12MP",
      "os": "iOS 17",
      "processor": "A17 Bionic",
      "ram": "8GB",
      "screen_size": "6.1 inches",
      "storage": "128GB",
      "weight": "187g"
    },
    "category": "Electronics",
    "brand": "Apple",
    "in_stock": true
  }
}
```

---

### Actualizar un producto existente
```http
PUT /meli/products/PROD002
Content-Type: application/json

{
  "name": "iPhone 15 Pro Max",
  "price": 1199.99,
  "specifications": {
    "screen_size": "6.7 inches",
    "battery": "4422mAh",
    "weight": "221g"
  },
  "in_stock": false
}
```

**Respuesta Exitosa (200 OK):**
```json
{
  "message": "Product updated successfully",
  "product": {
    "id": "PROD002",
    "name": "iPhone 15 Pro Max",
    "image_url": "https://example.com/images/iphone15pro.jpg",
    "description": "El último iPhone con chip A17 Pro",
    "price": 1199.99,
    "currency": "USD",
    "rating": 4.8,
    "review_count": 0,
    "specifications": {
      "5g": true,
      "battery": "4422mAh",
      "camera": "48MP + 12MP + 12MP",
      "os": "iOS 17",
      "processor": "A17 Bionic",
      "ram": "8GB",
      "screen_size": "6.7 inches",
      "storage": "128GB",
      "weight": "221g"
    },
    "category": "Electronics",
    "brand": "Apple",
    "in_stock": false
  }
}
```

---

### Eliminar un producto
```http
DELETE /meli/products/PROD002
```

**Respuesta Exitosa (200 OK):**
```json
{
  "message": "Product deleted successfully"
}
```

---

### Comparar productos
```http
POST /meli/products/compare
Content-Type: application/json

{
  "product_ids": ["PROD001", "NONEXISTENT"]
}
```

**Respuesta Exitosa (206 Partial Content):**
```json
{
  "products": [
    {
      "id": "PROD001",
      "name": "Samsung Galaxy S23 Ultra",
      "image_url": "https://example.com/images/galaxy-s23-ultra.jpg",
      "description": "El último y más avanzado smartphone de Samsung con cámara de 200MP",
      "price": 1199.99,
      "currency": "USD",
      "rating": 4.7,
      "review_count": 1523,
      "specifications": {
        "5g": true,
        "battery": "5000mAh",
        "camera": "108MP + 12MP + 10MP",
        "os": "Android 14",
        "processor": "Snapdragon 8 Gen 2",
        "ram": "12GB",
        "screen_size": "6.7 inches",
        "storage": "256GB",
        "weight": "195g"
      },
      "category": "Electronics",
      "brand": "Samsung",
      "in_stock": true
    }
  ],
  "count": 1
}
```

---

### Obtener productos por categoría
```http
GET /meli/products/category/Electronics
```

**Respuesta Exitosa (200 OK):**
```json
{
  "products": [
    {
      "id": "PROD001",
      "name": "Samsung Galaxy S23 Ultra",
      "image_url": "https://example.com/images/galaxy-s23-ultra.jpg",
      "description": "El último y más avanzado smartphone de Samsung con cámara de 200MP",
      "price": 1199.99,
      "currency": "USD",
      "rating": 4.7,
      "review_count": 1523,
      "specifications": {
        "5g": true,
        "battery": "5000mAh",
        "camera": "108MP + 12MP + 10MP",
        "os": "Android 14",
        "processor": "Snapdragon 8 Gen 2",
        "ram": "12GB",
        "screen_size": "6.7 inches",
        "storage": "256GB",
        "weight": "195g"
      },
      "category": "Electronics",
      "brand": "Samsung",
      "in_stock": true
    }
  ],
  "count": 1
}
```

---

## 🔄 Códigos de Estado HTTP

| Código | Descripción |
|--------|-------------|
| 200 OK | Operación exitosa |
| 201 Created | Recurso creado exitosamente |
| 204 No Content | Operación exitosa sin contenido para devolver |
| 206 Partial Content | Algunos recursos no se encontraron |
| 400 Bad Request | Solicitud incorrecta |
| 404 Not Found | Recurso no encontrado |
| 500 Internal Server Error | Error interno del servidor |

## 🔒 Manejo de Errores

Todos los errores siguen este formato:
```json
{
  "error": "Error Type",
  "message": "Descripción detallada del error",
  "code": 400
}
```

## 🌐 CORS

La API soporta CORS con los siguientes encabezados:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`
