# 记账 App 后端技术方案

## 1. 文档信息

- 项目名称：记账 App 后端
- 当前阶段：MVP
- 技术栈：Go、Fiber v3、PostgreSQL、Redis、sqlc、pgx、JWT
- 目标用户：移动端 App、Web 前端、自动化脚本 / 第三方 API 调用方

## 2. 建设目标

本期后端需支持以下核心能力：

- 邮箱注册、登录
- JWT 鉴权
- `access_token` 短期有效，`refresh_token` 长期有效
- 系统默认两级分类
- 创建收入 / 支出记录
- 每条记录带币种
- 查询月统计、日统计、分类统计
- 支持 API Key 调用接口
- 使用 Redis 做缓存与会话加速

## 3. 非目标

本期不纳入范围：

- GitHub / Google OAuth 登录
- 邮箱验证码、找回密码
- 汇率自动换算
- 用户自定义分类
- 预算、账户余额、资产账户体系
- 多租户
- 对外公开开发者平台

## 4. 总体架构

系统采用分层架构，保持与当前仓库一致：

- `handlers`
  - HTTP 请求解析、参数校验、响应格式化
- `services`
  - 业务逻辑、权限校验、缓存协同
- `repository/dbgen`
  - sqlc 生成的数据访问层
- `middleware`
  - JWT 鉴权、API Key 鉴权、错误处理、权限控制
- `config`
  - 环境变量配置
- PostgreSQL
  - 主数据存储，唯一事实源
- Redis
  - 会话缓存、分类缓存、统计缓存、API Key 校验缓存

请求链路如下：

`Client -> Fiber Router -> Middleware -> Handler -> Service -> Repository -> PostgreSQL/Redis`

## 5. 功能范围

### 5.1 用户认证

支持：

- 邮箱注册
- 邮箱登录
- 刷新 access token
- 当前用户信息获取
- 登出

认证策略：

- App/Web 用户端使用 JWT
- access token 用于接口鉴权
- refresh token 用于续期 access token

### 5.2 分类管理

支持系统默认分类读取：

- 支出分类
- 收入分类
- 两级树状结构

本期不支持：

- 用户自定义分类增删改

### 5.3 记账功能

支持：

- 新增收入 / 支出
- 查询账单列表
- 查询账单详情
- 修改账单
- 删除账单
- 按时间范围、分类、币种过滤

### 5.4 统计功能

支持：

- 月统计
- 日统计
- 分类统计

统计粒度：

- 按用户时区进行自然月 / 自然日聚合
- 按币种分组统计
- 本期不做汇率折算

### 5.5 API Key

支持：

- 创建 API Key
- 查询 API Key 列表
- 吊销 API Key
- 使用 API Key 调用部分业务接口

使用场景：

- 自动化脚本
- 第三方服务集成
- 内部工具调用

## 6. 认证与授权设计

### 6.1 双认证体系

系统支持两套认证方式：

1. 用户会话认证
- `Authorization: Bearer <access_token>`

2. 程序调用认证
- `X-API-Key: <api_key>`

业务接口可以复用同一套路由，但鉴权中间件需明确区分认证来源。

### 6.2 Access Token

建议：

- 格式：JWT
- 签名算法：HS256
- 有效期：15 分钟
- 用途：业务接口鉴权

Claims 建议包含：

- `sub`：用户 ID
- `email`
- `token_type=access`
- `exp`
- `iat`

### 6.3 Refresh Token

建议：

- 使用高熵随机字符串，不直接使用 JWT
- 有效期：30 天
- 服务端持久化并支持吊销
- 数据库和 Redis 中只保存摘要或哈希，不保存明文

用途：

- 换取新的 access token
- 支持多设备登录
- 支持服务端强制失效

### 6.4 API Key

建议格式：

- `sk_live_xxxxxxxxxxxxxxxxxxxxxxxxx`

特性：

- 只在创建时展示一次明文
- 数据库存储哈希
- 支持权限范围 `scopes`
- 支持过期时间
- 支持吊销
- 支持最近使用时间记录

### 6.5 权限模型

API Key 使用 scope 控制权限，第一版建议支持：

- `transactions:read`
- `transactions:write`
- `transactions:update`
- `transactions:delete`
- `stats:read`
- `categories:read`

JWT 用户默认按自身身份访问用户资源，不依赖 scope。

