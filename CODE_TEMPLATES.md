# Code Templates for IPL Backend Service

## 📝 Quick Start Templates

These templates follow the established coding standards and can be used to quickly scaffold new features.

---

## 1. Complete Feature Template

### Model (`internal/models/{entity}.go`)
```go
package models

import "time"

// {Entity} represents the {table_name} table
type {Entity} struct {
	ID          uint       `json:"id" gorm:"primarykey" example:"1"`
	DocumentID  *string    `json:"document_id" gorm:"column:document_id" example:"abc123"`
	Name        string     `json:"name" gorm:"column:name" example:"Example Name"`
	Description *string    `json:"description" gorm:"column:description" example:"Description"`
	IsActive    *bool      `json:"is_active" gorm:"column:is_active" example:"true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedByID *int       `json:"created_by_id"`
	UpdatedByID *int       `json:"updated_by_id"`
}

// TableName sets the insert table name for {Entity}
func ({Entity}) TableName() string {
	return "{table_name}"
}
```

### Repository Interface (`internal/repository/{entity}_repository.go`)
```go
package repository

import (
	"ipl-be-svc/internal/models"
	"gorm.io/gorm"
)

// {Entity}Repository defines the interface for {entity} data operations
type {Entity}Repository interface {
	GetByID(id uint) (*models.{Entity}, error)
	GetAll() ([]*models.{Entity}, error)
	Create({entity} *models.{Entity}) error
	Update({entity} *models.{Entity}) error
	Delete(id uint) error
	GetPaginated(page, limit int) ([]*models.{Entity}, int64, error)
}

// {entity}Repository implements {Entity}Repository
type {entity}Repository struct {
	db *gorm.DB
}

// New{Entity}Repository creates a new instance of {Entity}Repository
func New{Entity}Repository(db *gorm.DB) {Entity}Repository {
	return &{entity}Repository{
		db: db,
	}
}

// GetByID retrieves a {entity} record by ID
func (r *{entity}Repository) GetByID(id uint) (*models.{Entity}, error) {
	var {entity} models.{Entity}
	err := r.db.Where("id = ?", id).First(&{entity}).Error
	if err != nil {
		return nil, err
	}
	return &{entity}, nil
}

// GetAll retrieves all {entity} records
func (r *{entity}Repository) GetAll() ([]*models.{Entity}, error) {
	var {entities} []*models.{Entity}
	err := r.db.Find(&{entities}).Error
	return {entities}, err
}

// Create creates a new {entity} record
func (r *{entity}Repository) Create({entity} *models.{Entity}) error {
	return r.db.Create({entity}).Error
}

// Update updates an existing {entity} record
func (r *{entity}Repository) Update({entity} *models.{Entity}) error {
	return r.db.Save({entity}).Error
}

// Delete deletes a {entity} record by ID
func (r *{entity}Repository) Delete(id uint) error {
	return r.db.Delete(&models.{Entity}{}, id).Error
}

