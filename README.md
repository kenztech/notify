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
```

## How To Use

This section walks you through creating a simple real-time notification system using WebSockets and HTTP requests. The system consists of a frontend and a backend.

## Frontend (static/index.html)

### Description

The frontend is a simple web page that allows users to connect to the WebSocket server and send/receive notifications.

### Code

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Notify</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: sans-serif; background: white; color: black; }
        .container { display: grid; grid-template-columns: 320px 1fr; height: 100vh; }
        .sidebar { border-right: 1px solid #e0e0e0; padding: 1.5rem; background: white; }
        .main-content { background: #f5f5f5; display: flex; flex-direction: column; }
        header { padding: 1rem; background: white; display: flex; justify-content: space-between; }
        #status { font-size: 0.75rem; color: black; }
        .message-form { padding: 1rem; display: flex; gap: 0.75rem; }
        input, select, button { padding: 0.75rem; border-radius: 0.5rem; }
        button { background: black; color: white; cursor: pointer; }
    </style>
</head>
<body>
    <div class="container">
        <aside class="sidebar">
            <header><h1>Settings</h1></header>
            <form class="connection-form">
                <label for="userId">User ID</label>
                <input type="text" id="userId" value="user123" />
                <label for="groupId">Group ID</label>
                <input type="text" id="groupId" value="demo-users" />
                <button type="button" onclick="connect()">Connect</button>
            </form>
        </aside>
        <main class="main-content">
            <header><h1>Notifications</h1><div id="status">Disconnected</div></header>
            <div id="messages"></div>
            <footer class="message-form">
                <input type="text" id="message" placeholder="Type a message..." />
                <select id="type"><option value="info">Info</option><option value="alert">Alert</option></select>
                <button onclick="sendNotification()">Send</button>
            </footer>
        </main>
    </div>
    <script>
        let ws;
        function connect() {
            const userId = document.getElementById('userId').value;
            const groupId = document.getElementById('groupId').value;
            ws = new WebSocket(`ws://localhost:8080/ws?userId=${userId}&groupId=${groupId}`);
            ws.onopen = () => document.getElementById('status').textContent = `Connected as ${userId}`;
            ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                const messages = document.getElementById('messages');
                const msgEl = document.createElement('p');
                msgEl.textContent = `${data.type.toUpperCase()}: ${data.message}`;
                messages.appendChild(msgEl);
            };
        }
        function sendNotification() {
            const type = document.getElementById('type').value;
            const message = document.getElementById('message').value;
            fetch(`http://localhost:8080/notify/${type}/${encodeURIComponent(message)}`, { method: 'POST' });
        }
    </script>
</body>
</html>
```

## Backend (main.go)

### Description - Backend

The backend is a Go server using `go-chi` for routing and `notify` for handling real-time notifications.

### Code - Backend

```go
package main

import (
 "log"
 "net/http"
 "sync"
 "github.com/go-chi/chi/v5"
 "github.com/go-chi/chi/v5/middleware"
 "github.com/kenztech/notify"
 "github.com/kenztech/notify/models"
)

type InMemoryStore struct {
 data  map[string]models.Notification
 mutex sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
 return &InMemoryStore{data: make(map[string]models.Notification)}
}

func (m *InMemoryStore) SaveNotification(n models.Notification) error {
 m.mutex.Lock()
 defer m.mutex.Unlock()
 m.data[n.ID] = n
 return nil
}

func (m *InMemoryStore) GetNotifications(targetID string, groupIDs []string) ([]models.Notification, error) {
 m.mutex.RLock()
 defer m.mutex.RUnlock()
 var result []models.Notification
 for _, n := range m.data {
  if n.TargetID == targetID || containsAny(n.GroupIDs, groupIDs) {
   result = append(result, n)
  }
 }
 return result, nil
}

func containsAny(slice, items []string) bool {
 for _, s := range slice {
  for _, i := range items {
   if s == i {
    return true
   }
  }
 }
 return false
}

func main() {
 r := chi.NewRouter()
 r.Use(middleware.Logger)
 store := NewInMemoryStore()
 notifier := notify.NewNotifier(store)
 r.Route("/notify", func(r chi.Router) {
  r.Post("/{type}/{message}", func(w http.ResponseWriter, r *http.Request) {
   typeMsg := chi.URLParam(r, "type")
   message := chi.URLParam(r, "message")
   notifier.SendNotification(models.Notification{
    Type:    typeMsg,
    Message: message,
   })
   w.WriteHeader(http.StatusOK)
  })
 })

 r.HandleFunc("/ws", notifier.HandleWebSocket)

 log.Println("Server started on :8080")
 http.ListenAndServe(":8080", r)
}
```

## Running the System

### 1. Start the Backend

```sh
go run main.go
```

### 2. Start a Simple HTTP Server for Frontend

```sh
cd static
python3 -m http.server 8081
```

### 3. Open the Frontend

Visit `http://localhost:8081` and connect with a user ID.

### 4. Send Notifications

Use the input fields and buttons to send notifications.

## Conclusion

This notification system allows users to receive real-time messages using WebSockets. The backend efficiently handles message storage and delivery using an in-memory implementation.
