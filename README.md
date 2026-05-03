### Quick Start Guide
1. **Bring up Infrastructure**: `docker-compose up -d`
2. **Start Services**:
    * Order Service: `go run main.go`
    * Payment Service: `go run paymentMain.go`
    * Notification Service: `go run NotificationMain.go`
3. **Test DLQ Flow**:
    * Send a POST request to `:8081/orders` with `amount: 101`.
    * Check logs for `"Moving message..."`.
    * Verify message location in RabbitMQ UI (`localhost:15672`).

---

## System Architecture

The following diagram represents the final event-driven flow, showcasing both synchronous gRPC calls and asynchronous RabbitMQ event

---

# Microservices Platform: Order, Payment & Notification

This repository documents the evolution of a high-performance microservices-based platform, transitioning from simple REST communication to a robust, event-driven architecture using **gRPC** and **RabbitMQ**.

---

## Assignment 1: Two-Service Platform (REST)
*Focus: Clean Architecture & Service Decoupling*

The initial phase implements the foundation of the system using synchronous communication.

### Key Deliverables
* **Clean Architecture**: Strict separation of concerns across Domain, Usecase, Repository, and Transport layers.
* **Database Isolation**: Independent data persistence with `orders_db` and `payments_db`.
* **Microservice Autonomy**: Services are fully decoupled with no shared codebases.
* **Resilient REST**: Synchronous communication with a mandatory 2-second timeout for the Order-to-Payment flow.
* **Automatic Recovery**: Orders are automatically set to `Failed` if payment processing fails or times out.

---

## Assignment 2: gRPC Migration & Streaming
*Focus: Performance & Contract-First Development*

This stage optimizes inter-service communication by migrating to **gRPC**.

### Architecture Evolution
* **Contract-First Workflow**: Interfaces are strictly defined in Protocol Buffers before implementation.
* **Hybrid Service Roles**: The Order Service acts as both a gRPC client (Authorization) and a gRPC server (Status Updates).
* **Unary RPC**: Fast, synchronous payment validation with business rules (e.g., auto-decline for amounts > 100,000).
* **Server-Side Streaming**: Real-time status updates pushed via the `SubscribeToOrderUpdates` RPC stream.
* **Enhanced Resilience**: Implementation of a 5-second gRPC timeout and 503 error handling.

---

## Assignment 3: Event-Driven Architecture (RabbitMQ)
*Focus: Reliability & Failure Handling*

The final phase introduces asynchronous processing and advanced messaging reliability.

### Reliability & Messaging Logic
* **Manual Acknowledgments**: `auto-ack` is disabled; messages are only acknowledged (`d.Ack`) after confirmed processing.
* **Message Persistence**: Queues are marked as `durable` and messages use `Persistent` delivery mode.
* **Idempotency Strategy**: Thread-safe **In-Memory Tracking** using a `map` and `sync.Mutex` to prevent duplicate notifications.

### Advanced Failure Handling (DLQ)
* **Dead Letter Queue (DLQ)**: Configured with `x-dead-letter-exchange` to handle "poison messages".
* **Failure Trigger**: Permanent errors are simulated for specific cases (e.g., `amount: 101`).
* **DLX Routing**: Failed messages are routed via `d.Nack(false, false)` to the `payment.dlq` for manual inspection.

---

## Infrastructure & Setup

### Prerequisites
* **Go** 1.21+
* **PostgreSQL** (Databases: `orders_db`, `payments_db`)
* **Docker & Docker Compose**

### Environment Variables (.env)
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
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

```
### Architecture Diagram:
![Diagram_A3.png](Diagram_A3.png)

### Results:
* **Post_Orders:** ![POST_Orders.png](POST_Orders.png)

* **Post_Orders (Amount Logic):**![POST_Orders_Amount_Logic.png](POST_Orders_Amount_Logic.png)

* **Stream_Client:**![Stream_Client.png](Stream_Client.png)

* **WorkFlow:**![Workflow.png](Workflow.png)