// GetPaginated retrieves paginated {entity} records
func (r *{entity}Repository) GetPaginated(page, limit int) ([]*models.{Entity}, int64, error) {
	var {entities} []*models.{Entity}
	var total int64

	offset := (page - 1) * limit

	// Count total records
	if err := r.db.Model(&models.{Entity}{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Limit(limit).Offset(offset).Find(&{entities}).Error
	return {entities}, total, err
}
```

### Service Interface (`internal/service/{entity}_service.go`)
```go
package service

import (
	"fmt"
	"ipl-be-svc/internal/models"
	"ipl-be-svc/internal/repository"
	"ipl-be-svc/pkg/logger"
)

// {Entity}Service defines the interface for {entity} business operations
type {Entity}Service interface {
	Get{Entity}ByID(id uint) (*models.{Entity}, error)
	GetAll{Entities}() ([]*models.{Entity}, error)
	Create{Entity}(req *Create{Entity}Request) (*models.{Entity}, error)
	Update{Entity}(id uint, req *Update{Entity}Request) error
	Delete{Entity}(id uint) error
	GetPaginated{Entities}(page, limit int) ([]*models.{Entity}, int64, error)
}

// Create{Entity}Request represents {entity} creation input
type Create{Entity}Request struct {
	Name        string  `json:"name" binding:"required" example:"Example Name"`
	Description *string `json:"description" example:"Description"`
	IsActive    *bool   `json:"is_active" example:"true"`
}

// Update{Entity}Request represents {entity} update input
type Update{Entity}Request struct {
	Name        *string `json:"name" example:"Updated Name"`
	Description *string `json:"description" example:"Updated description"`
	IsActive    *bool   `json:"is_active" example:"false"`
}

// {entity}Service implements {Entity}Service
type {entity}Service struct {
	{entity}Repo repository.{Entity}Repository
	logger       *logger.Logger
}

// New{Entity}Service creates a new instance of {Entity}Service
func New{Entity}Service({entity}Repo repository.{Entity}Repository, logger *logger.Logger) {Entity}Service {
	return &{entity}Service{
		{entity}Repo: {entity}Repo,
		logger:       logger,
	}
}

// Get{Entity}ByID retrieves a {entity} by ID with business logic
func (s *{entity}Service) Get{Entity}ByID(id uint) (*models.{Entity}, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid {entity} ID")
	}

	{entity}, err := s.{entity}Repo.GetByID(id)
	if err != nil {
		s.logger.WithError(err).WithField("{entity}_id", id).Error("Failed to get {entity}")
		return nil, fmt.Errorf("{entity} not found")
	}

	return {entity}, nil
}

// GetAll{Entities} retrieves all {entities}
func (s *{entity}Service) GetAll{Entities}() ([]*models.{Entity}, error) {
	{entities}, err := s.{entity}Repo.GetAll()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get {entities}")
		return nil, fmt.Errorf("failed to retrieve {entities}")
	}
	return {entities}, nil
}

// Create{Entity} creates a new {entity}
func (s *{entity}Service) Create{Entity}(req *Create{Entity}Request) (*models.{Entity}, error) {
	// Business validation here
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	{entity} := &models.{Entity}{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := s.{entity}Repo.Create({entity}); err != nil {
		s.logger.WithError(err).Error("Failed to create {entity}")
		return nil, fmt.Errorf("failed to create {entity}")
	}

	s.logger.WithField("{entity}_id", {entity}.ID).Info("{Entity} created successfully")
	return {entity}, nil
}

// Update{Entity} updates an existing {entity}
func (s *{entity}Service) Update{Entity}(id uint, req *Update{Entity}Request) error {
	// Get existing {entity}
	{entity}, err := s.{entity}Repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("{entity} not found")
	}

	// Update fields if provided
	if req.Name != nil {
		{entity}.Name = *req.Name
	}
	if req.Description != nil {
		{entity}.Description = req.Description
	}
	if req.IsActive != nil {
		{entity}.IsActive = req.IsActive
	}

	if err := s.{entity}Repo.Update({entity}); err != nil {
		s.logger.WithError(err).Error("Failed to update {entity}")
		return fmt.Errorf("failed to update {entity}")
	}

	s.logger.WithField("{entity}_id", id).Info("{Entity} updated successfully")
	return nil
}

// Delete{Entity} deletes a {entity}
func (s *{entity}Service) Delete{Entity}(id uint) error {
	// Check if exists
	if _, err := s.{entity}Repo.GetByID(id); err != nil {
		return fmt.Errorf("{entity} not found")
	}

	if err := s.{entity}Repo.Delete(id); err != nil {
		s.logger.WithError(err).Error("Failed to delete {entity}")
		return fmt.Errorf("failed to delete {entity}")
	}

	s.logger.WithField("{entity}_id", id).Info("{Entity} deleted successfully")
	return nil
}

// GetPaginated{Entities} retrieves paginated {entities}
func (s *{entity}Service) GetPaginated{Entities}(page, limit int) ([]*models.{Entity}, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	{entities}, total, err := s.{entity}Repo.GetPaginated(page, limit)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get paginated {entities}")
		return nil, 0, fmt.Errorf("failed to retrieve {entities}")
	}

	return {entities}, total, nil
}
```

### Handler (`internal/handler/{entity}_handler.go`)
```go
package handler

import (
	"strconv"
	"ipl-be-svc/internal/service"
	"ipl-be-svc/pkg/logger"
	"ipl-be-svc/pkg/utils"
	"github.com/gin-gonic/gin"
)

// {Entity}Handler handles {entity}-related HTTP requests
type {Entity}Handler struct {
	{entity}Service service.{Entity}Service
	logger          *logger.Logger
}

