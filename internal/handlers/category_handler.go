package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vasujain275/expense-tracker-api/internal/models"
	"github.com/vasujain275/expense-tracker-api/internal/services"
)

type categoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *categoryHandler {
	return &categoryHandler{service: service}
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create a new transaction category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      models.Category  true  "Category object"
// @Success      201  {object}  models.Category
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /categories [post]
func (h *categoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name  string              `json:"name" binding:"required"`
		Type  models.CategoryType `json:"type" binding:"required,oneof=income expense"`
		Color string              `json:"color" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := h.service.CreateCategory(req.Name, req.Type, req.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, category)
}

// GetCategory godoc
// @Summary      Get category by ID
// @Description  Get category details by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /categories/{id} [get]
func (h *categoryHandler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}
	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// GetAllCategories godoc
// @Summary      Get all categories
// @Description  Get all transaction categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Category
// @Failure      500  {object}  ErrorResponse
// @Router       /categories [get]
func (h *categoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetCategoriesByType godoc
// @Summary      Get categories by type
// @Description  Get all categories of a specific type (income/expense)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        type   path      string  true  "Category Type (income/expense)"
// @Success      200  {array}   models.Category
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /categories/type/{type} [get]
func (h *categoryHandler) GetCategoriesByType(c *gin.Context) {
	typeStr := c.Query("type")
	categoryType := models.CategoryType(typeStr)
	categories, err := h.service.GetCategoriesByType(categoryType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// UpdateCategory godoc
// @Summary      Update category
// @Description  Update category information by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Param        category  body      models.Category  true  "Category object"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /categories/{id} [put]
func (h *categoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := h.service.UpdateCategory(id, req.Name, req.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary      Delete category
// @Description  Delete a category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      204  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /categories/{id} [delete]
func (h *categoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}
	if err := h.service.DeleteCategory(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
