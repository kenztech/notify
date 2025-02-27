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
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        :root {
            --primary: #000000;
            --bg-color: #ffffff;
            --text-color: #000000;
            --border-color: #e0e0e0;
            --hover-color: #333333;
            --success-color: green;
            --error-color: #666666;
            --input-bg: #ffffff;
            --message-bg: #f5f5f5;
            --message-sent: #000000;
            --message-received: #ffffff;
            --sidebar-width: 320px;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: var(--bg-color);
            color: var(--text-color);
            line-height: 1.4;
            height: 100vh;
            overflow: hidden;
        }

        .container {
            display: grid;
            grid-template-columns: var(--sidebar-width) 1fr;
            height: 100vh;
            background: var(--bg-color);
        }

        .sidebar {
            border-right: 1px solid var(--border-color);
            display: flex;
            flex-direction: column;
            background: var(--bg-color);
        }

        .main-content {
            display: flex;
            flex-direction: column;
            background: var(--message-bg);
        }

        header {
            padding: 1rem;
            background: var(--bg-color);
            color: black;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        header h1 {
            font-size: 1.25rem;
            font-weight: 500;
            letter-spacing: 0.5px;
        }

        #status {
            font-size: 0.75rem;
            opacity: 0.9;
            display: flex;
            align-items: center;
            gap: 0.375rem;
            color: black;
        }

        #status::before {
            content: '';
            display: inline-block;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background: var(--error-color);
        }

        #status.connected::before {
            background: var(--success-color);
        }

        .connection-panel {
            padding: 1.5rem;
            border-bottom: 1px solid var(--border-color);
            background: white;
        }

        #messages {
            flex: 1;
            padding: 1rem;
            overflow-y: auto;
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
        }

        #messages p {
            padding: 0.75rem 1rem;
            border-radius: 0.75rem;
            max-width: 80%;
            position: relative;
            font-size: 0.9375rem;
            animation: slideIn 0.2s ease-out;
            margin: 0.25rem 0;
        }

        #messages p[data-type="info"] {
            align-self: flex-end;
            background: var(--message-sent);
            color: white;
            border-top-right-radius: 0.25rem;
        }

        #messages p[data-type="alert"] {
            align-self: flex-start;
            background: var(--message-received);
            border: 1px solid var(--border-color);
            border-top-left-radius: 0.25rem;
        }

        @keyframes slideIn {
            from { opacity: 0; transform: translateY(8px); }
            to { opacity: 1; transform: translateY(0); }
        }

        #messages p strong {
            color: inherit;
            font-weight: 600;
            font-size: 0.75rem;
            display: block;
            margin-bottom: 0.25rem;
            opacity: 0.8;
        }

        .message-form {
            padding: 1rem;
            background: var(--bg-color);
            border-top: 1px solid var(--border-color);
            display: flex;
            gap: 0.75rem;
            align-items: center;
        }

        .input-group {
            display: flex;
            gap: 0.5rem;
            align-items: center;
            width: 100%;
        }

        .connection-form {
            display: flex;
            flex-direction: column;
            gap: 0.75rem;
        }

        label {
            font-size: 0.75rem;
            font-weight: 500;
            color: var(--text-color);
            opacity: 0.8;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        input, select {
            padding: 0.75rem;
            border: 1px solid var(--border-color);
            border-radius: 0.5rem;
            font-size: 0.9375rem;
            color: var(--text-color);
            background: var(--input-bg);
            transition: all 0.2s;
        }

        input:focus, select:focus {
            outline: none;
            border-color: var(--primary);
            box-shadow: 0 0 0 2px rgba(0, 0, 0, 0.1);
        }

        input[type="text"] {
            width: 100%;
        }

        select {
            padding-right: 2rem;
            background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='black' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
            background-repeat: no-repeat;
            background-position: right 0.5rem center;
            background-size: 1rem;
            appearance: none;
        }

        button {
            padding: 0.75rem 1.5rem;
            background: var(--primary);
            color: white;
            border: none;
            border-radius: 0.5rem;
            font-size: 0.9375rem;
            font-weight: 500;
            cursor: pointer;
            transition: background-color 0.2s;
            display: flex;
            align-items: center;
            gap: 0.5rem;
            letter-spacing: 0.5px;
        }

        button:hover {
            background: var(--hover-color);
        }

        .target-inputs {
            display: flex;
            gap: 0.5rem;
        }

        .target-inputs input {
            flex: 1;
        }

        @media (max-width: 768px) {
            .container {
                grid-template-columns: 1fr;
            }

            .sidebar {
                display: none;
            }

            .message-form {
                flex-direction: column;
            }

            .target-inputs {
                flex-direction: column;
            }

            button {
                width: 100%;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <aside class="sidebar">
            <header>
                <h1>Settings</h1>
            </header>
            <div class="connection-panel">
                <form class="connection-form">
                    <div>
                        <label for="userId">User ID</label>
                        <input type="text" id="userId" value="user123" />
                    </div>
                    <div>
                        <label for="groupId">Group ID</label>
                        <input type="text" id="groupId" value="demo-users" />
                    </div>
                    <button type="button" onclick="connect()">Connect</button>
                </form>
            </div>
        </aside>
        <main class="main-content">
            <header>
                <h1>Notifications</h1>
                <div id="status">Disconnected</div>
            </header>
            <div id="messages"></div>
            <footer class="message-form">
                <div class="input-group">
                    <input type="text" id="message" placeholder="Type your message..." />
                    <select id="type">
                        <option value="info">Info</option>
                        <option value="alert">Alert</option>
                    </select>
                </div>
                <div class="target-inputs">
                    <input type="text" id="targetId" placeholder="Target User" />
                    <input type="text" id="targetGroup" placeholder="Target Group" />
                </div>
                <button onclick="sendNotification()">Send</button>
            </footer>
        </main>
    </div>

    <script>
        let ws;

        function connect() {
            const userId = document.getElementById('userId').value;
            const groupId = document.getElementById('groupId').value;
            const statusEl = document.getElementById('status');
            const url = `ws://localhost:8080/ws?userId=${userId}&groupId=${groupId}`;
            
            if (ws) {
                ws.close();
            }

            ws = new WebSocket(url);

            ws.onopen = () => {
                statusEl.textContent = `Connected as ${userId}`;
                statusEl.classList.add('connected');
            };

            ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                const messages = document.getElementById('messages');
                const messageEl = document.createElement('p');
                messageEl.setAttribute('data-type', data.type);
                messageEl.innerHTML = `<strong>${data.type.toUpperCase()}</strong>${data.message}<br><small>ID: ${data.id}</small>`;
                messages.appendChild(messageEl);
                messages.scrollTop = messages.scrollHeight;
            };

            ws.onclose = () => {
                statusEl.textContent = 'Disconnected';
                statusEl.classList.remove('connected');
            };
        }

        function sendNotification() {
            const type = document.getElementById('type').value;
            const message = document.getElementById('message').value;
            const targetId = document.getElementById('targetId').value;
            const targetGroup = document.getElementById('targetGroup').value;

            if (!message.trim()) {
                return;
            }

            let url = `http://localhost:8080/notify/${type}/${encodeURIComponent(message)}`;
            const params = [];
            if (targetId) params.push(`targetId=${targetId}`);
            if (targetGroup) params.push(`groupId=${targetGroup}`);
            if (params.length > 0) url += `?${params.join('&')}`;

            fetch(url, { method: 'POST' })
                .then(() => {
                    document.getElementById('message').value = '';
                    document.getElementById('targetId').value = '';
                    document.getElementById('targetGroup').value = '';
                })
                .catch(err => console.error('Error sending notification:', err));
        }

        document.getElementById('message').addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                sendNotification();
            }
        });
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
 "time"

 "github.com/go-chi/chi/v5"
 "github.com/go-chi/chi/v5/middleware"
 "github.com/kenztech/notify"
 "github.com/kenztech/notify/models"
)