## 7. 路由与接口设计

### 7.1 认证接口

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/refresh`
- `POST /api/auth/logout`
- `GET /api/auth/me`

说明：

- `/api/auth/me` 仅允许 JWT
- `/api/auth/refresh` 使用 refresh token
- `/api/auth/logout` 用于吊销当前 refresh token

### 7.2 API Key 管理接口

- `POST /api/api-keys`
- `GET /api/api-keys`
- `POST /api/api-keys/:id/revoke`
- `DELETE /api/api-keys/:id`

说明：

- API Key 管理接口仅允许 JWT
- 创建成功时只返回一次完整 key
- 列表接口仅返回 prefix、name、scope、状态、过期时间、最后使用时间

### 7.3 分类接口

- `GET /api/categories?type=1`

说明：

- `type=1` 为支出
- `type=2` 为收入
- 支持 JWT 和 API Key
- API Key 需要 `categories:read`

### 7.4 记账接口

- `POST /api/transactions`
- `GET /api/transactions`
- `GET /api/transactions/:id`
- `PUT /api/transactions/:id`
- `DELETE /api/transactions/:id`

说明：

- JWT 可访问本人数据
- API Key 可访问所属用户数据
- API Key 需按 scope 控制读写权限

### 7.5 统计接口

- `GET /api/stats/monthly?month=2026-03`
- `GET /api/stats/daily?month=2026-03`
- `GET /api/stats/category?month=2026-03&type=1`

说明：

- 支持 JWT 和 API Key
- API Key 需要 `stats:read`

## 8. 数据库设计

### 8.1 users

用途：保存用户基础信息

字段：

- `id uuid primary key`
- `email varchar(255) unique not null`
- `username varchar(50) not null`
- `password_hash varchar(255) not null`
- `default_currency char(3) not null default 'CNY'`
- `timezone varchar(64) not null default 'Asia/Shanghai'`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`

### 8.2 refresh_tokens

用途：保存 refresh token 会话

字段：

- `id uuid primary key`
- `user_id uuid not null references users(id) on delete cascade`
- `token_hash varchar(255) not null`
- `device_id varchar(128) null`
- `user_agent varchar(255) null`
- `ip_address varchar(64) null`
- `expires_at timestamptz not null`
- `revoked_at timestamptz null`
- `created_at timestamptz not null default now()`

索引建议：

- `idx_refresh_tokens_user_id`
- `idx_refresh_tokens_token_hash`
- `idx_refresh_tokens_expires_at`

### 8.3 api_keys

用途：保存 API Key

字段：

- `id uuid primary key`
- `user_id uuid not null references users(id) on delete cascade`
- `name varchar(100) not null`
- `key_prefix varchar(20) not null`
- `key_hash varchar(255) not null`
- `scopes jsonb not null`
- `status varchar(20) not null default 'active'`
- `last_used_at timestamptz null`
- `expires_at timestamptz null`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`

索引建议：

- `unique(key_prefix)`
- `idx_api_keys_user_id`
- `idx_api_keys_status`

### 8.4 categories

用途：保存系统默认分类

字段：

- `id bigserial primary key`
- `user_id uuid null`
- `type smallint not null`
- `name varchar(50) not null`
- `parent_id bigint null references categories(id)`
- `level smallint not null`
- `sort_order int not null default 0`
- `icon varchar(50) null`
- `is_system boolean not null default false`
- `is_active boolean not null default true`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`

约束：

- 一级分类：`level=1 and parent_id is null`
- 二级分类：`level=2 and parent_id is not null`
- `type in (1, 2)`

### 8.5 transactions

用途：保存记账数据

字段：

