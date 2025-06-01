-- name: CreateAccount :one
INSERT INTO accounts (user_id, name, type, balance)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts 
WHERE id = $1;

-- name: GetAccountsByUserID :many
SELECT * FROM accounts 
WHERE user_id = $1 AND is_active = true
ORDER BY created_at DESC;

-- name: GetAllAccountsByUserID :many
SELECT * FROM accounts 
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateAccount :one
UPDATE accounts 
SET name = $2, type = $3, balance = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateAccountBalance :one
UPDATE accounts 
SET balance = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteAccount :one
UPDATE accounts 
SET is_active = false, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts 
WHERE id = $1;

-- name: GetAccountWithUser :one
SELECT a.*, u.name as username, u.email as user_email
FROM accounts a
JOIN users u ON a.user_id = u.id
WHERE a.id = $1;

-- name: GetUserTotalBalance :one
SELECT COALESCE(SUM(balance), 0) as total_balance
FROM accounts 
WHERE user_id = $1 AND is_active = true;
