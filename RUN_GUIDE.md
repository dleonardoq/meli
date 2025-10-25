# 🚀 Guía de Ejecución

Esta guía explica cómo ejecutar la aplicación Meli API

## 📦 Requisitos Previos

- Go 1.16 o superior instalado
- Puerto 8080 disponible (por defecto) o el que mandes como flag

## 📥 Instalación

1. **Descomprimir el archivo**
   ```bash
   unzip meli-api-v1.0.0.zip -d meli-api
   cd meli-api
   ```

2. **Instalar dependencias**
   ```bash
   go mod download
   ```

   Estructura del proyecto:
   ```
   meli-api/
   ├── data/               # Directorio para almacenamiento de datos
   │   └── products.json   # Archivo de datos de productos
   ├── handlers/           # Manejadores de peticiones HTTP
   ├── middleware/         # Middleware de la aplicación
   ├── orquestator/        # Orquestador de rutas
   ├── models/             # Modelos de datos
   ├── storage/            # Lógica de almacenamiento
   ├── main.go             # Punto de entrada de la aplicación
   ├── go.mod              # Archivo de módulo de Go
   └── README.md           # Documentación del proyecto
   ```

## 🖥️ Ejecución

### Modo Desarrollo
```bash
# Iniciar la aplicación con recarga automática (requiere air)
go install github.com/cosmtrek/air@latest
air

# O iniciar sin recarga automática
go run .
```

### Modo Producción
```bash
# Construir el binario
go build -o meli-api

# Ejecutar el binario
./meli-api
```

### Opciones de Línea de Comandos

```bash
# Cambiar puerto (por defecto: 8080)
go run . -http-port=3000

```

## 🌐 Verificar la Aplicación

Una vez en ejecución, puedes acceder a:
- API: http://localhost:8080/meli/products

## 🔄 Detener la Aplicación
Presiona `Ctrl + C` en la terminal donde se está ejecutando la aplicación.

## 🔄 Reiniciar la Aplicación
1. Detén la aplicación con `Ctrl + C`
2. Vuelve a ejecutarla con los comandos anteriores

## 🔒 Permisos

Asegúrate de que la aplicación tenga permisos de escritura en el directorio donde se ejecuta, especialmente en el directorio `data/`.

## 🚨 Solución de Problemas

Si la aplicación no inicia:
1. Verifica que el puerto 8080 esté disponible
2. Asegúrate de tener permisos de escritura en el directorio
3. Revisa los logs de la aplicación para mensajes de error
4. Verifica que el archivo de datos tenga el formato JSON correcto