- `id bigserial primary key`
- `user_id uuid not null references users(id) on delete cascade`
- `type smallint not null`
- `amount decimal(12,2) not null`
- `currency char(3) not null`
- `category_id bigint not null references categories(id)`
- `description text not null default ''`
- `occurred_at timestamptz not null`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`

约束：

- `type in (1, 2)`
- `amount > 0`
- `currency` 为 3 位大写 ISO 4217 编码

索引建议：

- `idx_transactions_user_occurred_at`
- `idx_transactions_user_type_occurred_at`
- `idx_transactions_user_category_occurred_at`
- `idx_transactions_user_currency_occurred_at`

## 9. 默认分类设计

### 9.1 支出分类

一级：

- 必要
- 非必要
- 其他

二级：

- 必要 -> 吃、穿、住、行
- 非必要 -> 旅行、娱乐、购物
- 其他 -> 其他

### 9.2 收入分类

一级：

- 工资收入
- 其他收入

二级：

- 工资收入 -> 工资、奖金
- 其他收入 -> 理财、兼职、红包、其他

初始化策略：

- 通过初始化 SQL 或启动迁移插入
- 标记为 `is_system=true`
- 本期只读，不允许删除

## 10. Redis 设计

Redis 用作性能加速层，不作为唯一数据源。

### 10.1 Refresh Token 缓存

Key：

- `auth:refresh:{token_hash}`

Value 示例：

```json
{
  "user_id": "uuid",
  "expires_at": "2026-04-30T00:00:00Z",
  "revoked": false
}
```

用途：

- 刷新 token 时优先查 Redis
- 登出时删除 Redis key 并更新数据库

TTL：

- 与 refresh token 剩余有效期一致

### 10.2 API Key 缓存

Key：

- `apikey:{key_prefix}`

Value 示例：

```json
{
  "api_key_id": "uuid",
  "user_id": "uuid",
  "key_hash": "xxx",
  "scopes": ["transactions:write", "stats:read"],
  "status": "active",
  "expires_at": "2026-12-31T23:59:59Z"
}
```

用途：

- API Key 鉴权加速
- 按 prefix 快速定位

TTL：

- 10 到 30 分钟

### 10.3 分类缓存

Key：

- `categories:system:type:{type}`

Value：

- 分类树 JSON

TTL：

- 6 到 24 小时

失效时机：

- 分类初始化变更
- 后续分类管理变更

### 10.4 统计缓存

Key：

- `stats:monthly:{user_id}:{month}`
- `stats:daily:{user_id}:{month}`
- `stats:category:{user_id}:{month}:{type}`

TTL：

- 5 到 15 分钟

失效时机：

- 新增记账
- 修改记账
- 删除记账

## 11. 缓存策略

统一原则：

- PostgreSQL 为事实源
- Redis 为旁路缓存
- 写入数据库成功后删除相关缓存
- 读取时先查 Redis，未命中再回源数据库并回填

不建议缓存：

- 账单明细列表
- access token
- 普通用户基础信息

## 12. 业务规则

### 12.1 注册

规则：

- 邮箱必填且唯一
- 密码长度至少 8 位
- 密码使用 bcrypt 存储
- 成功后自动登录并返回 token

### 12.2 登录

规则：

- 通过邮箱查询用户
- 校验密码
- 返回 access token 与 refresh token

### 12.3 创建账单

规则：

- `type` 只能为收入或支出
- `amount > 0`
- `currency` 必须为合法三位币种码
- `category_id` 必须存在
- 分类必须是二级分类
- 分类 `type` 必须与账单 `type` 一致
- 数据归属当前用户

### 12.4 统计

规则：

- 统计按用户时区进行日期切分
- 按币种独立统计
- 本期不做跨币种汇总折算

### 12.5 API Key

规则：

- 创建时生成一次明文 key
- 服务端只存 hash
- 支持吊销
- 支持过期
- 仅允许访问授予 scope 的接口

## 13. 接口鉴权边界

### 13.1 仅允许 JWT 的接口

- `/api/auth/me`
- `/api/auth/refresh`
- `/api/auth/logout`
- `/api/api-keys`
- `/api/api-keys/:id/revoke`
- `/api/api-keys/:id`

### 13.2 允许 JWT 或 API Key 的接口

- `/api/categories`
- `/api/transactions`
- `/api/stats/monthly`
- `/api/stats/daily`
- `/api/stats/category`

建议实现统一上下文字段：

- `user_id`
- `auth_type`
- `scopes`

这样 service 层不需要区分 JWT 或 API Key。

## 14. 中间件设计

建议新增以下中间件：

- `middleware/auth_jwt.go`
- `middleware/auth_api_key.go`
- `middleware/auth_optional_or.go`
- `middleware/require_scope.go`

职责：

- `auth_jwt`
  - 校验 access token
  - 注入 `user_id`
- `auth_api_key`
  - 校验 `X-API-Key`
  - 校验 scope
  - 注入 `user_id`
- `auth_optional_or`
  - 同一路由支持 JWT 或 API Key
- `require_scope`
  - 检查 API Key 权限范围

## 15. SQL 与查询设计建议

sqlc 建议增加以下查询类别：

- 用户
  - 创建用户
  - 按邮箱查询用户
  - 按 ID 查询用户

- Refresh Token
  - 创建
  - 按 token_hash 查询
  - 吊销
  - 删除过期 token

- API Key
  - 创建
  - 查询用户 API Key 列表
  - 按 prefix 查询
  - 更新 `last_used_at`
  - 吊销

- 分类
  - 按 type 查询系统分类
  - 查询分类详情

- 账单
  - 创建
  - 更新
  - 删除
  - 分页查询
  - 按 ID 查询

- 统计
  - 月统计
  - 日统计
  - 分类统计

## 16. 安全设计

最低要求：

- 密码使用 bcrypt
- refresh token 不存明文
- API Key 不存明文
- JWT 设置短有效期
- 所有查询必须绑定 `user_id`
- 金额使用 decimal，不使用 float
- 时间字段统一使用 `timestamptz`
- 记录关键鉴权失败日志
- API Key 记录 `last_used_at`

后续可扩展：

- 限流
- API Key IP 白名单
- 审计日志
- 邮箱验证
- 第三方登录绑定

## 17. 配置项设计

建议新增环境变量：

- `APP_ENV`
- `HTTP_PORT`
- `JWT_SECRET`
- `ACCESS_TOKEN_TTL_MINUTES`
- `REFRESH_TOKEN_TTL_HOURS`
- `REFRESH_TOKEN_PEPPER`
- `API_KEY_PEPPER`
- `POSTGRES_DSN`
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `REDIS_DB`

可选：

- `DEFAULT_TIMEZONE`
- `DEFAULT_CURRENCY`

## 18. 项目落地方案

结合当前仓库，建议如下改造：

新增或调整文件：

- `sqlc/schema.sql`
- `sqlc/queries.sql`
- `internal/services/auth.go`
- `internal/services/bill.go`
- `internal/services/stats.go`
- `internal/handlers/auth.go`
- `internal/handlers/bill.go`
- 新增 `internal/services/category.go`
- 新增 `internal/services/apikey.go`
- 新增 `internal/handlers/category.go`
- 新增 `internal/handlers/apikey.go`
- 新增 `internal/cache/`
- 新增 `middleware/auth_api_key.go`
- 调整 `internal/apps/router.go`

命名上如果需要保持最小改动，可以继续使用 `bill`，不必立刻全量改名成 `transaction`。

## 19. 分阶段实施计划

### 阶段一：数据库与认证基础

- 增加 `refresh_tokens`
- 增加 `api_keys`
- 增加 `categories`
- 升级 `bills` 表结构支持 `currency` 和 `category_id`
- 完成 register/login/refresh/logout/me

### 阶段二：分类与记账

- 初始化系统默认分类
- 分类树查询
- 创建 / 修改 / 删除 / 查询账单
- 分类与类型校验
- Redis 分类缓存

### 阶段三：统计与缓存

- 月统计
- 日统计
- 分类统计
- Redis 统计缓存
- 账单变更后的统计缓存失效

### 阶段四：API Key

- 创建与吊销 API Key
- API Key 鉴权中间件
- scope 权限控制
- Redis API Key 缓存

## 20. 风险与注意事项

主要风险：

- 当前仓库 `users.password` 与后续多登录方式扩展存在耦合
- 现有 `bills.category int` 需要平滑迁移到正式分类模型
- 时区统计如果处理不当，日统计和月统计会出现边界错误
- 多币种若未来要汇总，需要单独设计汇率体系

本期建议：

- 先保证模型边界清晰
- 不提前做汇率折算
- 不提前做自定义分类
- 不把 API Key 和 JWT 混成同一种会话逻辑

## 21. 结论

本方案适合当前项目的 MVP 阶段，能够在不推翻现有 Fiber + PostgreSQL + sqlc 架构的前提下，补齐：

- 用户会话认证
- 程序调用认证
- 两级分类
- 多币种记账
- 聚合统计
- Redis 缓存能力

下一步可基于此文档继续落地：

1. PostgreSQL 表结构 SQL
2. sqlc 查询定义
3. 接口请求 / 响应 JSON 契约
4. 代码目录改造清单
