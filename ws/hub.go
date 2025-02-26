package ws

import (
	"log"
	"sync"

	"github.com/kenztech/notify/messaging"
	"github.com/kenztech/notify/models"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan models.Notification
	register   chan *Client
	unregister chan *Client
	broker     messaging.Broker
	mu         sync.RWMutex
}

func NewHub(broker messaging.Broker) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan models.Notification, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broker:     broker,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.broker.TrackUser(client.userID, client.groupIDs)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
				h.broker.UntrackUser(client.userID, client.groupIDs)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			data, err := message.Marshal()
			if err != nil {
				log.Printf("Failed to marshal notification: %v", err)
				continue
			}
			if message.TargetID != "" {
				h.broker.Publish("notify:user:"+message.TargetID, data)
			} else if len(message.GroupIDs) > 0 {
				for _, groupID := range message.GroupIDs {
					members, _ := h.broker.GetGroupMembers(groupID)
					for _, userID := range members {
						h.broker.Publish("notify:user:"+userID, data)
					}
				}
			} else {
				users, _ := h.broker.GetActiveUsers()
				for _, userID := range users {
					h.broker.Publish("notify:user:"+userID, data)
				}
			}
		}
	}
}
