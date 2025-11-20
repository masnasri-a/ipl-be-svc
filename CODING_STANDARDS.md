# IPL Backend Service - Coding Standards & Best Practices

## 📋 Table of Contents
1. [Architecture Overview](#architecture-overview)
2. [Project Structure](#project-structure)
3. [Naming Conventions](#naming-conventions)
4. [Layer Responsibilities](#layer-responsibilities)
5. [Code Templates](#code-templates)
6. [Best Practices](#best-practices)
7. [Testing Guidelines](#testing-guidelines)

---

## 🏗️ Architecture Overview

This project follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────┐
│   Handler   │  ← HTTP Layer (Gin handlers, routes, middleware)
└──────┬──────┘
       │
┌──────▼──────┐
│   Service   │  ← Business Logic Layer
└──────┬──────┘
       │
┌──────▼──────┐
│ Repository  │  ← Data Access Layer (GORM, raw SQL)
└──────┬──────┘
       │
┌──────▼──────┐
│  Database   │  ← PostgreSQL
└─────────────┘
```

### Key Principles:
- **Dependency Inversion**: Higher layers depend on interfaces, not implementations
- **Single Responsibility**: Each layer has one clear responsibility
- **Interface Segregation**: Use focused interfaces for each concern
- **Testability**: Easy to mock dependencies for unit testing

---

## 📁 Project Structure

```
ipl-be-svc/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point & dependency injection
├── internal/
│   ├── config/
│   │   └── config.go            # Environment configuration
│   ├── database/
│   │   └── database.go          # Database connection & migrations
│   ├── models/
│   │   ├── base_model.go        # Common base model
│   │   ├── {entity}.go          # Domain entities (1 file per table)
│   │   ├── {entity}_link.go     # Link/junction tables
│   │   └── response/            # Response DTOs
│   │       └── {entity}_response.go
│   ├── repository/
│   │   └── {entity}_repository.go  # Data access per entity
│   ├── service/
│   │   └── {entity}_service.go     # Business logic per entity
│   ├── handler/
│   │   ├── {entity}_handler.go     # HTTP handlers per entity
│   │   └── routes.go               # Route registration
│   └── middleware/
│       ├── cors.go              # CORS configuration
│       ├── logger.go            # Request logging
│       └── error.go             # Error handling
├── pkg/
│   ├── logger/
│   │   └── logger.go            # Logger utilities
│   └── utils/
│       ├── response.go          # Standard response helpers
│       ├── params.go            # Parameter parsing
│       └── jwt.go               # JWT utilities
├── docs/                        # Swagger documentation (auto-generated)
├── .env                         # Environment variables (gitignored)
├── .env.example                 # Environment template
├── go.mod                       # Go module definition
└── README.md                    # Project documentation
```

---

## 🏷️ Naming Conventions

### Files & Directories
- **Snake case** for files: `user_repository.go`, `billing_service.go`
- **Lowercase** for directories: `internal/handler/`, `pkg/utils/`
- **One entity per file**: Separate concerns clearly
- **Link tables**: Use `{entity1}_{entity2}_link.go` pattern

### Go Code
```go
// Interfaces: PascalCase with "er" suffix or descriptive name
type UserService interface {}
type BillingRepository interface {}

// Structs: PascalCase
type User struct {}
type BillingHandler struct {}

// Private implementations: camelCase with interface prefix
type userService struct {}        // implements UserService
type billingRepository struct {}  // implements BillingRepository

// Functions/Methods: PascalCase (public), camelCase (private)
func NewUserService() UserService {}  // Constructor
func (s *userService) getUserByID() {} // Private method
func (s *userService) CreateUser() {}  // Public method

// Variables: camelCase
var userRepo repository.UserRepository
var appLogger *logger.Logger

// Constants: PascalCase or UPPER_SNAKE_CASE
const MaxRetries = 3
const DEFAULT_PAGE_SIZE = 20
```

### Database
- **Tables**: snake_case, plural: `billings`, `master_menus`, `up_users`
- **Columns**: snake_case: `created_at`, `nama_penghuni`, `document_id`
- **Link tables**: `{table1}_{table2}_lnk`: `billings_profile_id_lnk`

### API Endpoints
- **REST conventions**: 
  - `GET /api/v1/users` - List
  - `GET /api/v1/users/:id` - Get one
  - `POST /api/v1/users` - Create
  - `PUT /api/v1/users/:id` - Update
  - `DELETE /api/v1/users/:id` - Delete
- **Nested resources**: `/api/v1/users/:user_id/profile`
- **Actions**: `/api/v1/billings/bulk-monthly` (use verbs only for non-CRUD operations)

---

## 🎯 Layer Responsibilities

### 1. Models Layer (`internal/models/`)

**Purpose**: Define data structures and business entities

```go
package models

import "time"

// User represents the up_users table
type User struct {
	ID        uint       `json:"id" gorm:"primarykey"`
	Username  string     `json:"username" gorm:"column:username"`
	Email     string     `json:"email" gorm:"column:email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	// Relationships
	Roles []Role `json:"roles,omitempty" gorm:"many2many:up_users_role_lnk"`
}

// TableName sets the table name for GORM
func (User) TableName() string {
	return "up_users"
}
```

**Rules**:
- ✅ Pure data structures with tags
- ✅ GORM tags for database mapping: `gorm:"column:field_name"`
- ✅ JSON tags for API responses: `json:"field_name"`
- ✅ Swagger tags for documentation: `example:"value"`
- ✅ One model per database table
- ✅ TableName() method for explicit table naming
- ❌ NO business logic
- ❌ NO database operations
- ❌ NO external dependencies

**Response DTOs** (`internal/models/response/`):
```go
package response

// UserResponse represents API response structure
type UserResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" example:"john@example.com"`
	RoleName string `json:"role_name" example:"Admin"`
}
```

---

### 2. Repository Layer (`internal/repository/`)

**Purpose**: Data access abstraction, database queries

```go
package repository

import (
	"ipl-be-svc/internal/models"
	"gorm.io/gorm"
)

// UserRepository defines data access methods
type UserRepository interface {
	GetByID(id uint) (*models.User, error)
	GetAll() ([]*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	GetByEmail(email string) (*models.User, error)
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Complex query example with raw SQL
func (r *userRepository) GetUserWithRoles(userID uint) (*models.User, error) {
	var user models.User
	query := `
		SELECT u.*, r.name as role_name
		FROM up_users u
		INNER JOIN up_users_role_lnk url ON u.id = url.user_id
		INNER JOIN up_roles r ON url.role_id = r.id
		WHERE u.id = ?
	`
	err := r.db.Raw(query, userID).Scan(&user).Error
	return &user, err
}
```

**Rules**:
- ✅ Interface-based design
- ✅ GORM for simple queries
- ✅ Raw SQL for complex joins
- ✅ Return domain models, not DTOs
- ✅ Handle GORM errors appropriately
- ✅ Use transactions for multi-table operations
- ❌ NO business logic
- ❌ NO HTTP handling
- ❌ NO response formatting

**Common Patterns**:
```go
// Pagination
func (r *userRepository) GetPaginated(page, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64
	
	offset := (page - 1) * limit
	
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

// Filtering
func (r *userRepository) GetByRole(roleType string) ([]*models.User, error) {
	var users []*models.User
	err := r.db.Table("up_users").
		Joins("JOIN up_users_role_lnk url ON up_users.id = url.user_id").
		Joins("JOIN up_roles r ON url.role_id = r.id").
		Where("r.type = ?", roleType).
		Find(&users).Error
	return users, err
}
```

---

### 3. Service Layer (`internal/service/`)

**Purpose**: Business logic, orchestration, validation

```go
package service

import (
	"fmt"
	"ipl-be-svc/internal/models"
	"ipl-be-svc/internal/models/response"
	"ipl-be-svc/internal/repository"
	"ipl-be-svc/pkg/logger"
)

// UserService defines business operations
type UserService interface {
	GetUserByID(id uint) (*response.UserResponse, error)
	CreateUser(req *CreateUserRequest) (*models.User, error)
	UpdateUser(id uint, req *UpdateUserRequest) error
	DeleteUser(id uint) error
	GetActiveUsers() ([]*response.UserResponse, error)
}

// CreateUserRequest represents user creation input
type CreateUserRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
	RoleID   uint   `json:"role_id" binding:"required" example:"1"`
}

// userService implements UserService
type userService struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, logger *logger.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetUserByID retrieves user with business logic
func (s *userService) GetUserByID(id uint) (*response.UserResponse, error) {
	// Validation
	if id == 0 {
		return nil, fmt.Errorf("invalid user ID")
	}
	
	// Repository call
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", id).Error("Failed to get user")
		return nil, fmt.Errorf("user not found")
	}
	
	// Business logic: check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user is inactive")
	}
	
	// Transform to response DTO
	return &response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// CreateUser creates a new user with validation
func (s *userService) CreateUser(req *CreateUserRequest) (*models.User, error) {
	// Business validation
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}
	
	// Hash password (example)
	hashedPassword := hashPassword(req.Password)
	
	// Create model
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}
	
	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		s.logger.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user")
	}
	
	s.logger.WithField("user_id", user.ID).Info("User created successfully")
	return user, nil
}
```

**Rules**:
- ✅ Interface-based design
- ✅ Business logic and validation
- ✅ Transform between DTOs and domain models
- ✅ Logging with context
- ✅ Error handling and wrapping
- ✅ Transaction orchestration
- ❌ NO HTTP handling
- ❌ NO direct database access
- ❌ NO response formatting (use DTOs)

---

### 4. Handler Layer (`internal/handler/`)

**Purpose**: HTTP request/response handling, routing

```go
package handler

import (
	"strconv"
	"ipl-be-svc/internal/service"
	"ipl-be-svc/pkg/logger"
	"ipl-be-svc/pkg/utils"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService service.UserService
	logger      *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetUser handles GET /api/v1/users/:id
// @Summary Get user by ID
// @Description Get detailed information about a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.APIResponse{data=response.UserResponse} "User retrieved successfully"
// @Failure 400 {object} utils.APIResponse "Invalid user ID"
// @Failure 404 {object} utils.APIResponse "User not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	// Parse parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid user ID")
		utils.BadRequestResponse(c, "Invalid user ID", err)
		return
	}
	
	// Call service
	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		h.logger.WithError(err).WithField("user_id", id).Error("Failed to get user")
		
		// Handle specific errors
		if err.Error() == "user not found" {
			utils.NotFoundResponse(c, "User not found")
			return
		}
		
		utils.InternalServerErrorResponse(c, "Failed to get user", err)
		return
	}
	
	// Log success
	h.logger.WithField("user_id", id).Info("User retrieved successfully")
	
	// Return response
	utils.SuccessResponse(c, "User retrieved successfully", user)
}

// CreateUser handles POST /api/v1/users
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body service.CreateUserRequest true "User data"
// @Success 201 {object} utils.APIResponse{data=models.User} "User created successfully"
// @Failure 400 {object} utils.APIResponse "Invalid request data"
// @Failure 409 {object} utils.APIResponse "User already exists"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req service.CreateUserRequest
	
	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body")
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	
	// Call service
	user, err := h.userService.CreateUser(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create user")
		
		if err.Error() == "email already exists" {
			utils.ConflictResponse(c, "User already exists", err)
			return
		}
		
		utils.InternalServerErrorResponse(c, "Failed to create user", err)
		return
	}
	
	h.logger.WithField("user_id", user.ID).Info("User created successfully")
	utils.CreatedResponse(c, "User created successfully", user)
}
```

**Rules**:
- ✅ Swagger annotations for ALL endpoints
- ✅ Use `utils` response helpers
- ✅ Parameter validation and parsing
- ✅ Proper HTTP status codes
- ✅ Logging with context
- ✅ Error categorization
- ❌ NO business logic
- ❌ NO database access
- ❌ NO complex transformations

**Routes Registration** (`internal/handler/routes.go`):
```go
package handler

import (
	"github.com/gin-gonic/gin"
	"ipl-be-svc/internal/service"
	"ipl-be-svc/pkg/logger"
)

func SetupRoutes(
	router *gin.Engine,
	userService service.UserService,
	// ... other services
	logger *logger.Logger,
) {
	// Initialize handlers
	userHandler := NewUserHandler(userService, logger)
	
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", HealthCheck)
		
		// User routes
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}
}
```

---

## 📝 Standard Response Formats

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

### Paginated Response
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Using Response Helpers
```go
import "ipl-be-svc/pkg/utils"

// Success
utils.SuccessResponse(c, "User retrieved successfully", user)

// Created (201)
utils.CreatedResponse(c, "User created successfully", user)

// Bad Request (400)
utils.BadRequestResponse(c, "Invalid input", err)

// Not Found (404)
utils.NotFoundResponse(c, "User not found")

// Internal Server Error (500)
utils.InternalServerErrorResponse(c, "Failed to process request", err)

// Paginated
utils.PaginatedSuccessResponse(c, "Users retrieved", users, page, limit, total)
```

---

## 🔧 Configuration Management

**Environment Variables** (`.env`):
```env
# Server
PORT=8080
GIN_MODE=debug

# Database
DB_HOST=192.168.8.187
DB_PORT=54320
DB_USER=admin
DB_PASSWORD=secret
DB_NAME=strapi
DB_SSLMODE=disable

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json

# JWT
JWT_SECRET=your-secret-key

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

**Config Structure** (`internal/config/config.go`):
```go
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	JWT      JWTConfig
	CORS     CORSConfig
}
```

---

## 🛡️ Error Handling Best Practices

### 1. Repository Layer
```go
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err // Return GORM error as-is
	}
	return &user, nil
}
```

### 2. Service Layer
```go
func (s *userService) GetUserByID(id uint) (*response.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user")
		return nil, fmt.Errorf("user not found") // Wrap with business context
	}
	return transformToResponse(user), nil
}
```

### 3. Handler Layer
```go
func (h *UserHandler) GetUser(c *gin.Context) {
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			utils.NotFoundResponse(c, "User not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to get user", err)
		return
	}
	utils.SuccessResponse(c, "User retrieved successfully", user)
}
```

---

## 📊 Logging Best Practices

```go
// Import logger
import "ipl-be-svc/pkg/logger"

// Initialize in main.go
appLogger := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Format)

// Use in services/handlers
h.logger.Info("Simple message")
h.logger.WithField("user_id", 123).Info("Message with context")
h.logger.WithFields(map[string]interface{}{
	"user_id": 123,
	"action": "create",
}).Info("Message with multiple fields")
h.logger.WithError(err).Error("Error message")
```

**Log Levels**:
- `Debug`: Development debugging
- `Info`: Important events
- `Warn`: Warning conditions
- `Error`: Error conditions
- `Fatal`: Fatal errors (exits application)

---

## 🧪 Testing Guidelines

### Repository Tests
```go
func TestUserRepository_GetByID(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()
	
	repo := repository.NewUserRepository(db)
	
	// Create test data
	user := &models.User{Username: "test", Email: "test@example.com"}
	db.Create(user)
	
	// Test
	result, err := repo.GetByID(user.ID)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user.Username, result.Username)
}
```

### Service Tests (with mocks)
```go
func TestUserService_GetUserByID(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := service.NewUserService(mockRepo, mockLogger)
	
	// Setup expectations
	expectedUser := &models.User{ID: 1, Username: "test"}
	mockRepo.On("GetByID", uint(1)).Return(expectedUser, nil)
	
	// Test
	result, err := service.GetUserByID(1)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Username, result.Username)
	mockRepo.AssertExpectations(t)
}
```

---

## 🚀 Quick Reference: Creating a New Feature

### Step 1: Create Model
```go
// internal/models/product.go
package models

type Product struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"column:name"`
	Price       int64     `json:"price" gorm:"column:price"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Product) TableName() string {
	return "products"
}
```

### Step 2: Create Repository Interface & Implementation
```go
// internal/repository/product_repository.go
package repository

type ProductRepository interface {
	GetByID(id uint) (*models.Product, error)
	GetAll() ([]*models.Product, error)
	Create(product *models.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	return &product, err
}
```

### Step 3: Create Service Interface & Implementation
```go
// internal/service/product_service.go
package service

type ProductService interface {
	GetProductByID(id uint) (*models.Product, error)
}

type productService struct {
	productRepo repository.ProductRepository
	logger      *logger.Logger
}

func NewProductService(repo repository.ProductRepository, log *logger.Logger) ProductService {
	return &productService{productRepo: repo, logger: log}
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	return s.productRepo.GetByID(id)
}
```

### Step 4: Create Handler
```go
// internal/handler/product_handler.go
package handler

type ProductHandler struct {
	productService service.ProductService
	logger         *logger.Logger
}

func NewProductHandler(svc service.ProductService, log *logger.Logger) *ProductHandler {
	return &ProductHandler{productService: svc, logger: log}
}

// @Summary Get product by ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 200 {object} utils.APIResponse{data=models.Product}
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Product not found")
		return
	}
	utils.SuccessResponse(c, "Product retrieved successfully", product)
}
```

### Step 5: Register in main.go
```go
// cmd/server/main.go
// Initialize repository
productRepo := repository.NewProductRepository(db.DB)

// Initialize service
productService := service.NewProductService(productRepo, appLogger)

// Pass to routes
handler.SetupRoutes(router, ..., productService, ..., appLogger)
```

### Step 6: Add Routes
```go
// internal/handler/routes.go
func SetupRoutes(..., productService service.ProductService, ...) {
	productHandler := NewProductHandler(productService, logger)
	
	products := v1.Group("/products")
	{
		products.GET("/:id", productHandler.GetProduct)
	}
}
```

### Step 7: Generate Swagger Docs
```bash
swag init -g cmd/server/main.go -o docs/
```

---

## ✅ Code Review Checklist

- [ ] Follows Clean Architecture layers
- [ ] Interface-based design for services/repositories
- [ ] Proper error handling at each layer
- [ ] Logging with contextual information
- [ ] Swagger documentation for handlers
- [ ] Uses standard response helpers
- [ ] Follows naming conventions
- [ ] No business logic in handlers
- [ ] No HTTP logic in services
- [ ] Repository returns domain models
- [ ] Service transforms to DTOs
- [ ] Proper HTTP status codes
- [ ] Environment variables in config
- [ ] No hardcoded values
- [ ] Transactions for multi-step operations

---

## 📚 Additional Resources

- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Swagger/OpenAPI Spec](https://swagger.io/specification/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**Last Updated**: November 2025
**Version**: 1.0
