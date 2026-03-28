# GO-Eats Backend - Technical Architecture Documentation

## Executive Summary

**GO-Eats** is a production-ready Go-based food delivery backend system implementing clean architecture principles. The system serves RESTful APIs for user management, restaurant operations, cart/order processing, delivery personnel management with 2FA, reviews, and real-time notifications via WebSocket and NATS messaging.

---

## Technology Stack

| Component          | Technology                          | Version    |
|--------------------|-------------------------------------|------------|
| Language           | Go                                  | 1.24.4     |
| Web Framework      | Gin                                 | v1.10.1    |
| ORM                | Bun                                 | v1.2.15    |
| Database           | PostgreSQL                          | 16         |
| Message Broker     | NATS                                | v1.43.0    |
| Real-time Comm     | gorilla/websocket                   | v1.5.3     |
| Authentication     | golang-jwt/jwt                      | v5.2.3     |
| 2FA                | pquerna/otp (TOTP)                  | v1.5.0     |
| Validation         | go-playground/validator             | v10.27.0   |
| Password Hashing   | bcrypt                              | v0.40.0    |
| Logging            | slog-gin                            | v1.15.1    |
| Testing            | testify, testcontainers             | v0.38.0    |
| Containerization   | Docker Compose                      | -          |

---

## System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                            │
│              (Mobile App / Web / Postman / curl)                │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      GIN HTTP SERVER                            │
│  • Port: 8080                                                   │
│  • CORS Middleware                                              │
│  • Structured Logging (slog)                                    │
│  • JWT Authentication Middleware                                │
│  • Request Recovery                                             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      HANDLER LAYER                              │
│  ┌─────────┬──────────┬────────┬──────────┬──────────┐         │
│  │  User   │Restaurant│  Cart  │ Delivery │  Review  │         │
│  │ Handler │ Handler  │ Handler │ Handler  │ Handler  │         │
│  └─────────┴──────────┴────────┴──────────┴──────────┘         │
│  ┌──────────────────┬──────────────────┐                        │
│  │ Announcement     │ Notification     │                        │
│  │ Handler (SSE)    │ Handler (WS)     │                        │
│  └──────────────────┴──────────────────┘                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      SERVICE LAYER                              │
│  • Business Logic Implementation                                │
│  • Transaction Management                                       │
│  • NATS Publishing                                              │
│  • Validation                                                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                  REPOSITORY LAYER (Generic)                     │
│  • Insert, Delete, Select, Update                               │
│  • Relation Loading                                             │
│  • Raw SQL Support                                              │
│  • Health Check                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      DATA LAYER                                 │
│  ┌─────────────────┐  ┌─────────────────┐                      │
│  │   PostgreSQL    │  │   NATS Broker   │                      │
│  │   (Port 5432)   │  │   (Port 4222)   │                      │
│  └─────────────────┘  └─────────────────┘                      │
└─────────────────────────────────────────────────────────────────┘
```

---

## Project Structure

```
GO-Eats/
├── cmd/
│   └── api/
│       ├── main.go              # Application entry point, wire all components
│       └── middleware/
│           └── middleware.go    # JWT authentication middleware
│
├── pkg/
│   ├── abstract/                # Interface contracts (Domain-driven design)
│   │   ├── user/user.go         # User interface definitions
│   │   ├── restaurant/          # Restaurant interfaces
│   │   ├── cart/cart.go         # Cart interfaces
│   │   ├── delivery/delivery.go # Delivery interfaces
│   │   ├── review/review.go     # Review interfaces
│   │   ├── announcements/       # Announcement interfaces
│   │   └── config/env.go        # Environment loader
│   │
│   ├── database/
│   │   ├── database.go          # Generic repository implementation
│   │   └── models/              # Domain models (Bun ORM structs)
│   │       ├── user/user.go
│   │       ├── restaurant/      # Restaurant, MenuItem
│   │       ├── cart/            # Cart, CartItems
│   │       ├── order/           # Order, OrderItems
│   │       ├── delivery/        # DeliveryPerson, Deliveries
│   │       ├── review/          # Review
│   │       └── utils/           # Timestamp, common types
│   │
│   ├── handler/                 # HTTP layer (controllers)
│   │   ├── server.go            # Gin server setup, CORS, logging
│   │   ├── user/                # User HTTP handlers
│   │   ├── restaurant/          # Restaurant HTTP handlers
│   │   ├── cart/                # Cart & Order HTTP handlers
│   │   ├── delivery/            # Delivery HTTP handlers
│   │   ├── review/              # Review HTTP handlers
│   │   ├── announcements/       # SSE announcement handler
│   │   └── notification/        # WebSocket notification handler
│   │
│   ├── service/                 # Business logic layer
│   │   ├── user/                # User business logic
│   │   ├── restaurant/          # Restaurant business logic
│   │   ├── cart_order/          # Cart & Order business logic
│   │   ├── delivery/            # Delivery business logic (2FA, TOTP)
│   │   ├── review/              # Review business logic
│   │   ├── announcements/       # Announcement business logic
│   │   └── notification/        # NATS subscription logic
│   │
│   ├── nats/
│   │   └── nats_server.go       # NATS connection, pub/sub
│   │
│   ├── storage/                 # Image storage abstraction
│   │
│   └── tests/                   # Integration & unit tests
│       ├── user/
│       ├── restaurant/
│       ├── cart/
│       ├── delivery/
│       └── database/
│
├── httpclient/
│   ├── FoodDelivery.http        # HTTP request collection
│   └── FoodDelivery.postman_collection.json
│
├── docker-compose.yml           # PostgreSQL container
├── go.mod                       # Go module definition
├── .env                         # Environment configuration
└── uploads/                     # Local file storage
```

---

## Database Schema

### Entity Relationship Overview

```
┌─────────────┐       ┌───────────────┐       ┌─────────────┐
│    users    │       │  restaurant   │       │  menu_item  │
├─────────────┤       ├───────────────┤       ├─────────────┤
│ id (PK)     │       │ restaurant_id │       │ menu_id     │
│ name        │       │ name          │       │ restaurant_id (FK)
│ email       │       │ photo         │       │ name        │
│ password    │       │ description   │       │ description │
│ timestamps  │       │ address       │       │ price       │
│             │       │ city, state   │       │ category    │
│             │       │ timestamps    │       │ available   │
│             │       │               │       │ timestamps  │
└─────────────┘       └───────────────┘       └─────────────┘
       │                      │                      │
       │                      └──────────────────────┘
       │                              │
       ▼                              ▼