// New{Entity}Handler creates a new {entity} handler
func New{Entity}Handler({entity}Service service.{Entity}Service, logger *logger.Logger) *{Entity}Handler {
	return &{Entity}Handler{
		{entity}Service: {entity}Service,
		logger:          logger,
	}
}

// Get{Entity} handles GET /api/v1/{entities}/:id
// @Summary Get {entity} by ID
// @Description Get detailed information about a specific {entity}
// @Tags {entities}
// @Accept json
// @Produce json
// @Param id path int true "{Entity} ID"
// @Success 200 {object} utils.APIResponse{data=models.{Entity}} "{Entity} retrieved successfully"
// @Failure 400 {object} utils.APIResponse "Invalid {entity} ID"
// @Failure 404 {object} utils.APIResponse "{Entity} not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/{entities}/{id} [get]
func (h *{Entity}Handler) Get{Entity}(c *gin.Context) {
	// Parse parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id_param", idParam).Error("Invalid {entity} ID")
		utils.BadRequestResponse(c, "Invalid {entity} ID", err)
		return
	}

	// Call service
	{entity}, err := h.{entity}Service.Get{Entity}ByID(uint(id))
	if err != nil {
		h.logger.WithError(err).WithField("{entity}_id", id).Error("Failed to get {entity}")

		if err.Error() == "{entity} not found" {
			utils.NotFoundResponse(c, "{Entity} not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to get {entity}", err)
		return
	}

	h.logger.WithField("{entity}_id", id).Info("{Entity} retrieved successfully")
	utils.SuccessResponse(c, "{Entity} retrieved successfully", {entity})
}

// GetAll{Entities} handles GET /api/v1/{entities}
// @Summary Get all {entities}
// @Description Get list of all {entities}
// @Tags {entities}
// @Accept json
// @Produce json
// @Success 200 {object} utils.APIResponse{data=[]models.{Entity}} "{Entities} retrieved successfully"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/{entities} [get]
func (h *{Entity}Handler) GetAll{Entities}(c *gin.Context) {
	{entities}, err := h.{entity}Service.GetAll{Entities}()
	if err != nil {
		h.logger.WithError(err).Error("Failed to get {entities}")
		utils.InternalServerErrorResponse(c, "Failed to get {entities}", err)
		return
	}

	h.logger.WithField("count", len({entities})).Info("{Entities} retrieved successfully")
	utils.SuccessResponse(c, "{Entities} retrieved successfully", {entities})
}

// Create{Entity} handles POST /api/v1/{entities}
// @Summary Create a new {entity}
// @Description Create a new {entity} with the provided information
// @Tags {entities}
// @Accept json
// @Produce json
// @Param {entity} body service.Create{Entity}Request true "{Entity} data"
// @Success 201 {object} utils.APIResponse{data=models.{Entity}} "{Entity} created successfully"
// @Failure 400 {object} utils.APIResponse "Invalid request data"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/{entities} [post]
func (h *{Entity}Handler) Create{Entity}(c *gin.Context) {
	var req service.Create{Entity}Request

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body")
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	{entity}, err := h.{entity}Service.Create{Entity}(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create {entity}")
		utils.InternalServerErrorResponse(c, "Failed to create {entity}", err)
		return
	}

	h.logger.WithField("{entity}_id", {entity}.ID).Info("{Entity} created successfully")
	utils.CreatedResponse(c, "{Entity} created successfully", {entity})
}

// Update{Entity} handles PUT /api/v1/{entities}/:id
// @Summary Update a {entity}
// @Description Update an existing {entity} with the provided information
// @Tags {entities}
// @Accept json
// @Produce json
// @Param id path int true "{Entity} ID"
// @Param {entity} body service.Update{Entity}Request true "{Entity} data"
// @Success 200 {object} utils.APIResponse "{Entity} updated successfully"
// @Failure 400 {object} utils.APIResponse "Invalid request data"
// @Failure 404 {object} utils.APIResponse "{Entity} not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/{entities}/{id} [put]
func (h *{Entity}Handler) Update{Entity}(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).Error("Invalid {entity} ID")
		utils.BadRequestResponse(c, "Invalid {entity} ID", err)
		return
	}

	var req service.Update{Entity}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body")
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	if err := h.{entity}Service.Update{Entity}(uint(id), &req); err != nil {
		h.logger.WithError(err).Error("Failed to update {entity}")

		if err.Error() == "{entity} not found" {
			utils.NotFoundResponse(c, "{Entity} not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to update {entity}", err)
		return
	}

	h.logger.WithField("{entity}_id", id).Info("{Entity} updated successfully")
	utils.SuccessResponse(c, "{Entity} updated successfully", nil)
}

