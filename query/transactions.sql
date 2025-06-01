-- name: CreateTransaction :one
INSERT INTO transactions (account_id, category_id, amount, description, date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transactions 
WHERE id = $1;

-- name: GetTransactionWithDetails :one
SELECT t.*, a.name as account_name, a.type as account_type, c.name as category_name, c.type as category_type, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
WHERE t.id = $1;

-- name: ListTransactions :many
SELECT t.*, a.name as account_name, a.type as account_type, c.name as category_name, c.type as category_type, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
ORDER BY t.date DESC, t.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTransactionsByUserID :many
SELECT t.*, a.name as account_name, a.type as account_type, c.name as category_name, c.type as category_type, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
WHERE a.user_id = $1
ORDER BY t.date DESC, t.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTransactionsByAccountID :many
SELECT t.*, a.name as account_name, c.name as category_name, c.type as category_type, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
WHERE t.account_id = $1
ORDER BY t.date DESC, t.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTransactionsByCategory :many
SELECT t.*, a.name as account_name, a.type as account_type, c.name as category_name, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
WHERE t.category_id = $1
ORDER BY t.date DESC, t.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTransactionsByDateRange :many
SELECT t.*, a.name as account_name, a.type as account_type, c.name as category_name, c.type as category_type, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
WHERE a.user_id = $1 AND t.date >= $2 AND t.date <= $3
ORDER BY t.date DESC, t.created_at DESC
LIMIT $4 OFFSET $5;

-- name: GetTransactionsByAmountRange :many
SELECT t.*, a.name as account_name, a.type as account_type, c.name as category_name, c.type as category_type, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
WHERE a.user_id = $1 AND ABS(t.amount) >= $2 AND ABS(t.amount) <= $3
ORDER BY t.date DESC, t.created_at DESC
LIMIT $4 OFFSET $5;

-- name: SearchTransactionsByDescription :many
SELECT t.*, a.name as account_name, a.type as account_type, c.name as category_name, c.type as category_type, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
WHERE a.user_id = $1 AND t.description ILIKE '%' || $2 || '%'
ORDER BY t.date DESC, t.created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateTransaction :one
UPDATE transactions 
SET account_id = $2, category_id = $3, amount = $4, description = $5, date = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions 
WHERE id = $1;

-- name: GetTransactionCountByUser :one
SELECT COUNT(*) as total_transactions
FROM transactions t
JOIN accounts a ON t.account_id = a.id
WHERE a.user_id = $1;

-- name: GetSpendingSummaryByCategory :many
SELECT c.name as category_name, c.color as category_color, c.type as category_type,
       COUNT(t.id) as transaction_count,
       COALESCE(SUM(ABS(t.amount)), 0) as total_amount
FROM categories c
LEFT JOIN transactions t ON c.id = t.category_id
LEFT JOIN accounts a ON t.account_id = a.id
WHERE (a.user_id = $1 OR a.user_id IS NULL)
  AND ($2::date IS NULL OR t.date >= $2)
  AND ($3::date IS NULL OR t.date <= $3)
GROUP BY c.id, c.name, c.color, c.type
HAVING COUNT(t.id) > 0 OR $4::boolean = true
ORDER BY total_amount DESC;

-- name: GetIncomeVsExpense :one
SELECT 
    COALESCE(SUM(CASE WHEN t.amount > 0 THEN t.amount ELSE 0 END), 0) as total_income,
    COALESCE(SUM(CASE WHEN t.amount < 0 THEN ABS(t.amount) ELSE 0 END), 0) as total_expense,
    COALESCE(SUM(t.amount), 0) as net_amount
FROM transactions t
JOIN accounts a ON t.account_id = a.id
WHERE a.user_id = $1
  AND ($2::date IS NULL OR t.date >= $2)
  AND ($3::date IS NULL OR t.date <= $3);

-- name: GetMonthlySpendingTrend :many
SELECT 
    DATE_TRUNC('month', t.date) as month,
    COALESCE(SUM(CASE WHEN t.amount > 0 THEN t.amount ELSE 0 END), 0) as income,
    COALESCE(SUM(CASE WHEN t.amount < 0 THEN ABS(t.amount) ELSE 0 END), 0) as expense,
    COALESCE(SUM(t.amount), 0) as net
FROM transactions t
JOIN accounts a ON t.account_id = a.id
WHERE a.user_id = $1
  AND t.date >= $2
  AND t.date <= $3
GROUP BY DATE_TRUNC('month', t.date)
ORDER BY month DESC;

-- name: GetRecentTransactions :many
SELECT t.*, a.name as account_name, a.type as account_type, c.name as category_name, c.type as category_type, c.color as category_color
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN categories c ON t.category_id = c.id
WHERE a.user_id = $1
ORDER BY t.created_at DESC
LIMIT $2;
