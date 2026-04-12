# Assignment 2: gRPC Migration & Contract-First Development

## 1. Overview
This project is an evolution of a microservice-based system, migrating inter-service communication from REST to **gRPC**. It follows a **Contract-First** approach, where service interfaces are defined in Protocol Buffers before implementation.

The system consists of:
* **Order Service**: A REST-facing gateway that acts as a gRPC client to authorize payments and a gRPC server for real-time status updates.
* **Payment Service**: A gRPC server that processes transactions based on business logic.

---

## 2. Architecture
The project adheres to **Clean Architecture** (Domain, Usecase, Repository, Transport), ensuring the core business logic remains independent of the communication protocol.

### System Design
* **Order Service**:
    * **REST API (Gin)**: Handles external client requests.
    * **gRPC Client**: Sends unary requests to the Payment Service.
    * **gRPC Server**: Provides a server-side stream (`SubscribeToOrderUpdates`) for real-time notifications.
* **Payment Service**:
    * **gRPC Server**: Implements `ProcessPayment` to validate and store transactions.

---

## 3. Contract-First Workflow
We utilize a dedicated workflow to ensure service contracts are the single source of truth.

* **Proto Repository**: [https://github.com/Saltanat4/protos-repo](https://github.com/Saltanat4/protos-repo)
    * Contains the `.proto` files defining the data structures and service methods.
* **Generated Code Repository**: [https://github.com/Saltanat4/gen-repo](https://github.com/Saltanat4/gen-repo)
    * Contains the compiled Go code (`.pb.go`), automated via GitHub Actions.

---

## 4. Setup & Installation

### Prerequisites
* **Go 1.21+**
* **PostgreSQL** (Two databases: `orders_db` and `payments_db`)

### Step-by-Step Run Guide
1.  **Clone the Repository**:
    ```bash
    git clone <your-repo-link>
    cd AP2_assignment2
    ```
2.  **Configure Environment Variables**:
    Create a `.env` file in the root directory:
    ```env
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_HOST=localhost
    DB_PORT=5432
    DB_NAME_ORDER=orders_db
    DB_NAME_PAYMENT=payments_db

    PAYMENT_GRPC_PORT=50051
    ORDER_HTTP_PORT=8081
    ORDER_GRPC_PORT=50052
    PAYMENT_GRPC_ADDRESS=localhost:50051
    ```
3.  **Install Dependencies**:
    ```bash
    go mod tidy
    ```
4.  **Run Payment Service**:
    ```bash
    go run paymentMain.go
    ```
5.  **Run Order Service**:
    ```bash
    go run main.go
    ```

---

## 5. API Reference & Business Rules
* **Unary RPC (`ProcessPayment`)**: Used for synchronous authorization. Payments exceeding **100,000** are automatically **Declined**.
* **Streaming RPC (`SubscribeToOrderUpdates`)**: A server-side stream that pushes database status changes to the client in real-time.
* **Resilience**: The Order Service uses a **5-second gRPC timeout**. If the Payment Service is unavailable, the REST API returns **503 Service Unavailable**.

---

## 6. Evidences

### Evidence 1: Successful Unary gRPC Call
* **What it proves**: Successful communication between Order Client and Payment Server.
* **Source**: Postman REST output and Payment terminal logs.

### Evidence 2: Real-Time Status Streaming
* **What it proves**: Functionality of `SubscribeToOrderUpdates` pushing updates (Pending -> Paid).
* **Source**: Postman gRPC Client interface.

### Evidence 3: Resilience & Error Handling
* **What it proves**: System returns **503** when the Payment Service is down.
* **Source**: Postman REST response.