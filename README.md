# Notify - Real-Time Notification System

`Notify` is a scalable, modular, and abstract real-time notification system built in Go. It leverages WebSocket for client communication and provides flexible interfaces for persistence and messaging, allowing you to use it with various databases (e.g., MongoDB, PostgreSQL, in-memory) and messaging brokers (e.g., Redis, RabbitMQ, in-memory). Whether you're building a simple chat app or a complex enterprise notification service, `Notify` adapts to your needs.

## Features

- **Real-Time Notifications**: Deliver messages instantly to connected WebSocket clients.
- **Targeted Delivery**: Send notifications to individual users or groups.
- **Scalability**: Abstracted messaging layer supports scalable brokers like Redis or RabbitMQ.
- **Persistence**: Abstracted storage layer allows durable notification history with databases like MongoDB or PostgreSQL.
- **Modularity**: Swap out persistence and messaging implementations without changing core code.
- **Easy Integration**: Works seamlessly with the Chi router for HTTP and WebSocket endpoints.

## Installation

To use `Notify` in your Go project, install it via `go get`:

```bash
go get github.com/kenztech/notify@v0.1.0
