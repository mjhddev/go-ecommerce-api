# Go E-Commerce API

A RESTful E-Commerce API built with **Go**, **Gin**, **GORM**, and **PostgreSQL**.

This project demonstrates modern backend development practices using a layered architecture, JWT authentication, role-based authorization, Docker, Swagger, and unit testing.

---

## Features

### Authentication

- User Registration
- User Login
- JWT Authentication
- Role-Based Authorization

### User

- Get User Profile

### Category

- Create Category
- Get Categories
- Get Category by ID
- Update Category
- Delete Category

### Product

- Create Product
- Get Products
- Get Product by ID
- Update Product
- Delete Product

### Shopping Cart

- Add Product to Cart
- View Cart
- Update Cart Item
- Remove Cart Item

### Order

- Checkout
- Order History
- Order Detail

### Documentation

- Swagger API Documentation

### Quality

- Docker Support
- Unit Testing

---

## Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go |
| Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL |
| Authentication | JWT |
| Documentation | Swagger |
| Containerization | Docker & Docker Compose |
| Testing | Testify |

---

## Project Structure

```text
.
├── .github/
│   └── workflows/
├── cmd/
│   └── api/
│       └── main.go
├── configs/
├── docs/
├── internal/
│   ├── dto/
│   ├── errs/
│   ├── handlers/
│   ├── middleware/
│   ├── migrations/
│   ├── models/
│   ├── repositories/
│   ├── response/
│   ├── routes/
│   ├── services/
│   ├── utils/
│   └── validators/
├── .dockerignore
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

---

## Architecture

This project follows a layered architecture.

```text
                 HTTP Request
                      │
                      ▼
                 Gin Router
                      │
                      ▼
                 Middleware
                      │
                      ▼
                   Handler
                      │
                      ▼
                   Service
                      │
                      ▼
                 Repository
                      │
                      ▼
                PostgreSQL
```

### Design Patterns

- Layered Architecture
- Repository Pattern
- Dependency Injection
- DTO Pattern
- Service Layer
- Transaction Management

---

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL
- Docker (Optional)

---

### Clone Repository

```bash
git clone https://github.com/mjhddev/go-ecommerce-api.git
```

```bash
cd go-ecommerce-api
```

---

### Install Dependencies

```bash
go mod tidy
```

---

### Environment Variables

Copy the example environment file.

```bash
cp .env.example .env
```

Update the values according to your local environment.

Example:

```env
APP_NAME=Go E-Commerce API
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ecommerce_db
DB_SSLMODE=disable

JWT_SECRET=your_jwt_secret
```

---

### Run Application

```bash
go run ./cmd/api
```

The server will start on:

```
http://localhost:8080
```

---

## Running with Docker

Build and start all services.

```bash
docker compose up --build
```

Or if the image has already been built.

```bash
docker compose up
```

Stop all containers.

```bash
docker compose down
```

---

## API Documentation

Swagger UI is available at:

```
http://localhost:8080/swagger/index.html
```

---

## API Endpoints

### Authentication

| Method | Endpoint |
|--------|----------|
| POST | /api/v1/auth/register |
| POST | /api/v1/auth/login |

### Profile

| Method | Endpoint |
|--------|----------|
| GET | /api/v1/profile |

### Categories

| Method | Endpoint |
|--------|----------|
| GET | /api/v1/categories |
| GET | /api/v1/categories/:id |
| POST | /api/v1/categories |
| PUT | /api/v1/categories/:id |
| DELETE | /api/v1/categories/:id |

### Products

| Method | Endpoint |
|--------|----------|
| GET | /api/v1/products |
| GET | /api/v1/products/:id |
| POST | /api/v1/products |
| PUT | /api/v1/products/:id |
| DELETE | /api/v1/products/:id |

### Cart

| Method | Endpoint |
|--------|----------|
| GET | /api/v1/cart |
| POST | /api/v1/cart |
| PUT | /api/v1/cart/:id |
| DELETE | /api/v1/cart/:id |

### Orders

| Method | Endpoint |
|--------|----------|
| POST | /api/v1/orders/checkout |
| GET | /api/v1/orders |
| GET | /api/v1/orders/:id |

---

## Sample Response

```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "name": "Mechanical Keyboard",
    "description": "RGB Mechanical Keyboard",
    "price": 750000,
    "stock": 10,
    "category": {
      "id": 2,
      "name": "Gaming"
    }
  }
}
```

---

## Testing

Run all unit tests.

```bash
go test ./...
```

Current test coverage includes:

- User Service
- Category Service
- Product Service
- Cart Service
- Order Service

---

## Future Improvements

- Integration Testing
- GitHub Actions CI
- Refresh Token
- Payment Gateway Integration
- Redis Caching
- Product Search
- Pagination
- Email Verification

---

## License

This project is created for learning purposes and portfolio development.