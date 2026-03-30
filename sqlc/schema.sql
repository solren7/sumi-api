CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE users (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email             VARCHAR(255) NOT NULL UNIQUE,
    username          VARCHAR(50) NOT NULL,
    password_hash     VARCHAR(255) NOT NULL,
    default_currency  CHAR(3) NOT NULL DEFAULT 'CNY',
    timezone          VARCHAR(64) NOT NULL DEFAULT 'Asia/Shanghai',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_users_modtime
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE refresh_tokens (
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

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

CREATE TABLE api_keys (
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

CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_status ON api_keys(status);

CREATE TRIGGER update_api_keys_modtime
BEFORE UPDATE ON api_keys
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE categories (
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

CREATE INDEX idx_categories_type_parent_sort ON categories(type, parent_id, sort_order);
CREATE INDEX idx_categories_user_type ON categories(user_id, type);

CREATE TRIGGER update_categories_modtime
BEFORE UPDATE ON categories
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE bills (
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

CREATE INDEX idx_bills_user_occurred_at ON bills(user_id, occurred_at DESC);
CREATE INDEX idx_bills_user_type_occurred_at ON bills(user_id, type, occurred_at DESC);
CREATE INDEX idx_bills_user_category_occurred_at ON bills(user_id, category_id, occurred_at DESC);
CREATE INDEX idx_bills_user_currency_occurred_at ON bills(user_id, currency, occurred_at DESC);

CREATE TRIGGER update_bills_modtime
BEFORE UPDATE ON bills
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

INSERT INTO categories (id, type, name, parent_id, level, sort_order, is_system, is_active)
VALUES
    (1001, 1, '必要', NULL, 1, 1, TRUE, TRUE),
    (1002, 1, '非必要', NULL, 1, 2, TRUE, TRUE),
    (1003, 1, '其他', NULL, 1, 3, TRUE, TRUE),
    (1101, 1, '吃', 1001, 2, 1, TRUE, TRUE),
    (1102, 1, '穿', 1001, 2, 2, TRUE, TRUE),
    (1103, 1, '住', 1001, 2, 3, TRUE, TRUE),
    (1104, 1, '行', 1001, 2, 4, TRUE, TRUE),
    (1201, 1, '旅行', 1002, 2, 1, TRUE, TRUE),
    (1202, 1, '娱乐', 1002, 2, 2, TRUE, TRUE),
    (1203, 1, '购物', 1002, 2, 3, TRUE, TRUE),
    (1301, 1, '其他', 1003, 2, 1, TRUE, TRUE),
    (2001, 2, '工资收入', NULL, 1, 1, TRUE, TRUE),
    (2002, 2, '其他收入', NULL, 1, 2, TRUE, TRUE),
    (2101, 2, '工资', 2001, 2, 1, TRUE, TRUE),
    (2102, 2, '奖金', 2001, 2, 2, TRUE, TRUE),
    (2201, 2, '理财', 2002, 2, 1, TRUE, TRUE),
    (2202, 2, '兼职', 2002, 2, 2, TRUE, TRUE),
    (2203, 2, '红包', 2002, 2, 3, TRUE, TRUE),
    (2204, 2, '其他', 2002, 2, 4, TRUE, TRUE)
ON CONFLICT (id) DO NOTHING;

SELECT setval('categories_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM categories), 1), true);