┌─────────────┐       ┌───────────────┐       ┌─────────────┐
│   reviews   │       │     cart      │       │  cart_items │
├─────────────┤       ├───────────────┤       ├─────────────┤
│ review_id   │       │ cart_id (PK)  │       │ cart_item_id│
│ user_id (FK)│       │ user_id (FK)  │       │ cart_id (FK)│
│ restaurant_id│      │ timestamps    │       │ item_id (FK)│
│ rating      │       │               │       │ restaurant_id│
│ comment     │       │               │       │ quantity    │
│ timestamps  │       │               │       │ timestamps  │
└─────────────┘       └───────────────┘       └─────────────┘
                              │
                              ▼
┌─────────────┐       ┌───────────────┐       ┌─────────────┐
│   orders    │◄──────│  order_items  │       │ deliveries  │
├─────────────┤       ├───────────────┤       ├─────────────┤
│ order_id    │       │ order_item_id │       │ delivery_id │
│ user_id (FK)│       │ order_id (FK) │       │ delivery_person_id│
│ order_status│       │ item_id (FK)  │       │ order_id (FK)│
│ total_amount│       │ restaurant_id │       │ delivery_status│
│ address     │       │ quantity      │       │ delivery_time│
│ timestamps  │       │ price         │       │ timestamps  │
│             │       │ timestamps    │       │             │
└─────────────┘       └───────────────┘       └─────────────┘
                              ▲
                              │
                    ┌─────────────────┐
                    │ delivery_person │
                    ├─────────────────┤
                    │ delivery_person_id│
                    │ name            │
                    │ phone (unique)  │
                    │ vehicle_details │
                    │ status          │
                    │ auth_key (TOTP) │
                    │ is_auth_set     │
                    │ timestamps      │
                    └─────────────────┘
