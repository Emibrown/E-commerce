# E-Commerce API

A sample **Golang** + **Gin** e-commerce RESTful API with **GORM** and **Swagger** documentation.  

## Features

- **User Management** (register, login, JWT auth)
- **Product Management** (CRUD, admin-only)
- **Order Management** (create, list, cancel, update status)
- **PostgreSQL**
- **Swagger**-based API documentation

---

## Prerequisites

1. **Go** (1.18+ recommended)  
2. **Git**  
3. **PostgreSQL** (or another supported DB, if you adapt the code)  
4. **[Swag CLI](https://github.com/swaggo/swag)** (optional, only if you’ll regenerate docs)

---

## Installation & Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/Emibrown/E-commerce
   cd E-commerce

2. Initialize/update Go modules:
    ```bash
    go mod tidy

3. Create your database (e.g., ecommerce_db) in PostgreSQL.

4. Set environment variables (either in a .env file or export them in your shell).

## Environment Variables

You can configure the application by setting the following variables:

    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=ecommerce_db
    DB_PORT=5432
    JWT_SECRET=mysecretkey
    HOST=localhost
    PORT=8080
    ADMIN_SECRET=supersecret


Place these in a .env file (recommended) or export them directly into your environment


## Running the App

    go run cmd/main.go

## Swagger Documentation

1. View the API docs in your browser at:

    ```bash
    http://<host>:<port>/swagger/index.html

2. If you edit the Swagger annotations or add new endpoints, regenerate the docs:

    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    swag init -g cmd/main.go

