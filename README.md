# Content Backoffice API

REST API in Go for managing AI agent workflows and content automation. Handles agents, workflows, steps, executions, and content reviews.

## Tech Stack

- **Go** 1.24.0
- **Framework:** Gin v1.10.1
- **ORM:** GORM v1.25.7
- **Database:** MySQL (production), SQLite (tests)
- **Auth:** JWT (golang-jwt/jwt v5) + bcrypt
- **Config:** godotenv

## Project Structure

```
.
├── main.go               # Entry point, router setup, AutoMigrate
├── config/               # Env var loading
├── client/               # MySQL GORM connection init
├── handler/              # HTTP route handlers
├── middleware/           # JWT auth + scope-based RBAC
├── model/                # GORM structs
│   └── dto/              # Request/response DTOs
├── service/              # Shared utilities (pagination, filters)
└── scripts/              # Utility scripts (create_admin)
```

## Environment Variables

| Variable     | Default     | Description                  |
|--------------|-------------|------------------------------|
| `DB_USER`    | `nuser`     | MySQL username               |
| `DB_PASS`    | `npass`     | MySQL password               |
| `DB_HOST`    | `localhost` | MySQL host                   |
| `DB_PORT`    | `3306`      | MySQL port                   |
| `DB_NAME`    | `ndb`       | MySQL database name          |
| `JWT_SECRET` | —           | Required. JWT signing secret |

## Running

```bash
# Start the server (port 8080)
go run main.go

# Run tests (uses SQLite in-memory, no MySQL required)
go test ./...
```

## Create Admin User

```bash
go run scripts/create_admin.go <username> <password>
# Example:
go run scripts/create_admin.go admin admin123456
```

Admin users get all scopes: `agents:read/write`, `workflows:read/write`, `steps:read/write`, `step-executions:read/write`, `n:read/write`, `users:admin`.

## Authentication

All protected routes require a Bearer JWT token in the `Authorization` header.

```
Authorization: Bearer <token>
```

### Login

```
POST /auth/login
```

```json
{
  "username": "admin",
  "password": "admin123456"
}
```

Response:
```json
{
  "token": "<jwt>",
  "expires_at": "...",
  "user_id": "admin",
  "scopes": ["workflows:read", "workflows:write", "..."]
}
```

Tokens are valid for **24 hours**. Scopes are always taken from the database — never from the request body.

### Validate Token

```
POST /auth/validate
Authorization: Bearer <token>
```

## API Endpoints

### Agents — scope: `agents:read` / `agents:write`

| Method | Path          | Scope          | Description             |
|--------|---------------|----------------|-------------------------|
| GET    | `/agents`     | `agents:read`  | List agents (paginated) |
| POST   | `/agents`     | `agents:write` | Create agent            |
| PUT    | `/agents/:id` | `agents:write` | Update agent            |
| DELETE | `/agents/:id` | `agents:write` | Delete agent            |

Query params: `page`, `size`

### Workflows — scope: `workflows:read` / `workflows:write`

| Method | Path                     | Scope             | Description                         |
|--------|--------------------------|-------------------|-------------------------------------|
| GET    | `/workflows`             | `workflows:read`  | List workflows (paginated, filterable) |
| POST   | `/workflows`             | `workflows:write` | Create workflow                     |
| PUT    | `/workflows/:id`         | `workflows:write` | Update workflow                     |
| PATCH  | `/workflows/:id/enabled` | `workflows:write` | Enable/disable workflow             |
| DELETE | `/workflows/:id`         | `workflows:write` | Delete workflow                     |

Query params: `page`, `size`, `enabled` (true/false)

### Steps — scope: `steps:read` / `steps:write`

| Method | Path         | Scope         | Description |
|--------|--------------|---------------|-------------|
| GET    | `/steps`     | `steps:read`  | List steps  |
| POST   | `/steps`     | `steps:write` | Create step |
| PUT    | `/steps/:id` | `steps:write` | Update step |
| DELETE | `/steps/:id` | `steps:write` | Delete step |

### Step Executions — scope: `step-executions:read` / `step-executions:write`

| Method | Path                       | Scope                   | Description                          |
|--------|----------------------------|-------------------------|--------------------------------------|
| GET    | `/step-executions-grouped` | `step-executions:read`  | List executions grouped by execution |
| PUT    | `/step-executions-grouped` | `step-executions:write` | Update a step execution (`?id=`)     |

GET filters: `status`, `name`, `workflowId`, `execution_id`, `stepId`, `from` (RFC3339), `to` (RFC3339), `page`, `pageSize`

### Content Reviews — scope: `content-reviews:read` / `content-reviews:write`

| Method | Path                   | Scope                   | Description                              |
|--------|------------------------|-------------------------|------------------------------------------|
| GET    | `/content-reviews`     | `content-reviews:read`  | List content reviews (paginated, filterable) |
| PUT    | `/content-reviews/:id` | `content-reviews:write` | Update content review                    |

GET filters: `status`, `execution_id`, `category`, `from` (YYYY-MM-DD), `to` (YYYY-MM-DD), `page`, `limit` (max 200), `sort`

### Users — scope: `users:admin`

| Method | Path           | Scope         | Description  |
|--------|----------------|---------------|--------------|
| GET    | `/users`       | `users:admin` | List users   |
| POST   | `/users`       | `users:admin` | Create user  |
| PUT    | `/users/:id`   | `users:admin` | Update user  |
| DELETE | `/users/:id`   | `users:admin` | Delete user  |

## Domain Models

- **User** — username, password hash, scopes (JSON array), is_active
- **Agent** — provider, secret
- **Workflow** — name, description, enabled flag, steps (cascade delete)
- **Step** — order, name, operation type, prompt, belongs to Workflow
- **Execution** — workflow_id, status, timestamps
- **StepExecution** — step_id, execution_id, status, output
- **Content** — title, short_description, message, type, sub_type, category, sub_category, image_url, image_prompt, status, execution_id

Schema is managed via GORM `AutoMigrate` on startup — no migration files.

## CORS

All origins are allowed (`*`). Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH. Headers: `Authorization`, `Content-Type`.