// Delete{Entity} handles DELETE /api/v1/{entities}/:id
// @Summary Delete a {entity}
// @Description Delete an existing {entity}
// @Tags {entities}
// @Accept json
// @Produce json
// @Param id path int true "{Entity} ID"
// @Success 200 {object} utils.APIResponse "{Entity} deleted successfully"
// @Failure 400 {object} utils.APIResponse "Invalid {entity} ID"
// @Failure 404 {object} utils.APIResponse "{Entity} not found"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/{entities}/{id} [delete]
func (h *{Entity}Handler) Delete{Entity}(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.WithError(err).Error("Invalid {entity} ID")
		utils.BadRequestResponse(c, "Invalid {entity} ID", err)
		return
	}

	if err := h.{entity}Service.Delete{Entity}(uint(id)); err != nil {
		h.logger.WithError(err).Error("Failed to delete {entity}")

		if err.Error() == "{entity} not found" {
			utils.NotFoundResponse(c, "{Entity} not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to delete {entity}", err)
		return
	}

	h.logger.WithField("{entity}_id", id).Info("{Entity} deleted successfully")
	utils.SuccessResponse(c, "{Entity} deleted successfully", nil)
}

// Get{Entities}Paginated handles GET /api/v1/{entities}/paginated
// @Summary Get paginated {entities}
// @Description Get paginated list of {entities}
// @Tags {entities}
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} utils.PaginatedResponse{data=[]models.{Entity}} "{Entities} retrieved successfully"
// @Failure 400 {object} utils.APIResponse "Invalid pagination parameters"
// @Failure 500 {object} utils.APIResponse "Internal server error"
// @Router /api/v1/{entities}/paginated [get]
func (h *{Entity}Handler) Get{Entities}Paginated(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	{entities}, total, err := h.{entity}Service.GetPaginated{Entities}(page, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get paginated {entities}")
		utils.InternalServerErrorResponse(c, "Failed to get {entities}", err)
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": total,
	}).Info("{Entities} retrieved successfully")

	utils.PaginatedSuccessResponse(c, "{Entities} retrieved successfully", {entities}, page, limit, total)
}
```

### Routes Registration (add to `internal/handler/routes.go`)
```go
// Initialize handler
{entity}Handler := New{Entity}Handler({entity}Service, logger)

// {Entity} routes
{entities} := v1.Group("/{entities}")
{
	{entities}.GET("", {entity}Handler.GetAll{Entities})
	{entities}.GET("/:id", {entity}Handler.Get{Entity})
	{entities}.POST("", {entity}Handler.Create{Entity})
	{entities}.PUT("/:id", {entity}Handler.Update{Entity})
	{entities}.DELETE("/:id", {entity}Handler.Delete{Entity})
	{entities}.GET("/paginated", {entity}Handler.Get{Entities}Paginated)
}
```

### Main.go Registration (add to `cmd/server/main.go`)
```go
// Initialize repositories
{entity}Repo := repository.New{Entity}Repository(db.DB)

// Initialize services
{entity}Service := service.New{Entity}Service({entity}Repo, appLogger)

// Pass to routes
handler.SetupRoutes(router, ..., {entity}Service, ..., appLogger)
```

---

## 2. Quick Snippets

### Basic Model
```go
type {Entity} struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Name      string    `json:"name" gorm:"column:name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ({Entity}) TableName() string {
	return "{table_name}"
}
```

### Repository with Raw SQL
```go
func (r *{entity}Repository) CustomQuery(param string) ([]*models.{Entity}, error) {
	var results []*models.{Entity}
	query := `
		SELECT e.*
		FROM {table_name} e
		WHERE e.column = ?
	`
	err := r.db.Raw(query, param).Scan(&results).Error
	return results, err
}
```