```

---

## API Reference

### User Module

| Method | Endpoint       | Request Body                    | Response              | Auth |
|--------|----------------|---------------------------------|-----------------------|------|
| POST   | `/user/`       | `{name, email, password}`       | `201: {message}`      | No   |
| POST   | `/user/login`  | `{email, password}`             | `200: {token}`        | No   |
| DELETE | `/user/:id`    | -                               | `204: No Content`     | Yes  |

### Restaurant Module

| Method | Endpoint                    | Request Body                    | Response              | Auth |
|--------|-----------------------------|---------------------------------|-----------------------|------|
| POST   | `/restaurant`               | `multipart/form-data`           | `201: {message}`      | No   |
| GET    | `/restaurant`               | -                               | `200: []Restaurant`   | No   |
| GET    | `/restaurant/:id`           | -                               | `200: Restaurant`     | No   |
| DELETE | `/restaurant/:id`           | -                               | `204: No Content`     | No   |
| POST   | `/restaurant/menu`          | `{restaurant_id, name, desc, price, category, available}` | `201` | No   |
| GET    | `/restaurant/menu`          | -                               | `200: []MenuItem`     | No   |
| DELETE | `/restaurant/menu/:rid/:mid`| -                               | `204: No Content`     | No   |

### Cart & Order Module

| Method | Endpoint                        | Request Body              | Response              | Auth |
|--------|---------------------------------|---------------------------|-----------------------|------|
| POST   | `/cart/add`                     | `{item_id, restaurant_id, quantity}` | `201`         | Yes  |
| GET    | `/cart/list`                    | -                         | `200: {items}`        | Yes  |
| DELETE | `/cart/remove/:id`              | -                         | `204`                 | Yes  |
| POST   | `/cart/order/new`               | -                         | `201: {order}`        | Yes  |
| GET    | `/cart/orders`                  | -                         | `200: []Order`        | Yes  |
| GET    | `/cart/orders/:id`              | -                         | `200: OrderItems`     | Yes  |
| GET    | `/cart/orders/deliveries/:id`   | -                         | `200: []Delivery`     | Yes  |

### Review Module

| Method | Endpoint                    | Request Body              | Response              | Auth |
|--------|-----------------------------|---------------------------|-----------------------|------|
| POST   | `/review/:restaurant_id`    | `{rating (1-5), comment}` | `201: {message}`      | Yes  |
| GET    | `/review/:restaurant_id`    | -                         | `200: []Review`       | Yes  |
| DELETE | `/review/:review_id`        | -                         | `204`                 | Yes  |

### Delivery Module

| Method | Endpoint                    | Request Body              | Response              | Auth |
|--------|-----------------------------|---------------------------|-----------------------|------|
| POST   | `/delivery/add`             | `{name, phone, vehicle_details}` | `201`              | No   |
| POST   | `/delivery/login`           | `{phone, otp}`            | `200: {token}`        | No   |
| POST   | `/delivery/update-order`    | `{order_id, status}`      | `200`                 | Yes  |
| GET    | `/delivery/deliveries/:order_id` | -                      | `200: []Delivery`     | Yes  |

### Announcement Module

| Method | Endpoint                    | Request Body              | Response              | Auth |
|--------|-----------------------------|---------------------------|-----------------------|------|
| GET    | `/announcements/events`     | -                         | `SSE stream`          | No   |

### Notification Module

| Method | Endpoint                    | Request Body              | Response              | Auth |
|--------|-----------------------------|---------------------------|-----------------------|------|
| GET    | `/notify/ws`                | Query: `?token=<jwt>`     | `WebSocket`           | Yes  |

---

## Authentication Implementation

### User JWT Flow
```go
// middleware/middleware.go
type UserClaims struct {
    UserID int64  `json:"user_id"`
    Name   string `json:"name"`
    jwt.RegisteredClaims
}

