# Retail Management with Go | Semi-Microservices

A retail management backend implementing a Hybrid Architecture (Monolith + Microservice). The system separates core business logic from inventory management using gRPC for inter-service communication.

## Tech Stack

* **Language:** Go (Golang)
* **HTTP Framework:** Fiber
* **RPC Framework:** gRPC & Protobuf
* **Database:** MySQL
* **Architecture:** Monolith with Inventory Microservice
* **Key Libraries:** ULID, JWT, Logrus, Godotenv, Validator

## Project Structure

* **retail-monolith/**: Monolith service (HTTP API, Users, Transactions).
* **inventory-service/**: Microservice (gRPC Server, Stock Management).
* **database/**: SQL dumps for schema and seeding.
* **docs/**: OpenAPI specification and Postman collection.

## Setup & Installation

### 1. Database Setup
Create two MySQL databases and import the provided SQL files:
* **retail_management**: Import `database/retail_manager.sql`
* **retail_inventory**: Import `database/retail_inventory.sql`

### 2. Environment Configuration
**inventory-service/.env**
```ini
DB_USER=root
DB_PASS=your-password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=retail_inventory
DB_PARAMS="parseTime=true&loc=UTC"
GRPC_PORT=50051
````

**retail-monolith/.env**

```ini
SERVER_PORT=":3000"
ALLOWED_ORIGIN=your-allowed-origin or just leave this

DB_USER=root
DB_PASS=your-password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=retail_manager
DB_PARAMS="parseTime=true&loc=UTC"

JWT_SECRET_KEY=your-jwt-pw
```

### 3\. Running the Services

**Step A: Start Inventory Microservice**
This service must be running first.

```bash
cd inventory-service
go mod tidy
go run main.go
```

**Step B: Start Management API**
Open a new terminal.

```bash
cd retail-monolith
go mod tidy
go run main.go
```

## API Endpoints

Also can be accessed in `docs/openapi.json`.


| Module | Method | Endpoint | Description |
| :--- | :--- | :--- | :--- |
| **Auth** | POST | `/auth/login` | Login User & Get Token |
| | GET | `/auth/me` | Get Current Profile |
| **Users** | POST | `/users` | Create User (Admin only) |
| | GET | `/users` | Get All Users (Admin only) |
| | GET | `/users/:userId` | Get User by ID (Admin only) |
| | PATCH | `/users/:userId` | Update User (Admin only) |
| | DELETE | `/users/:userId` | Delete User (Admin only) |
| | GET | `/roles` | Get All Roles (Admin only) |
| **Categories** | POST | `/categories` | Create Category (Admin only) |
| | GET | `/categories` | Get All Categories |
| | PUT | `/categories/:categoryId` | Update Category (Admin only) |
| | DELETE | `/categories/:categoryId` | Delete Category (Admin only) |
| **Suppliers** | POST | `/suppliers` | Create Supplier (Admin only) |
| | GET | `/suppliers` | Get All Suppliers |
| | PATCH | `/suppliers/:supplierId` | Update Supplier (Admin only) |
| | DELETE | `/suppliers/:supplierId` | Delete Supplier (Admin only) |
| **Products** | POST | `/products` | Create Product + **Sync Stock (gRPC)** (Admin only) |
| | GET | `/products` | Get All Products + **Live Stock (gRPC)** |
| | GET | `/products/:productId` | Get Product by ID + **Live Stock (gRPC)** |
| | PATCH | `/products/:productId` | Update Product Details (Admin only) |
| | DELETE | `/products/:productId` | Delete Product (Admin only) |
| **Inventory** | POST | `/inventory/adjust` | Manual Stock Adjustment (**Proxy to gRPC**) (Admin only) |
| **Transactions**| POST | `/transactions` | Create Transaction + **Decrease Stock (gRPC)** (Cashier) |
| | GET | `/transactions` | Get Transaction History |
| | GET | `/transactions/:transactionId`| Get Transaction Detail by ID |

## Testing Flow

1.  Import `docs/postman_collection.json` into Postman.
2.  Run `POST /auth/login` to generate a token.
3.  Set the token in Authorization header (Bearer Token).
4.  Run `GET /products` to verify data aggregation from MySQL and gRPC.
5.  Run `POST /transactions` to verify distributed state changes (Stock decrement).