### Service with Transaction
```go
func (s *{entity}Service) CreateWithRelations(req *Request) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create main entity
		entity := &models.{Entity}{...}
		if err := tx.Create(entity).Error; err != nil {
			return err
		}
		
		// Create related entities
		relation := &models.Relation{...}
		if err := tx.Create(relation).Error; err != nil {
			return err
		}
		
		return nil
	})
}
```

### Handler with Query Parameters
```go
func (h *{Entity}Handler) Filter{Entities}(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	
	{entities}, err := h.{entity}Service.Filter({status}, category)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to filter {entities}", err)
		return
	}
	
	utils.SuccessResponse(c, "{Entities} retrieved successfully", {entities})
}
```

---

## 3. Common Replacements

When using these templates, replace:
- `{Entity}` → PascalCase entity name (e.g., `Product`, `User`, `Billing`)
- `{entity}` → camelCase entity name (e.g., `product`, `user`, `billing`)
- `{entities}` → plural camelCase (e.g., `products`, `users`, `billings`)
- `{Entities}` → plural PascalCase (e.g., `Products`, `Users`, `Billings`)
- `{table_name}` → database table name (e.g., `products`, `up_users`, `billings`)

---

## 4. VSCode Snippets

Add to `.vscode/go.code-snippets`:

```json
{
	"IPL Model Template": {
		"prefix": "iplmodel",
		"body": [
			"package models",
			"",
			"import \"time\"",
			"",
			"// ${1:Entity} represents the ${2:table_name} table",
			"type ${1:Entity} struct {",
			"\tID        uint      `json:\"id\" gorm:\"primarykey\"`",
			"\t${3:Name}  string    `json:\"${4:name}\" gorm:\"column:${4:name}\"`",
			"\tCreatedAt time.Time `json:\"created_at\"`",
			"\tUpdatedAt time.Time `json:\"updated_at\"`",
			"}",
			"",
			"func (${1:Entity}) TableName() string {",
			"\treturn \"${2:table_name}\"",
			"}"
		]
	},
	"IPL Repository Template": {
		"prefix": "iplrepo",
		"body": [
			"package repository",
			"",
			"import (",
			"\t\"ipl-be-svc/internal/models\"",
			"\t\"gorm.io/gorm\"",
			")",
			"",
			"type ${1:Entity}Repository interface {",
			"\tGetByID(id uint) (*models.${1:Entity}, error)",
			"}",
			"",
			"type ${2:entity}Repository struct {",
			"\tdb *gorm.DB",
			"}",
			"",
			"func New${1:Entity}Repository(db *gorm.DB) ${1:Entity}Repository {",
			"\treturn &${2:entity}Repository{db: db}",
			"}"
		]
	},
	"IPL Service Template": {
		"prefix": "iplservice",
		"body": [
			"package service",
			"",
			"import (",
			"\t\"ipl-be-svc/internal/repository\"",
			"\t\"ipl-be-svc/pkg/logger\"",
			")",
			"",
			"type ${1:Entity}Service interface {",
			"\tGet${1:Entity}ByID(id uint) error",
			"}",
			"",
			"type ${2:entity}Service struct {",
			"\t${2:entity}Repo repository.${1:Entity}Repository",
			"\tlogger *logger.Logger",
			"}",
			"",
			"func New${1:Entity}Service(repo repository.${1:Entity}Repository, log *logger.Logger) ${1:Entity}Service {",
			"\treturn &${2:entity}Service{${2:entity}Repo: repo, logger: log}",
			"}"
		]
	},
	"IPL Handler Template": {
		"prefix": "iplhandler",
		"body": [
			"package handler",
			"",
			"import (",
			"\t\"ipl-be-svc/internal/service\"",
			"\t\"ipl-be-svc/pkg/logger\"",
			"\t\"ipl-be-svc/pkg/utils\"",
			"\t\"github.com/gin-gonic/gin\"",
			")",
			"",
			"type ${1:Entity}Handler struct {",
			"\t${2:entity}Service service.${1:Entity}Service",
			"\tlogger *logger.Logger",
			"}",
			"",
			"func New${1:Entity}Handler(svc service.${1:Entity}Service, log *logger.Logger) *${1:Entity}Handler {",
			"\treturn &${1:Entity}Handler{${2:entity}Service: svc, logger: log}",
			"}"
		]
	}
}
```

---

**Last Updated**: November 2025
