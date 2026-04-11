# AP2_assignment1: Two-Service Platform (Order & Payment)

This project implements a simplified two-service platform using microservices architecture, Clean Architecture principles, and Go.

## Deliverables Checklist (Assignment 1 Requirements)

### 1. Core Architecture
* **Architecture:** Clean Architecture (Domain, Usecase, Repository, Transport) is implemented in both services.
* **Database Separation:** The services use decoupled databases (`orders_db` and `payments_db`).
* **Microservice Decomposition:** Services communicate via REST APIs and do not share code/packages.

### 2. Communication & Failure Handling
* **Communication:** Synchronous REST communication (Order Service calls Payment Service).
* **Resilience:** The HTTP client in Order Service includes a 2-second timeout.
* **Failure Handling:** If the Payment fails or is declined, the Order status is automatically updated to 'Failed'.

### 3. Repository
* Implements CRUD operations for Orders and Transactions.
* Correctly handles data types (e.g., storing `CreatedAt` as `time.Time` and mapping to SQL timestamp).

### 4. Code Quality & Formatting
* Code is logically organized and follows standard Go project structure.
* Includes `go.mod` and `go.sum` files.

## Local Setup & Run

1.  **Start Database:** Ensure PostgreSQL is running. Create `orders_db` and `payments_db`.
2.  **Environment:** Create a `.env` file in the root directory (refer to `.env.example`).
3.  **Run Payment Service:** `go run cmd/payment/main.go` (Runs on port 8080).
4.  **Run Order Service:** `go run cmd/order/main.go` (Runs on port 8081).