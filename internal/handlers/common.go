package handlers

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

// BalanceResponse represents an account balance response
type BalanceResponse struct {
	Balance string `json:"balance" example:"1000.00"`
}

// TransactionSummary represents a transaction summary response
type TransactionSummary struct {
	CategoryID   string `json:"category_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CategoryName string `json:"category_name" example:"Groceries"`
	Total        string `json:"total" example:"500.00"`
}

// MonthlyTotalResponse represents a monthly total response
type MonthlyTotalResponse struct {
	Total string `json:"total" example:"1500.00"`
}
