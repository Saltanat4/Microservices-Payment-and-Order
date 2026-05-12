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

### Assignment 4: Performance, Security & Distributed Caching
*Focus: Infrastructure Resilience, Distributed Caching & Rate Limiting*

The final evolution introduces Redis as a high-performance middleware layer to protect the system from traffic spikes and ensure message delivery guarantees.

### Key Deliverables
* **Distributed Rate Limiting:** Integrated Redis-based middleware to prevent API abuse (Limit: 10 requests/min per IP).

* **Infrastructure-Level Idempotency:** Replaced in-memory tracking with a distributed Redis SetNX strategy for the Notification Service to handle duplicate events safely.

* **Resilient Retry Mechanism:** Implemented exponential backoff for external integrations (Email Mock) to handle 20% simulated network failures.

* **Enhanced Monitoring:** Structured logging for tracking message flow across RabbitMQ, Redis, and SMTP mocks.

### Invalidation & Reliability Logic
* **Rate Limiting Strategy:** Uses a Fixed Window algorithm in Redis.
     * Every request triggers an INCR on a key tied to the Client IP. 
     * If the counter exceeds the threshold, the system returns a 429 Too Many Requests.

* **Idempotency Strategy:** The Notification Worker uses notified:<order_id> as a unique key in Redis. 
  * Before processing, it executes SetNX (Set if Not Exists) with a 24-hour TTL. 
  * If the key already exists, the message is acknowledged as a duplicate and ignored.

* **Retry Logic and Trigger:** A simulated 20% chance of failure in the EmailSender mock.
  * **Algorithm:** Custom backoff starting at 5s, increasing by 5s for each step (5s -> 10s -> 15s). 
  * **Threshold:** After 3 failed attempts, the message is moved to the Dead Letter Queue (DLQ) via d.Nack(false, false).

### WorkFlow (The Final Chain)
1. **Request:** Client sends POST /orders. Redis Middleware checks IP rate limits. 
2. **Order:** Order Service saves to orders_db and calls Payment Service via gRPC. 
3. **Payment:** Payment Service validates and publishes a payment.completed event to RabbitMQ. 
4. **Notification:**
     * Worker consumes the message and checks Redis for duplicates. 
     * If Amount == 101, the message is immediately routed to DLQ. 
     * Otherwise, it attempts email delivery with a Retry loop for transient errors.
5. Completion: Upon success, a fmt.Printf log confirms the final delivery details including the total amount

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
---
# System Architecture
![Diagram_4.png](png/Diagram_4.png)

--- 
# Demo Screenshoots
* **Post_Orders:** ![POST_Orders.png](png/POST_Orders.png)

* **Post_Orders (Amount Logic):**![POST_Orders_Amount_Logic.png](png/POST_Orders_Amount_Logic.png)

* **Stream_Client:**![Stream_Client.png](png/Stream_Client.png)

* **WorkFlow:**![Wo![retries.png](png/retries.png)rkflow.png](png/Workflow.png)
* **Redis Rate Limit (429 Error):** ![429error.png](png/429error.png)
* **Worker Retries & Success:**![retries.png](png/retries.png)
* **Idempotency (Duplicate Ignored):**![duplicate.png](png/duplicate.png)

* **RabbitMQ DashBoard (DLQ):**
![RabbitMG DashBoard.png](png/RabbitMG%20DashBoard.png)