// Token expires after 2 hours
ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour))
```

### Delivery Person 2FA Flow
1. **Registration:** TOTP secret generated via `github.com/pquerna/otp/totp`
2. **QR Code:** Google Authenticator URL stored in `auth_key_url`
3. **Login:** OTP validated against stored secret
4. **Token:** JWT issued upon successful validation

---

## Real-time Architecture

### NATS Topics
```
orders.new.<user_id>     → New order placement notifications
orders.status.<user_id>  → Order status change notifications
```

### WebSocket Message Format
```
USER_ID:<id>|MESSAGE:<content>
```

### SSE Events (Announcements)
```
event: message
data: <event_message>
```

---

## Design Patterns Applied

### 1. Repository Pattern
Generic database interface with type-safe operations:
```go
type Database interface {
    Insert(ctx context.Context, model any) (sql.Result, error)
    Delete(ctx context.Context, tableName string, filter Filter) (sql.Result, error)
    Select(ctx context.Context, model any, columnName string, parameter any) error
    Update(ctx context.Context, tableName string, Set Filter, Condition Filter) (sql.Result, error)
    // ...
}
```

### 2. Service Layer Pattern
Separation of concerns:
- **Handler:** HTTP handling, validation, response formatting
- **Service:** Business logic, transactions, external calls
- **Repository:** Data persistence

### 3. Dependency Injection
Constructor-based DI:
```go
func NewUserService(db database.Database, env string) *UsrService
func NewCartService(db database.Database, env string, nats *nats.NATS) *CartService
```

### 4. Interface Segregation
Domain-specific interfaces in `pkg/abstract/`:
```go
type Cart interface {
    Create(ctx context.Context, cart *cart.Cart) (*cart.Cart, error)
    GetCartId(ctx context.Context, UserId int64) (*cart.Cart, error)
    AddItem(ctx context.Context, Item *cart.CartItems) (*cart.CartItems, error)
}
```

---

## Testing Strategy

### Test Pyramid
```
         ┌─────────────┐
         │   E2E/      │  ← Minimal (manual via Postman)
         │   Manual    │
         └─────────────┘
         ┌─────────────┐
         │Integration  │  ← HTTP handlers + testcontainers
         │   Tests     │
         └─────────────┘
         ┌─────────────┐
         │   Unit      │  ← Service layer logic
         │   Tests     │
         └─────────────┘
```

### Test Infrastructure
- **testcontainers-go:** Spin up PostgreSQL in Docker for integration tests
- **testify/assert:** BDD-style assertions
- **go-faker/faker:** Generate realistic test data

### Example Test Pattern
```go
func TestMain(m *testing.M) {
    setup()     // Start PostgreSQL container
    result := m.Run()
    teardown()  // Cleanup containers
}

func TestAddUser(t *testing.T) {
    testDB := tests.Setup()
    testServer := handler.NewServer(testDB, false)
    // ... arrange, act, assert
}
```

---

## Configuration Management

### Environment Variables
| Variable           | Description                    | Example            |
|--------------------|--------------------------------|--------------------|
| `APP_ENV`          | Environment mode               | `dev`, `prod`      |
| `DB_HOST`          | PostgreSQL host                | `localhost`        |
| `DB_USERNAME`      | Database username              | `postgres`         |
| `DB_PASSWORD`      | Database password              | `***`              |
| `DB_NAME`          | Database name                  | `food-delivery`    |
| `DB_PORT`          | Database port                  | `5432`             |
| `STORAGE_TYPE`     | Image storage backend          | `local`            |
| `STORAGE_DIRECTORY`| Upload endpoint path           | `uploads`          |
| `LOCAL_STORAGE_PATH`| Filesystem path for uploads   | `/path/to/uploads` |
| `JWT_SECRET_KEY`   | JWT signing secret             | `***`              |

---

## Running the Application

### Start Database
```bash
docker-compose up -d
# Creates: go-eats-db (PostgreSQL 16 on port 5432)
```

### Run Server
```bash
go run cmd/api/main.go
# Server starts on :8080
```

### Run Tests
```bash
# Run all tests
go test ./... -v

# Run specific package tests
go test ./pkg/tests/user -v

