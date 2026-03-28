# CLAUDE.md — GO-Eats Backend

## Project Overview

GO-Eats is a Go-based food delivery backend implementing clean architecture. It serves
RESTful APIs for user management, restaurant operations, cart/order processing, delivery
personnel management (with 2FA), reviews, and real-time notifications via WebSocket and
NATS messaging.

**Current Stage:** MVP / Pre-Production (~40% production ready)
**Goal:** Resolve all 10 critical blockers and deploy live on Railway.

---

## Tech Stack

| Component      | Technology                  |
|----------------|-----------------------------|
| Language       | Go 1.24.4                   |
| Framework      | Gin v1.10.1                 |
| ORM            | Bun v1.2.15                 |
| Database       | PostgreSQL 16               |
| Message Broker | NATS v1.43.0                |
| Real-time      | gorilla/websocket           |
| Auth           | JWT (HS256) + TOTP 2FA      |
| Containerization | Docker Compose            |

---

## Project Structure

```
GO-Eats/
├── cmd/api/
│   ├── main.go              # Entry point, wires all components
│   └── middleware/
│       └── middleware.go    # JWT auth middleware
├── pkg/
│   ├── abstract/            # Domain interfaces (DDD contracts)
│   ├── database/
│   │   ├── database.go      # Generic repository + DB connection
│   │   └── models/          # Bun ORM structs (user, restaurant, cart, order, delivery, review)
│   ├── handler/             # HTTP layer (Gin handlers)
│   │   └── server.go        # Gin setup, CORS, logging
│   ├── service/             # Business logic layer
│   ├── nats/
│   │   └── nats_server.go   # NATS connection and pub/sub
│   ├── storage/             # Image storage abstraction (local only)
│   └── tests/               # Integration tests (testcontainers)
├── httpclient/
│   ├── FoodDelivery.http
│   └── FoodDelivery.postman_collection.json
├── docker-compose.yml
├── go.mod
└── .env
```

---

## Running Locally

```bash
# Start PostgreSQL + NATS
docker-compose up -d

# Run the server (port 8080)
go run cmd/api/main.go

# Run all tests
go test ./... -v

# Health check
curl http://localhost:8080/healthz
```

---

## Environment Variables

These must be set in `.env` locally or injected via Railway in production:

```
APP_ENV=dev
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_NAME=food-delivery
JWT_SECRET_KEY=your_secret
STORAGE_TYPE=local
STORAGE_DIRECTORY=uploads
LOCAL_STORAGE_PATH=/path/to/uploads
```

**Important:** `.env` must NEVER be committed to git. Always use `.env.example` for documentation.

---

## 🔴 Critical Blockers (Must Fix Before Deployment)

These 10 issues block production deployment. Fix them in this order:

| # | Issue | Location |
|---|-------|----------|
| 1 | CORS allows all origins (`*`) | `pkg/handler/server.go:25` |
| 2 | Hardcoded delivery address "New Delhi" | `pkg/service/cart_order/place_order.go:33` |
| 3 | No rate limiting middleware | N/A — needs to be added |
| 4 | No password validation | `pkg/database/models/user/user.go` |
| 5 | No database migration system | `pkg/database/database.go:201` |
| 6 | Missing FK indexes on all foreign key columns | All model files |
| 7 | No DB connection pool configuration | `pkg/database/database.go:185` |
| 8 | Incomplete input validation on handlers | Multiple handler files |
| 9 | No order status state machine | `pkg/service/cart_order/` |
| 10 | JWT secret in `.env` file in repo | `.env:11` |

---

## 🟡 Post-Launch (Do NOT work on these until blockers are resolved)

- File uploads: stay local for now, S3 later
- Pagination on list endpoints
- OpenAPI/Swagger docs
- Redis caching
- Payment integration (Stripe/Razorpay)
- Email notifications
- API versioning (`/api/v1/`)

---

## Deployment Target: Railway

Services to deploy:
1. **Go app** — built from `Dockerfile` (needs to be created)
2. **PostgreSQL** — Railway managed instance
3. **NATS** — Railway service

Railway reads env variables from its dashboard — do not rely on `.env` file in production.

### Pre-deploy checklist
- [ ] `.env` removed from git history
- [ ] `Dockerfile` created (multi-stage build preferred)
- [ ] `golang-migrate` integrated for schema migrations
- [ ] All 10 critical blockers resolved
- [ ] CORS restricted to Railway frontend domain
- [ ] `/healthz` endpoint confirmed working

---

## Architecture Layers (Top to Bottom)

```
Client → Gin HTTP Server → Handler Layer → Service Layer → Repository Layer → PostgreSQL / NATS
```

- **Handlers** — HTTP only, no business logic
- **Services** — All business logic lives here
- **Repository** — Generic DB operations via Bun ORM
- **Models** — Bun ORM structs in `pkg/database/models/`
- **Interfaces** — Defined in `pkg/abstract/` (dependency inversion)

---

## Key Conventions

- All routes go through JWT middleware except `/user/register` and `/user/login`
- Delivery personnel use TOTP 2FA on top of JWT
- NATS is used for order status events (pub/sub)
- WebSocket handles real-time notifications to clients
- SSE handles announcements

---

## Testing

```bash
# All tests (uses testcontainers — Docker must be running)
go test ./... -v

# Specific package
go test ./pkg/tests/user -v

# With race detector
go test ./... -race
```

Integration tests spin up a real PostgreSQL container via testcontainers-go.

---

## What Claude Should NOT Do

- Do not modify the `httpclient/` folder — these are test files only
- Do not add new features until all 10 critical blockers are resolved
- Do not hardcode any secrets or credentials
- Do not change the clean architecture layer separation
- Do not commit `.env` to git under any circumstance
