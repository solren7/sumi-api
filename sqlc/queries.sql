-- name: CreateUser :one
INSERT INTO users (
    username, email, password
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateBill :one
INSERT INTO bills (
    user_id, amount, description, bill_type, category, record_date
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateBill :one
-- sqlc 会自动处理 updated_at (通过数据库触发器)，不需要手动传
UPDATE bills
SET 
    amount = $2,
    description = $3,
    bill_type = $4,
    category = $5,
    record_date = $6
WHERE id = $1 AND user_id = $7
RETURNING *;

-- name: DeleteBill :exec
DELETE FROM bills
WHERE id = $1 AND user_id = $2;

-- name: ListBills :many
SELECT * FROM bills
WHERE user_id = $1
ORDER BY record_date DESC
LIMIT $2 OFFSET $3;

-- name: GetMonthlyStats :one
-- 首页顶部统计：依然是按时间范围查，逻辑不变
SELECT 
    COALESCE(SUM(CASE WHEN bill_type = 2 THEN amount ELSE 0 END), 0)::DECIMAL as total_income,
    COALESCE(SUM(CASE WHEN bill_type = 1 THEN amount ELSE 0 END), 0)::DECIMAL as total_expense
FROM bills
WHERE user_id = $1 
  AND record_date >= @start_time::TIMESTAMPTZ 
  AND record_date <= @end_time::TIMESTAMPTZ;

-- name: GetDailyStats :many
-- 首页列表：核心变化在这里
-- 我们需要把 record_date 强转为 DATE 类型来分组
SELECT 
    record_date::DATE as date, -- 强转为日期，去掉时分秒
    COALESCE(SUM(CASE WHEN bill_type = 2 THEN amount ELSE 0 END), 0)::DECIMAL as daily_income,
    COALESCE(SUM(CASE WHEN bill_type = 1 THEN amount ELSE 0 END), 0)::DECIMAL as daily_expense
FROM bills
WHERE user_id = $1 
  AND record_date >= @start_time::TIMESTAMPTZ 
  AND record_date <= @end_time::TIMESTAMPTZ
GROUP BY record_date::DATE
ORDER BY date DESC;