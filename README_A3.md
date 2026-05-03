# Microservices Order & Payment System (Event-Driven)

## Overview
This project is a microservices-based system consisting of three services: **Order Service**, **Payment Service**, and **Notification Service**. It demonstrates synchronous communication via **gRPC** and asynchronous, event-driven communication using **RabbitMQ**.

## Architecture Components
*   **Order Service**: Handles order creation via REST API and communicates with the Payment Service using gRPC.
*   **Payment Service**: Processes payments, updates the `payments_db`, and acts as a **Producer**, publishing `payment.completed` events to RabbitMQ[cite: 2].
*   **Notification Service**: Acts as a **Consumer**, listening for payment events to send simulated email notifications[cite: 1].
*   **RabbitMQ**: The message broker managing event distribution and reliability[cite: 2].

---

## Technical Implementations

### 1. Reliability & ACK Logic
To ensure no messages are lost during transit or processing, we implemented **Manual Acknowledgments**:
*   **Disabled Auto-Ack**: `auto-ack` is set to `false` in the Notification Service to prevent messages from disappearing if the service crashes[cite: 1].
*   **Explicit Confirmation**: The service only sends a `d.Ack(false)` after the event is successfully parsed and logged[cite: 1].
*   **Durability**: Queues are declared as `durable: true`, and messages are published with `amqp.Persistent` to ensure they survive a RabbitMQ broker restart[cite: 2].

### 2. Idempotency Strategy
To prevent duplicate notifications (e.g., if a network error occurs during an ACK and the message is redelivered), the Notification Service uses an **In-Memory Tracking** strategy:
*   **Tracking Map**: A `map[string]bool` stores `order_id`s that have already been processed[cite: 1].
*   **Thread Safety**: A `sync.Mutex` is used to lock the map, preventing race conditions when processing multiple messages concurrently[cite: 1].
*   **Logic**: If a received `order_id` already exists in the map, the service logs a "Duplicate message ignored" warning and sends an `Ack` immediately without re-sending the notification[cite: 1].

### 3. Advanced Failure Handling (Dead Letter Queue)
We implemented a **Dead Letter Queue (DLQ)** to handle "poison messages" or permanent business logic failures:
*   **DLX Configuration**: The `payment.completed` queue is configured with `x-dead-letter-exchange` (pointing to `payment.dlx`)[cite: 2].
*   **Failure Trigger**: In `NotificationMain_2.go`, we simulate a permanent error for specific cases (e.g., `amount == 101` or `amount <= 0`)[cite: 1].
*   **Routing**: When a failure is detected, the service calls `d.Nack(false, false)`. RabbitMQ then automatically routes the message to the `payment.dlq` for manual developer inspection[cite: 1, 2].

---

## Infrastructure (Docker)
The environment is orchestrated using `docker-compose`, which manages:
*   **RabbitMQ**: Using the `rabbitmq:3-management` image to provide a web UI on `localhost:15672`.
*   **Databases**: Dedicated PostgreSQL instances for `orders_db` and `payments_db`.

## How to Run
1.  **Start the environment**:
    ```bash
    docker-compose up -d
    ```
2.  **Run the services**:
    *   Start the **Order Service**.
    *   Start the **Payment Service** (this creates the DLX and Queues).
    *   Start the **Notification Service**.
3.  **Demonstrate DLQ**:
    *   Send a POST request to the Order Service with `amount: 101`.
    *   Observe the Notification Service log: `"Moving message..."`.
    *   Check the RabbitMQ Management UI to see the message in the `payment.dlq` queue.

---