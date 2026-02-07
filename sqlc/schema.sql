-- 1. 定义自动更新时间的函数 (PostgreSQL 标准做法)
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 2. 用户表
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username    VARCHAR(50) NOT NULL,
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 为 users 表绑定触发器
CREATE TRIGGER update_users_modtime 
BEFORE UPDATE ON users 
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- 3. 账单表
CREATE TABLE bills (
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    amount      DECIMAL(12, 2) NOT NULL,
    
    description TEXT NOT NULL,
    
    -- 1: 支出, 2: 收入
    bill_type   SMALLINT NOT NULL, 
    
    -- 具体分类 ID
    category    INT NOT NULL,
    
    -- 记录时间：精确到时分秒 (例如：2026-02-07 14:30:00+08)
    record_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 为 bills 表绑定触发器
CREATE TRIGGER update_bills_modtime 
BEFORE UPDATE ON bills 
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- 索引优化：因为 record_date 变为了时间戳，查询时通常是范围查询
CREATE INDEX idx_bills_user_record_date ON bills(user_id, record_date);