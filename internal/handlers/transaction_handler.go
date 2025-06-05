package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vasujain275/expense-tracker-api/internal/services"
)

type transactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *transactionHandler {
	return &transactionHandler{service: service}
}

// CreateTransaction godoc
// @Summary      Create a new transaction
// @Description  Create a new transaction and update account balance
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        transaction  body      models.Transaction  true  "Transaction object"
// @Success      201  {object}  models.Transaction
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /transactions [post]
func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var req struct {
		UserID      uuid.UUID       `json:"user_id" binding:"required"`
		AccountID   uuid.UUID       `json:"account_id" binding:"required"`
		CategoryID  uuid.UUID       `json:"category_id" binding:"required"`
		Amount      decimal.Decimal `json:"amount" binding:"required"`
		Description string          `json:"description" binding:"required"`
		Date        string          `json:"date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parsedDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, must be RFC3339"})
		return
	}
	serviceReq := services.TransactionCreateRequest{
		UserID:      req.UserID,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        parsedDate,
	}
	transaction, err := h.service.CreateTransaction(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, transaction)
}

// GetTransaction godoc
// @Summary      Get transaction by ID
// @Description  Get transaction details by its ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /transactions/{id} [get]
func (h *transactionHandler) GetTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}
	transaction, err := h.service.GetTransactionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

// GetTransactions godoc
// @Summary      Get transactions with filters
// @Description  Get transactions with optional filters and pagination
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        user_id     query     string  false  "User ID"
// @Param        account_id  query     string  false  "Account ID"
// @Param        category_id query     string  false  "Category ID"
// @Param        start_date  query     string  false  "Start Date (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "End Date (YYYY-MM-DD)"
// @Param        min_amount  query     number  false  "Minimum Amount"
// @Param        max_amount  query     number  false  "Maximum Amount"
// @Param        limit       query     int     false  "Limit"
// @Param        offset      query     int     false  "Offset"
// @Success      200  {array}   models.Transaction
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /transactions [get]
func (h *transactionHandler) GetTransactions(c *gin.Context) {
	var req struct {
		UserID     uuid.UUID        `form:"user_id" binding:"required"`
		AccountID  *uuid.UUID       `form:"account_id"`
		CategoryID *uuid.UUID       `form:"category_id"`
		StartDate  *string          `form:"start_date"`
		EndDate    *string          `form:"end_date"`
		MinAmount  *decimal.Decimal `form:"min_amount"`
		MaxAmount  *decimal.Decimal `form:"max_amount"`
		Limit      int              `form:"limit,default=20"`
		Offset     int              `form:"offset,default=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, must be YYYY-MM-DD"})
			return
		}
		startDate = &t
	}
	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, must be YYYY-MM-DD"})
			return
		}
		endDate = &t
	}
	serviceReq := services.TransactionListRequest{
		UserID:     req.UserID,
		AccountID:  req.AccountID,
		CategoryID: req.CategoryID,
		StartDate:  startDate,
		EndDate:    endDate,
		MinAmount:  req.MinAmount,
		MaxAmount:  req.MaxAmount,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}
	transactions, count, err := h.service.GetTransactions(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transactions": transactions, "count": count})
}

// UpdateTransaction godoc
// @Summary      Update transaction
// @Description  Update transaction information by its ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Param        transaction  body      models.Transaction  true  "Transaction object"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /transactions/{id} [put]
func (h *transactionHandler) UpdateTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}
	var req struct {
		AccountID   *uuid.UUID       `json:"account_id"`
		CategoryID  *uuid.UUID       `json:"category_id"`
		Amount      *decimal.Decimal `json:"amount"`
		Description *string          `json:"description"`
		Date        *string          `json:"date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var parsedDate *time.Time
	if req.Date != nil {
		t, err := time.Parse(time.RFC3339, *req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, must be RFC3339"})
			return
		}
		parsedDate = &t
	}
	serviceReq := services.TransactionUpdateRequest{
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        parsedDate,
	}
	transaction, err := h.service.UpdateTransaction(id, serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction godoc
// @Summary      Delete transaction
// @Description  Delete a transaction by its ID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      204  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /transactions/{id} [delete]
func (h *transactionHandler) DeleteTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}
	if err := h.service.DeleteTransaction(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetTransactionSummary godoc
// @Summary      Get transaction summary
// @Description  Get spending summary by category for a date range
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        user_id    query     string  true  "User ID"
// @Param        start_date query     string  false  "Start Date (YYYY-MM-DD)"
// @Param        end_date   query     string  false  "End Date (YYYY-MM-DD)"
// @Success      200  {array}   TransactionSummary
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /transactions/summary [get]
func (h *transactionHandler) GetTransactionSummary(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	var startDate, endDate *time.Time
	if s := c.Query("start_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, must be YYYY-MM-DD"})
			return
		}
		startDate = &t
	}
	if s := c.Query("end_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, must be YYYY-MM-DD"})
			return
		}
		endDate = &t
	}
	summary, err := h.service.GetTransactionSummary(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

// GetMonthlyTotal godoc
// @Summary      Get monthly total
// @Description  Get total transactions for a specific month
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        user_id query     string  true  "User ID"
// @Param        year    query     int     true  "Year"
// @Param        month   query     int     true  "Month (1-12)"
// @Success      200  {object}  MonthlyTotalResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /transactions/monthly-total [get]
func (h *transactionHandler) GetMonthlyTotal(c *gin.Context) {
	userIDStr := c.Query("user_id")
	year, yearErr := c.GetQuery("year")
	month, monthErr := c.GetQuery("month")
	userID, err := uuid.Parse(userIDStr)
	if err != nil || !yearErr || !monthErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, year, and month are required"})
		return
	}
	y, err := time.Parse("2006", year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year format, must be YYYY"})
		return
	}
	m, err := time.Parse("01", month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month format, must be MM"})
		return
	}
	// Compose the first day of the month
	monthTime := time.Date(y.Year(), m.Month(), 1, 0, 0, 0, 0, time.UTC)
	total, err := h.service.GetMonthlyTotal(userID, monthTime.Year(), int(monthTime.Month()))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total})
}
