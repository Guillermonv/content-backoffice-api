# Content Backoffice API

## Project Overview
Go REST API for managing AI agent workflows and content automation. Handles agents, workflows, steps, executions, and content reviews.

## Tech Stack
- **Framework:** Gin v1.10.1
- **ORM:** GORM v1.25.7 (MySQL in prod, SQLite in tests)
- **Auth:** JWT (golang-jwt/jwt v5) + bcrypt
- **Config:** godotenv
- **Go:** 1.24.0
- **Server:** `:8080`

## Architecture
`model/` → `handler/` → `middleware/` → `service/`

- **handler/** — HTTP route handlers and request logic
- **model/** — GORM structs (User, Agent, Workflow, Step, Execution, StepExecution, Content)
- **model/dto/** — Request/response Data Transfer Objects
- **middleware/** — JWT auth (`AuthMiddleware`) + scope checks (`RequireScopes()`)
- **service/** — Shared utilities: `Paginate()`, `ApplyStepExecutionFilters()`
- **config/** — Env var loading
- **client/** — MySQL GORM connection init

## Domain Models
- **User** — username, password hash, scopes, is_active
- **Agent** — provider, secret
- **Workflow** — name, description, steps, enabled flag
- **Step** — order, name, operation type, prompt, belongs to Workflow (cascade delete)
- **Execution** — workflow_id, status, timestamps
- **StepExecution** — step_id, execution_id, status, output
- **Content** — title, description, message, type, category, image data

## Auth & Authorization
- Bearer JWT tokens, validated in `AuthMiddleware`
- Scope-based RBAC via `RequireScopes()` (e.g., `"workflows:read"`, `"content-reviews:write"`)
- User scopes stored in JWT claims and DB

## Environment Variables
```
DB_USER=     # default: nuser
DB_PASS=     # default: npass
DB_HOST=     # default: localhost
DB_PORT=     # default: 3306
DB_NAME=     # default: ndb
JWT_SECRET=  # required, no safe default
```

## Running Tests
```bash
go test ./...
```
Tests use SQLite in-memory (`file::memory:?cache=shared`) — no MySQL required.
- `handler/handler_test.go` — integration tests with `setupTestServer()` and `performJSONRequest()`
- `service/pagination_test.go` — unit tests for pagination

## DB Initialization
`AutoMigrate` runs on startup for all models. No migration files — schema is managed via GORM structs.
