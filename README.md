
## Backend README Completo

```markdown:h%3A%5CSTOCK-
ANALYZER%5CBackend%5CREADME.m
d
# Stock Analyzer Backend API

Una API REST robusta basada 
en Go para análisis de 
acciones y recomendaciones de 
inversión construida con 
framework Gin, PostgreSQL y 
un patrón de arquitectura 
limpia. Este backend 
proporciona gestión completa 
de datos de acciones, 
sincronización en tiempo real 
y endpoints de API seguros.

## 🚀 Características 
Implementadas

- **API RESTful**: Endpoints 
limpios y bien documentados 
con documentación Swagger
- **Arquitectura Limpia**: 
Diseño modular con separación 
de responsabilidades
- **Integración de Base de 
Datos**: PostgreSQL con 
migraciones automatizadas
- **Autenticación JWT**: 
Sistema de autenticación 
seguro basado en tokens
- **Integración de API 
Externa**: Sincronización de 
datos de acciones en tiempo 
real
- **Soporte CORS**: Cross-
origin resource sharing 
habilitado
- **Health Checks**: 
Endpoints de monitoreo de 
salud integrados
- **Documentación Swagger**: 
Documentación interactiva de 
API
- **Testing Comprensivo**: 
Tests unitarios con 
capacidades de mocking
- **Filtrado de Datos**: 
Soporte avanzado de filtrado 
y paginación

## 🛠️ Tecnologías Utilizadas

### Framework Principal
- **Go** (v1.24.4): Lenguaje 
de programación de alto 
rendimiento
- **Gin** (v1.10.1): 
Framework web HTTP para Go
- **PostgreSQL**: Base de 
datos relacional robusta

### Autenticación y Seguridad
- **JWT** (github.com/golang-
jwt/jwt): 
Implementación JSON Web Token
- **CORS** (github.com/gin-
contrib/cors v1.7.5): Cross-
Origin Resource Sharing

### Base de Datos
- **lib/pq** (v1.10.9): 
Driver PostgreSQL para Go
- **Migraciones de Base de 
Datos**: Gestión automatizada 
de esquemas

### Documentación
- **Swagger/OpenAPI**: 
Generación de documentación 
de API
- **gin-swagger** (v1.6.0): 
Integración Gin para Swagger
- **swaggo/swag** (v1.16.4): 
Generador de documentación 
Swagger
- **swaggo/files** (v1.0.1): 
Servicio de archivos 
estáticos para Swagger UI


### Configuración
- **godotenv**: 
Carga de variables de entorno

## 📁 Estructura del 
Proyecto Implementada

```
Backend/
├── main.go                      # Punto de entrada de la aplicación
├── go.mod                       # Dependencias del módulo Go
├── go.sum                       # Checksums de dependencias
├── docs/                        # Documentación Swagger generada
│   ├── docs.go                  # Documentación generada
│   ├── swagger.json             # Especificación OpenAPI JSON
│   └── swagger.yaml             # Especificación OpenAPI YAML
├── internal/
│   ├── api/
│   │   └── routes.go            # Definiciones de rutas HTTP
│   ├── config/
│   │   └── config.go            # Configuración de la aplicación
│   ├── database/
│   │   └── database.go          # Conexión y migraciones de BD
│   ├── entity/
│   │   └── jwt.go               # Estructuras relacionadas con JWT
│   ├── middleware/
│   │   └── jwt.go               # Middleware de autenticación JWT
│   ├── models/
│   │   └── stock.go             # Modelos y estructuras de datos
│   ├── services/
│   │   ├── api_client.go        # Cliente de API externa
│   │   └── stock_service.go     # Capa de lógica de negocio
│   └── main_test.go             # Archivo principal de tests
└── config.tf                    # Configuración Terraform

```

## 🔌 Endpoints de API 
```
### Health Check
```
GET /health
```
Descripción : Estado de salud de la API

