# Go E-Commerce API

A RESTful E-Commerce API built with Go, Gin, GORM, and PostgreSQL following Clean Architecture principles.

## Features

- JWT Authentication & Authorization
- Product Management
- Category Management
- Shopping Cart
- Checkout
- Order History
- Dashboard Statistics
- Product Image Upload
- Swagger Documentation
- Docker & Docker Compose
- Database Seeder

## Tech Stack

- Go
- Gin
- GORM
- PostgreSQL
- JWT
- Swagger
- Docker

## Project Structure

```text
.
├── cmd
├── configs
├── database
├── docs
├── internal
│   ├── dto
│   ├── handlers
│   ├── middleware
│   ├── models
│   ├── repositories
│   ├── routes
│   ├── services
│   ├── storage
│   └── utils
├── uploads
├── Dockerfile
└── docker-compose.yml
```

## Getting Started

Clone the repository.

```bash
git clone https://github.com/mjhddev/go-ecommerce-api.git
cd go-ecommerce-api
```

Start the application with Docker.

```bash
docker compose up --build
```

Run the database seeder.

```bash
docker compose run --rm seed
```

Open Swagger UI.

```
http://localhost:8080/swagger/index.html
```

## Future Improvements

- Redis Caching
- Payment Gateway Integration
- Unit Testing
- CI/CD Pipeline

```mermaid
erDiagram
    USERS ||--o{ CART_ITEMS : has
    USERS ||--o{ ORDERS : places

    CATEGORIES ||--o{ PRODUCTS : contains

    PRODUCTS ||--o{ CART_ITEMS : added_to
    PRODUCTS ||--o{ ORDER_ITEMS : purchased_as

    ORDERS ||--o{ ORDER_ITEMS : contains

    USERS {
        uint id PK
        string name
        string email
        string password
        string role
    }

    CATEGORIES {
        uint id PK
        string name
    }

    PRODUCTS {
        uint id PK
        uint category_id FK
        string name
        float price
        int stock
        string image_url
    }

    CART_ITEMS {
        uint id PK
        uint user_id FK
        uint product_id FK
        int quantity
    }

    ORDERS {
        uint id PK
        uint user_id FK
        float total_amount
        string status
    }

    ORDER_ITEMS {
        uint id PK
        uint order_id FK
        uint product_id FK
        int quantity
        float price
    }
```