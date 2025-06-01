-- name: CreateCategory :one
INSERT INTO categories (name, type, color)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories 
WHERE id = $1;

-- name: GetCategoryByName :one
SELECT * FROM categories 
WHERE name = $1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY name ASC;

-- name: ListCategoriesByType :many
SELECT * FROM categories
WHERE type = $1
ORDER BY name ASC;

-- name: UpdateCategory :one
UPDATE categories 
SET name = $2, type = $3, color = $4
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories 
WHERE id = $1;

-- name: GetCategoryUsageCount :one
SELECT COUNT(*) as usage_count
FROM transactions 
WHERE category_id = $1;

-- name: GetCategoriesWithTransactionCount :many
SELECT c.*, COUNT(t.id) as transaction_count
FROM categories c
LEFT JOIN transactions t ON c.id = t.category_id
GROUP BY c.id, c.name, c.type, c.color, c.created_at
ORDER BY c.name ASC;
