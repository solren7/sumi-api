-- +goose Up
CREATE TABLE IF NOT EXISTS configs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type        VARCHAR(50) NOT NULL,
    key         VARCHAR(100) NOT NULL,
    user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
    value       JSONB NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'configs_status_chk'
    ) THEN
        ALTER TABLE configs
        ADD CONSTRAINT configs_status_chk CHECK (status IN ('active', 'inactive'));
    END IF;
END $$;
-- +goose StatementEnd

CREATE UNIQUE INDEX IF NOT EXISTS idx_configs_unique_key
ON configs(type, key, user_id);

CREATE INDEX IF NOT EXISTS idx_configs_type_key ON configs(type, key);
CREATE INDEX IF NOT EXISTS idx_configs_user_type ON configs(user_id, type);

DROP INDEX IF EXISTS idx_categories_type_parent_sort;

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'categories' AND column_name = 'is_system'
    ) THEN
        DELETE FROM categories WHERE user_id IS NULL;
    END IF;
END $$;
-- +goose StatementEnd

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
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM configs
WHERE type = 'category_template' AND key = 'default_categories' AND user_id IS NULL;

DROP INDEX IF EXISTS idx_configs_user_type;
DROP INDEX IF EXISTS idx_configs_type_key;
DROP INDEX IF EXISTS idx_configs_unique_key;

DROP TABLE IF EXISTS configs;