// InMemoryStore implements persistence.Store
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

// InMemoryBroker implements messaging.Broker
type InMemoryBroker struct {
 subs   map[string]chan []byte
 users  map[string]bool
 groups map[string]map[string]bool
 mutex  sync.RWMutex
}

func NewInMemoryBroker() *InMemoryBroker {
 return &InMemoryBroker{
  subs:   make(map[string]chan []byte),
  users:  make(map[string]bool),
  groups: make(map[string]map[string]bool),
 }
}

func (b *InMemoryBroker) Publish(channel string, message []byte) error {
 b.mutex.RLock()
 defer b.mutex.RUnlock()
 if ch, ok := b.subs[channel]; ok {
  select {
  case ch <- message:
  default:
  }
 }
 return nil
}

func (b *InMemoryBroker) Subscribe(channel string) (chan []byte, func(), error) {
 b.mutex.Lock()
 defer b.mutex.Unlock()
 ch := make(chan []byte, 256)
 b.subs[channel] = ch
 return ch, func() {
  b.mutex.Lock()
  defer b.mutex.Unlock()
  close(b.subs[channel])
  delete(b.subs, channel)
 }, nil
}

func (b *InMemoryBroker) TrackUser(userID string, groupIDs []string) error {
 b.mutex.Lock()
 defer b.mutex.Unlock()
 b.users[userID] = true
 for _, groupID := range groupIDs {
  if _, ok := b.groups[groupID]; !ok {
   b.groups[groupID] = make(map[string]bool)
  }
  b.groups[groupID][userID] = true
 }
 return nil
}

