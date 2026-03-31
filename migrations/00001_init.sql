-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS users (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email             VARCHAR(255) NOT NULL UNIQUE,
    username          VARCHAR(50) NOT NULL,
    password_hash     VARCHAR(255) NOT NULL,
    default_currency  CHAR(3) NOT NULL DEFAULT 'CNY',
    timezone          VARCHAR(64) NOT NULL DEFAULT 'Asia/Shanghai',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DROP TRIGGER IF EXISTS update_users_modtime ON users;
CREATE TRIGGER update_users_modtime
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL UNIQUE,
    device_id   VARCHAR(128),
    user_agent  VARCHAR(255),
    ip_address  VARCHAR(64),
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

CREATE TABLE IF NOT EXISTS api_keys (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name          VARCHAR(100) NOT NULL,
    key_prefix    VARCHAR(20) NOT NULL UNIQUE,
    key_hash      VARCHAR(255) NOT NULL,
    scopes        TEXT[] NOT NULL DEFAULT '{}',
    status        VARCHAR(20) NOT NULL DEFAULT 'active',
    last_used_at  TIMESTAMPTZ,
    expires_at    TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_status ON api_keys(status);

DROP TRIGGER IF EXISTS update_api_keys_modtime ON api_keys;
CREATE TRIGGER update_api_keys_modtime
BEFORE UPDATE ON api_keys
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE IF NOT EXISTS configs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type        VARCHAR(50) NOT NULL,
    key         VARCHAR(100) NOT NULL,
    user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
    value       JSONB NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT configs_status_chk CHECK (status IN ('active', 'inactive')),
    CONSTRAINT configs_unique_key UNIQUE (type, key, user_id)
);

CREATE INDEX IF NOT EXISTS idx_configs_type_key ON configs(type, key);
CREATE INDEX IF NOT EXISTS idx_configs_user_type ON configs(user_id, type);

DROP TRIGGER IF EXISTS update_configs_modtime ON configs;
CREATE TRIGGER update_configs_modtime
BEFORE UPDATE ON configs
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE IF NOT EXISTS categories (
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
    type        SMALLINT NOT NULL,
    name        VARCHAR(50) NOT NULL,
    parent_id   BIGINT REFERENCES categories(id) ON DELETE CASCADE,
    level       SMALLINT NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    icon        VARCHAR(50),
    is_system   BOOLEAN NOT NULL DEFAULT FALSE,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT categories_type_chk CHECK (type IN (1, 2)),
    CONSTRAINT categories_level_chk CHECK (level IN (1, 2))
);

CREATE INDEX IF NOT EXISTS idx_categories_type_parent_sort ON categories(type, parent_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_categories_user_type ON categories(user_id, type);

DROP TRIGGER IF EXISTS update_categories_modtime ON categories;
CREATE TRIGGER update_categories_modtime
BEFORE UPDATE ON categories
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE IF NOT EXISTS bills (
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        SMALLINT NOT NULL,
    amount      DECIMAL(12, 2) NOT NULL,
    currency    CHAR(3) NOT NULL,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    description TEXT NOT NULL DEFAULT '',
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT bills_type_chk CHECK (type IN (1, 2)),
    CONSTRAINT bills_amount_chk CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_bills_user_occurred_at ON bills(user_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_bills_user_type_occurred_at ON bills(user_id, type, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_bills_user_category_occurred_at ON bills(user_id, category_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_bills_user_currency_occurred_at ON bills(user_id, currency, occurred_at DESC);

DROP TRIGGER IF EXISTS update_bills_modtime ON bills;
CREATE TRIGGER update_bills_modtime
BEFORE UPDATE ON bills
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO configs (type, key, value, status)
VALUES (
    'category_template',
    'default_categories',
    '{
      "expense": [
        {
          "name": "必要",
          "sort_order": 1,
          "children": [
            { "name": "吃", "sort_order": 1 },
            { "name": "穿", "sort_order": 2 },
            { "name": "住", "sort_order": 3 },
            { "name": "行", "sort_order": 4 }
          ]
        },
        {
          "name": "非必要",
          "sort_order": 2,
          "children": [
            { "name": "旅行", "sort_order": 1 },
            { "name": "娱乐", "sort_order": 2 },
            { "name": "购物", "sort_order": 3 }
          ]
        },
        {
          "name": "其他",
          "sort_order": 3,
          "children": [
            { "name": "其他", "sort_order": 1 }
          ]
        }
      ],
      "income": [
        {
          "name": "工资收入",
          "sort_order": 1,
          "children": [
            { "name": "工资", "sort_order": 1 },
            { "name": "奖金", "sort_order": 2 }
          ]
        },
        {
          "name": "其他收入",
          "sort_order": 2,
          "children": [
            { "name": "理财", "sort_order": 1 },
            { "name": "兼职", "sort_order": 2 },
            { "name": "红包", "sort_order": 3 },
            { "name": "其他", "sort_order": 4 }
          ]
        }
      ]
    }'::jsonb,
    'active'
)
ON CONFLICT (type, key, user_id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS bills;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS configs;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS update_modified_column();
