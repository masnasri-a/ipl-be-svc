# IPL Backend Service

A RESTful API built with Go and Gin framework for menu management system.

## Features

- **Clean Architecture**: Separation of concerns with proper layering
- **PostgreSQL Integration**: Database configuration from MCP
- **GORM ORM**: Database operations with raw SQL support
- **Swagger Documentation**: Interactive API documentation
- **Structured Logging**: JSON/Text logging with logrus
- **Error Handling**: Comprehensive error handling and responses
- **CORS Support**: Cross-origin resource sharing
- **Graceful Shutdown**: Proper application lifecycle management
- **Environment Configuration**: Environment-based configuration

## Project Structure

```
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection and setup
│   ├── models/
│   │   └── models.go            # Data models (MasterMenu)
│   ├── repository/
│   │   └── menu_repository.go   # Menu data access layer
│   ├── service/
│   │   └── menu_service.go      # Menu business logic
│   ├── handler/
│   │   ├── menu_handler.go      # Menu HTTP handlers
│   │   └── routes.go            # Route definitions
│   └── middleware/
│       ├── logger.go            # Logging middleware
│       ├── error.go             # Error handling middleware
│       └── cors.go              # CORS middleware
├── pkg/
│   ├── logger/
│   │   └── logger.go            # Logger utilities
│   └── utils/
│       ├── response.go          # Response utilities
│       └── params.go            # Parameter utilities
├── docs/                        # Swagger documentation
├── .env.example                 # Environment variables example
└── go.mod                       # Go module definition
```

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL database
- Environment variables or .env file

### Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your settings
3. Install dependencies:

```bash
go mod tidy
```

4. Run the application:

```bash
go run cmd/server/main.go
```

### Environment Variables

```env
# Server Configuration
PORT=8080
GIN_MODE=debug

# Database Configuration (from MCP)
DB_HOST=192.168.8.187
DB_PORT=54320
DB_USER=admin
DB_PASSWORD=secret
DB_NAME=strapi
DB_SSLMODE=disable

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
```

## API Endpoints

### Health Check
- `GET /api/v1/health` - Health check endpoint

### Menus
- `GET /api/v1/menus/user/:user_id` - Get menus by user ID

### Swagger Documentation
- `GET /swagger/index.html` - Interactive API documentation

## API Documentation

The API uses the following SQL query to fetch menus by user ID:

```sql
SELECT DISTINCT ON (mm.document_id) mm.*
FROM up_users_role_lnk uurl
INNER JOIN role_menus_role_lnk rmrl ON rmrl.role_id = uurl.role_id
INNER JOIN role_menus_master_menu_lnk rmmml ON rmrl.role_menu_id = rmmml.role_menu_id
INNER JOIN master_menus mm ON rmmml.master_menu_id = mm.id
WHERE uurl.user_id = ?
ORDER BY mm.document_id, mm.id;
```

### Request/Response Examples

#### Get Menus by User ID
```bash
GET /api/v1/menus/user/1
```

**Response:**
```json
{
  "success": true,
  "message": "Menus retrieved successfully",
  "data": [
    {
      "id": 1,
      "document_id": "dashboard",
      "name": "Dashboard",
      "url": "/dashboard",
      "icon": "fas fa-tachometer-alt",
      "order_num": 1,
      "parent_id": null,
      "is_active": true
    }
  ]
}
```

#### Health Check
```bash
GET /api/v1/health
```

**Response:**
```json
{
  "status": "ok",
  "message": "Server is running",
  "service": "IPL Backend Service"
}
```

## Swagger Documentation

After starting the server, visit `http://localhost:8080/swagger/index.html` to access the interactive API documentation.

## Database Schema

The application expects the following database tables:
- `master_menus` - Main menu table
- `up_users_role_lnk` - User-role relationship
- `role_menus_role_lnk` - Role-menu relationship
- `role_menus_master_menu_lnk` - Menu linking table

## Architecture Principles

### Clean Architecture
- **Models**: Define data structures and business entities
- **Repository**: Data access layer with raw SQL queries
- **Service**: Business logic layer with filtering
- **Handler**: HTTP transport layer with Swagger documentation
- **Middleware**: Cross-cutting concerns

### Dependency Injection
- Services depend on repository interfaces
- Handlers depend on service interfaces
- Easy testing and mocking

### Error Handling
- Structured error responses
- Appropriate HTTP status codes
- Detailed logging

## To Run the Application

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Configure environment** (already set up with MCP database):
   ```bash
   cp .env.example .env
   # Edit .env if needed
   ```

3. **Run the server:**
   ```bash
   go run cmd/server/main.go
   ```
   
   Or use VS Code tasks: `Ctrl+Shift+P` → "Tasks: Run Task" → "Run Server"

4. **Test the API:**
   ```bash
   curl http://localhost:8080/api/v1/health
   curl http://localhost:8080/api/v1/menus/user/1
   ```

5. **View Swagger Documentation:**
   Open `http://localhost:8080/swagger/index.html` in your browser

## Features Implemented

- ✅ Clean architecture with menu management
- ✅ PostgreSQL integration with MCP configuration
- ✅ Raw SQL query execution for complex joins
- ✅ Swagger documentation with gin-swagger
- ✅ Health check endpoint
- ✅ Menu retrieval by user ID
- ✅ Structured error handling
- ✅ Environment-based configuration
- ✅ Logging with correlation IDs
- ✅ CORS configuration
- ✅ Graceful shutdown

## License

This project is licensed under the MIT License.