func (b *InMemoryBroker) UntrackUser(userID string, groupIDs []string) error {
 b.mutex.Lock()
 defer b.mutex.Unlock()
 delete(b.users, userID)
 for _, groupID := range groupIDs {
  if members, ok := b.groups[groupID]; ok {
   delete(members, userID)
   if len(members) == 0 {
    delete(b.groups, groupID)
   }
  }
 }
 return nil
}

func (b *InMemoryBroker) GetGroupMembers(groupID string) ([]string, error) {
 b.mutex.RLock()
 defer b.mutex.RUnlock()
 var members []string
 if group, ok := b.groups[groupID]; ok {
  for userID := range group {
   members = append(members, userID)
  }
 }
 return members, nil
}

func (b *InMemoryBroker) GetActiveUsers() ([]string, error) {
 b.mutex.RLock()
 defer b.mutex.RUnlock()
 var users []string
 for userID := range b.users {
  users = append(users, userID)
 }
 return users, nil
}

func main() {
 r := chi.NewRouter()
 r.Use(middleware.Logger)

 // Initialize store and broker
 store := NewInMemoryStore()
 broker := NewInMemoryBroker()

 // Configure notification system
 cfg := notify.Config{
  Store:  store,
  Broker: broker,
 }
 ns := notify.NewNotify(cfg)

 // Register notification routes
 ns.RegisterRoutes(r)

 // Serve static HTML frontend
 r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
 r.Get("/", func(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "static/index.html")
 })

 // Send a welcome notification after startup
 go func() {
  time.Sleep(2 * time.Second)
  req, _ := http.NewRequest("POST", "http://localhost:8080/notify/info/Welcome%20to%20Notify%20Demo!?groupId=demo-users", nil)
  http.DefaultClient.Do(req)
 }()

 // Start server
 log.Println("Server starting on :8080...")
 if err := http.ListenAndServe(":8080", r); err != nil {
  log.Fatalf("Server failed: %v", err)
 }
}

```

## Running the System

### 1. Start the Backend

```sh
go run main.go
```

### 2. Open the Frontend

Visit `http://localhost:8081` and connect with a user ID.

### 3. Send Notifications

Use the input fields and buttons to send notifications.

## Conclusion

This notification system allows users to receive real-time messages using WebSockets. The backend efficiently handles message storage and delivery using an in-memory implementation.
