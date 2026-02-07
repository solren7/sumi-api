-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccountsByUserId :many
SELECT * FROM accounts
WHERE user_id = $1 
  AND deleted_at IS NULL
ORDER BY sort_order ASC, created_at DESC;

-- name: CreateAccount :one
INSERT INTO accounts (
  user_id, name, type, currency, icon, color, status, is_default, sort_order
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateAccount :one
UPDATE accounts
SET 
  name = $2,
  type = $3,
  icon = $4,
  color = $5,
  status = $6,
  is_default = $7,
  sort_order = $8,
  updated_at = NOW()
WHERE id = $1
RETURNING *;
-- name: CreateTransaction :one
INSERT INTO transactions (
  user_id, type, amount, currency, description, notes, 
  transaction_date, transaction_id, category_id, account_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetRecentTransactions :many
-- 获取某个用户最近的交易记录，支持分页
SELECT * FROM transactions
WHERE user_id = $1 
  AND deleted_at IS NULL
ORDER BY transaction_date DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: SoftDeleteTransaction :exec
UPDATE transactions
SET deleted_at = NOW()
WHERE id = $1 AND user_id = $2;