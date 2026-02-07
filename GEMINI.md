# GEMINI.md - Project Context

## Project Overview
This project is a Go-based web application built using the **Fiber (v3)** framework. It follows a structured approach for building scalable and maintainable microservices or monolithic APIs.

### Main Technologies
- **Language:** Go (1.25+)
- **Web Framework:** [Fiber v3](https://github.com/gofiber/fiber)
- **Database:** PostgreSQL with [pgx/v5](https://github.com/jackc/pgx)
- **SQL Generator:** [sqlc](https://sqlc.dev/) for type-safe Go code from SQL.
- **Cache:** [Redis](https://github.com/redis/go-redis)
- **CLI Wrapper:** [Cobra](https://github.com/spf13/cobra) for command management.
- **Configuration:** [caarlos0/env](https://github.com/caarlos0/env) and [godotenv](https://github.com/joho/godotenv).
- **Logging:** [zerolog](https://github.com/rs/zerolog) via a custom wrapper in `pkg/logx`.

## Architecture & Directory Structure
The project follows a layered architecture aimed at decoupling concerns:

- `cmd/`: Entry points for the application. Uses Cobra for CLI commands (e.g., `api`).
- `config/`: Configuration loading from environment variables and `.env` files.
- `internal/`: Private application code.
  - `apps/`: Application bootstrapping and graceful shutdown logic.
  - `database/`: Infrastructure layer for PostgreSQL and Redis connection management.
  - `repository/`: Data access layer.
    - `dbgen/`: Generated code by `sqlc`.
  - `services/`, `handlers/`, `entity/`: (Referenced in README) Intended layers for business logic, HTTP handling, and domain entities.
- `middleware/`: Custom Fiber middlewares.
- `pkg/`: Public utility packages (e.g., `logx` for logging).
- `sqlc/`: Contains `schema.sql` and `queries.sql` used by SQLC.

## Building and Running

### Prerequisites
- Go 1.25 or higher
- PostgreSQL and Redis
- `sqlc` installed for code generation

### Key Commands
The project uses a `Makefile` to manage common tasks:

| Command | Description |
| :--- | :--- |
| `make generate` | Runs `sqlc generate` to update database code in `internal/repository/dbgen`. |
| `make run` | Runs the main application. |
| `make api` | Specifically starts the Fiber API server. |
| `make build` | Compiles the application into a binary (`bin/server`). |

## Development Conventions

### Database Workflow
1.  Modify `sqlc/schema.sql` (for schema changes) or `sqlc/queries.sql` (for new queries).
2.  Run `make generate` to update the Go code.
3.  Use the generated functions in the repository layer.

### Configuration
Configuration is managed via environment variables. Create a `.env` file in the root directory based on the following (inferred):
- `SERVER_PORT`: Port for the API server (e.g., `:3000`).
- `DB_DSN`: PostgreSQL connection string.
- `REDIS_ADDR`: Redis address (default: `localhost:6379`).

### Logging
Always use the `pkg/logx` package for logging to ensure consistency and proper trace integration.

### Graceful Shutdown
The application implements graceful shutdown in `internal/apps/app.go`. It listens for `SIGINT` and `SIGTERM` signals to allow active requests to finish and close database connections cleanly.
