package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vasujain275/expense-tracker-api/internal/models"
	"github.com/vasujain275/expense-tracker-api/internal/services"
)

type accountHandler struct {
	service services.AccountService
}

func NewAccountHandler(service services.AccountService) *accountHandler {
	return &accountHandler{service: service}
}

// CreateAccount godoc
// @Summary      Create a new account
// @Description  Create a new account for a user
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account  body      models.Account  true  "Account object"
// @Success      201  {object}  models.Account
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /accounts [post]
func (h *accountHandler) CreateAccount(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	var req struct {
		Name           string             `json:"name" binding:"required"`
		Type           models.AccountType `json:"type" binding:"required,oneof=bank cash credit_card"`
		InitialBalance decimal.Decimal    `json:"initial_balance" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := h.service.CreateAccount(userID, req.Name, req.Type, req.InitialBalance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, account)
}

// GetAccount godoc
// @Summary      Get account by ID
// @Description  Get account details by its ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  models.Account
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /accounts/{id} [get]
func (h *accountHandler) GetAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}
	account, err := h.service.GetAccountByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

// GetUserAccounts godoc
// @Summary      Get all accounts for a user
// @Description  Get all accounts associated with a user
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        user_id  query     string  true  "User ID"
// @Success      200  {array}   models.Account
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /accounts [get]
func (h *accountHandler) GetUserAccounts(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	activeOnly := c.DefaultQuery("active_only", "false") == "true"
	accounts, err := h.service.GetUserAccounts(userID, activeOnly)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

// UpdateAccount godoc
// @Summary      Update account
// @Description  Update account information by its ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Param        account  body      models.Account  true  "Account object"
// @Success      200  {object}  models.Account
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /accounts/{id} [put]
func (h *accountHandler) UpdateAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}
	var req struct {
		Name     string             `json:"name"`
		Type     models.AccountType `json:"type"`
		IsActive bool               `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := h.service.UpdateAccount(id, req.Name, req.Type, req.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

// DeleteAccount godoc
// @Summary      Delete account
// @Description  Delete an account by its ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      204  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /accounts/{id} [delete]
func (h *accountHandler) DeleteAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}
	if err := h.service.DeleteAccount(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetAccountBalance godoc
// @Summary      Get account balance
// @Description  Get the current balance of an account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  BalanceResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /accounts/{id}/balance [get]
func (h *accountHandler) GetAccountBalance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}
	balance, err := h.service.GetAccountBalance(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