Response :

```
{
  "status": "ok",
  "message": "API funcionando 
correctamente"
}
```
### Datos de Acciones
```
GET /api/v1/stocks
```
Descripción : Recuperar acciones con filtrado y paginación

Parámetros de Query :

- ticker (string): Símbolo ticker de la acción
- company (string): Nombre de la empresa
- brokerage (string): Firma de corretaje
- action (string): Acción recomendada (buy, sell, hold)
- rating (string): Rating de la acción
- sort_by (string): Campo de ordenamiento
- order (string): Orden de clasificación (asc, desc)
- page (int): Número de página para paginación
- limit (int): Número de elementos por página
- today (string): Filtro para datos de hoy
Response :

```
{
  "items": [
    {
      "id": 1,
      "ticker": "AAPL",
      "company": "Apple 
Inc.",
      "brokerage": "Goldman 
Sachs",
      "action": "buy",
      "rating_from": "A",
      "rating_to": "A+",
      "target_from": 
"150.00",
      "target_to": "180.00",
      "score": 8.5,
      "confidence": 0.85
    }
  ],
  "next_page": 
"/api/v1/stocks?page=2"
}
```
### Recomendaciones
```
GET /api/v1/recommendations
```
Descripción : Obtener recomendaciones de acciones

Response :

```
{
  "recommendations": [
    {
      "ticker": "AAPL",
      "company": "Apple 
Inc.",
      "score": 8.5,
      "reason": "Strong 
quarterly earnings",
      "target_price": 
"180.00",
      "current_rating": "A+",
      "confidence": 0.85
    }
  ]
}
```
### Documentación
```
GET /swagger/*
```
Descripción : Documentación interactiva Swagger
URL : http://localhost:8080/swagger/index.html


## 🛠️ Instalación y Desarrollo
### Prerrequisitos
- Go 1.24.4 o superior
- Base de datos PostgreSQL
- Variables de entorno configuradas
### Configuración
```
# Clonar el repositorio
git clone <repository-url>
cd Backend

# Instalar dependencias
go mod tidy

# Configurar variables de 
entorno
cp .env.example .env
# Editar .env con tu 
configuración

# Generar documentación 
Swagger
swag init

# Ejecutar la aplicación
go run main.go
```
### Variables de Entorno Requeridas
```
DATABASE_URL=postgres://user:
password@localhost/dbname?
sslmode=disable
JWT_SECRET_KEY=your-secret-
key
API_KEY=external-api-key
API_BASE_URL=https://api.exam
ple.com
PORT=8080
ENVIRONMENT=development


```
## 🔧 Configuración Implementada
### Configuración de Base de Datos
- Migraciones automatizadas al inicio
- Pool de conexiones
- Manejo de errores y recuperación
- Queries optimizadas con índices

### Configuración JWT
- Generación segura de tokens
- Expiración configurable
- Claims basados en usuario
- Middleware de validación

## 📈 Características Implementadas
### Sistema de Autenticación
- Autenticación basada en JWT
- Gestión de roles de usuario
- Validación segura de tokens
- Middleware de autorización
### Documentación de API
- Documentación completa Swagger/OpenAPI
- Interfaz de testing interactiva
- Ejemplos de request/response
- Integración de autenticación
### Gestión de Base de Datos
- Migraciones automatizadas de esquema
- Gestión de conexiones
- Optimización de queries
- Manejo robusto de errores


## 📝 Documentación de API
Accede a la documentación interactiva Swagger en: http://localhost:8080/swagger/index.html

La documentación incluye:
- Descripciones completas de endpoints
- Esquemas de request/response
- Códigos de error y respuestas

## 🔒 Seguridad Implementada
- Validación de Input : Sanitización de parámetros de entrada
- CORS Configurado : Orígenes permitidos específicos
- Variables de Entorno : Configuración sensible en variables de entorno
- Manejo de Errores : No exposición de información sensible