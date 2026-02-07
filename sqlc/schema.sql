-- ----------------------------
-- Table: users
-- ----------------------------
CREATE TABLE users (
  id         BIGINT PRIMARY KEY,
  email      VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  nickname   VARCHAR(50) NOT NULL,
  avatar_url VARCHAR(500),
  default_currency CHAR(3) NOT NULL DEFAULT 'USD',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

-- ----------------------------
-- Table: accounts
-- ----------------------------
CREATE TABLE accounts (
  id         SERIAL PRIMARY KEY, -- 使用 SERIAL 以自动处理序列
  user_id    BIGINT NOT NULL,
  name       VARCHAR(50) NOT NULL,
  type       VARCHAR(20) NOT NULL,
  currency   CHAR(3) NOT NULL DEFAULT 'USD',
  icon       JSONB NOT NULL DEFAULT '{}',
  color      VARCHAR(7),
  status     VARCHAR(20) NOT NULL DEFAULT 'active',
  is_default BOOLEAN NOT NULL DEFAULT FALSE,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,

  CONSTRAINT accounts_status_check CHECK (status IN ('active', 'inactive', 'archived')),
  CONSTRAINT accounts_type_check CHECK (type IN ('cash', 'bank', 'credit_card', 'alipay', 'wechat', 'investment', 'other'))
);

-- ----------------------------
-- Table: transactions
-- ----------------------------
CREATE TABLE transactions (
  id               SERIAL PRIMARY KEY,
  user_id          BIGINT NOT NULL,
  type             VARCHAR(10) NOT NULL,
  amount           NUMERIC(15,2) NOT NULL,
  currency         CHAR(3) NOT NULL,
  description      VARCHAR(200),
  notes            TEXT,
  transaction_date DATE NOT NULL,
  transaction_id   UUID,
  category_id      BIGINT,
  account_id       BIGINT,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at       TIMESTAMPTZ,

  CONSTRAINT transactions_type_check CHECK (type IN ('income', 'expense', 'transfer')),
  CONSTRAINT transactions_amount_check CHECK (amount > 0)
);