-- name: UpdateAccountBalanceFromTransactions :exec
UPDATE accounts 
SET balance = (
    SELECT COALESCE(SUM(t.amount), 0)
    FROM transactions t
    WHERE t.account_id = accounts.id
), updated_at = NOW()
WHERE id = $1;

-- name: RecalculateAllAccountBalances :exec
UPDATE accounts 
SET balance = subquery.total_amount, updated_at = NOW()
FROM (
    SELECT a.id, COALESCE(SUM(t.amount), 0) as total_amount
    FROM accounts a
    LEFT JOIN transactions t ON a.id = t.account_id
    WHERE a.user_id = $1
    GROUP BY a.id
) subquery
WHERE accounts.id = subquery.id;