# Run with race detector
go test ./... -race
```

### API Testing
```bash
# Health check
curl http://localhost:8080/healthz

# Create user
curl -X POST http://localhost:8080/user/ \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"secret123"}'
```

---

## Security Posture

### Current Implementation

| Aspect               | Status          | Details                               |
|----------------------|-----------------|---------------------------------------|
| Password Storage     | ✅ Implemented  | bcrypt with default cost              |
| JWT Signing          | ✅ Implemented  | HS256 symmetric signing               |
| Token Expiry         | ✅ Implemented  | 2 hours                               |
| Input Validation     | ⚠️ Partial      | go-playground/validator v10 (incomplete coverage) |
| SQL Injection        | ✅ Protected    | Bun ORM (parameterized queries)       |
| CORS                 | ❌ Insecure     | Allows all origins (`*`)              |
| 2FA                  | ✅ Implemented  | TOTP for delivery personnel           |
| Rate Limiting        | ❌ Missing      | No protection against brute force     |
| Password Policy      | ❌ Missing      | No strength requirements              |
| Audit Logging        | ⚠️ Basic        | slog structured logging               |

### Required Security Fixes (Pre-Production)

1. **CORS Restriction** - Replace `[]string{"*"}` with allowed frontend domains
2. **Rate Limiting Middleware** - Add `gin-rate-limit` or custom middleware
3. **Password Validation** - Add custom validator for min length (8) + complexity
4. **JWT Secret Management** - Use environment injection, not `.env` file

---

## Performance Characteristics

### Current Implementation

| Component            | Status          | Consideration                           |
|----------------------|-----------------|-----------------------------------------|
| Database             | ⚠️ Default       | Bun ORM with default connection pool    |
| Messaging            | ✅ Implemented  | NATS (low-latency pub/sub)              |
| Static Files         | ✅ Implemented  | Gin static middleware                   |
| Real-time            | ✅ Implemented  | WebSocket + NATS subscription           |
| Indexes              | ❌ Incomplete   | PKs auto-indexed; FKs NOT indexed       |
| Caching              | ❌ Missing      | No caching layer                        |

### Required Performance Fixes (Pre-Production)

1. **Connection Pool Configuration** - Set `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`
2. **Foreign Key Indexes** - Add indexes on all FK columns (`user_id`, `restaurant_id`, `order_id`, etc.)
3. **N+1 Query Prevention** - Use `SelectWithRelation` appropriately, add `JOIN` indexes

---

## Known Limitations & Tech Debt

### 🔴 Critical Blockers (Must Fix Before Production)

These issues **block production deployment** and must be resolved:

| # | Issue | File/Location | Risk |
|---|-------|---------------|------|
| 1 | **CORS allows all origins** | `pkg/handler/server.go:25` | Any website can make authenticated requests on behalf of users |
| 2 | **Hardcoded delivery address** | `pkg/service/cart_order/place_order.go:33` | All orders use "New Delhi" regardless of user input |
| 3 | **No rate limiting** | N/A | API vulnerable to brute force and DoS attacks |
| 4 | **No password validation** | `pkg/database/models/user/user.go` | Weak passwords allowed |
| 5 | **No migration system** | `pkg/database/database.go:201` | Schema changes unversioned, rollback impossible |
| 6 | **Missing FK indexes** | All model files | N+1 queries, severe performance degradation at scale |
| 7 | **No connection pool config** | `pkg/database/database.go:185` | Connection exhaustion under load |
| 8 | **Incomplete input validation** | Multiple handlers | Malformed data can reach database |
| 9 | **No order status state machine** | `pkg/service/cart_order/` | Invalid status transitions possible |
| 10 | **JWT secret in .env file** | `.env:11` | Secrets should not be in repository |

---

### 🟡 Future Scope (Can Deploy Without)

These improvements can be added post-launch:

| # | Issue | When to Address |
|---|-------|-----------------|
| 1 | File uploads: Local storage only | When scale requires S3/CDN |
| 2 | No pagination on list endpoints | When tables grow >1000 rows |
| 3 | No monitoring/observability | After first 100 users |
| 4 | No OpenAPI/Swagger spec | When onboarding external devs |
| 5 | No API versioning (`/api/v1/`) | When multiple API consumers exist |
| 6 | No payment integration | When online payments required |
| 7 | No email notifications | When user retention needs it |
| 8 | No idempotency keys | When payment integration added |
| 9 | No unit test coverage | When refactoring core logic |
| 10 | No distributed tracing | When microservices added |

---

## Deployment Considerations

### Production Readiness Checklist

#### 🔴 Critical (Blockers)
- [ ] **CORS origins restricted** to specific frontend domains
- [ ] **Rate limiting middleware** implemented
- [ ] **Password validation** with strength requirements
- [ ] **Migration system** (golang-migrate) integrated
- [ ] **FK indexes** added to all foreign key columns
- [ ] **Connection pool** configured for production
- [ ] **Hardcoded delivery address** removed from order placement
- [ ] **Order status state machine** implemented
- [ ] **Input validation** coverage on all handlers
- [ ] **JWT secret** via secure env injection (not in repo)

#### 🟡 Recommended (Post-Launch)
- [ ] TLS/HTTPS termination
- [ ] Log aggregation (ELK/Loki)
- [ ] Health check endpoints for load balancer
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Prometheus metrics endpoint
- [ ] API documentation (OpenAPI/Swagger)
- [ ] Pagination on list endpoints
- [ ] Docker multi-stage build

---

## Future Enhancements (Post-Production)

These features are **not required for initial launch** and can be added based on user feedback and scale:

| Priority | Feature                    | Description                          | Trigger to Build              |
|----------|----------------------------|--------------------------------------|-------------------------------|
| High     | Payment Integration        | Stripe/Razorpay integration          | When online payments required |
| High     | Email Notifications        | Order confirmations via SendGrid     | When user retention needs it  |
| Medium   | Admin Dashboard APIs       | Restaurant/order management          | When ops team hired           |
| Medium   | Redis Caching              | Cache frequently accessed data       | When query load increases     |
| Medium   | API Versioning             | `/api/v1/` prefix                    | When multiple API consumers   |
| Medium   | S3/Cloud Storage           | Migrate from local file storage      | When CDN needed               |
| Low      | GraphQL                    | Alternative to REST                  | When complex queries needed   |
| Low      | gRPC                       | Internal service communication       | When microservices added      |
| Low      | SMS Notifications          | Twilio integration                   | When SMS marketing needed     |
| Low      | Analytics Pipeline         | User behavior tracking               | When data-driven decisions    |

---

## Project Stage Assessment

### Current Development Stage: **MVP / Pre-Production**

The backend is functional for development and testing but **NOT production-ready**. Here's the assessment:

| Area              | Status     | Summary                                          |
|-------------------|------------|--------------------------------------------------|
| Core Features     | ✅ Complete | User, restaurant, cart, order, review, delivery  |
| Authentication    | ✅ Complete | JWT for users, TOTP+JWT for delivery personnel   |
| Real-time         | ✅ Complete | WebSocket + NATS pub/sub working                 |
| Testing           | ⚠️ Basic    | Integration tests exist, unit coverage limited   |
| Security          | ❌ Incomplete | Missing rate limiting, CORS, password policy   |
| Database          | ❌ Incomplete | No migrations, missing indexes                 |
| Business Logic    | ❌ Incomplete | Hardcoded values, no state machines            |
| DevOps            | ❌ Missing  | No CI/CD, no containerized production image      |
| Observability     | ❌ Missing  | Basic logging only, no metrics or tracing        |

### Production Readiness Score: **~40%**

**10 critical blockers** must be resolved before deployment. Estimated effort: **10-12 hours** of focused development.

---

## Version History

| Version | Date       | Changes                                    |
|---------|------------|-------------------------------------------|
| 1.1     | 2026-03-24 | Production readiness audit, critical vs future scope separation |
| 1.0     | 2026-03-24 | Initial architecture documentation        |

---

**Document Version:** 1.1
**Generated:** 2026-03-24
**Project:** GO-Eats Food Delivery Backend
**Repository:** github.com/Ayocodes24/GO-Eats
