package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kenztech/notify/messaging"
	"github.com/kenztech/notify/models"
)

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan models.Notification
	userID    string
	groupIDs  []string
	broker    messaging.Broker
	closeChan chan struct{}
}

// NewClient creates a new client instance, accepting groupIDs explicitly
func NewClient(hub *Hub, conn *websocket.Conn, userID string, groupIDs []string, broker messaging.Broker) *Client {
	return &Client{
		hub:       hub,
		conn:      conn,
		send:      make(chan models.Notification, 256),
		userID:    userID,
		groupIDs:  groupIDs,
		broker:    broker,
		closeChan: make(chan struct{}),
	}
}
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		close(c.closeChan)
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	msgChan, cleanup, err := c.broker.Subscribe("notify:user:" + c.userID)
	if err != nil {
		log.Printf("Failed to subscribe for user %s: %v", c.userID, err)
		return
	}
	defer cleanup()

	go func() {
		for {
			select {
			case msg, ok := <-msgChan:
				if !ok {
					return
				}
				var n models.Notification
				if err := n.Unmarshal(msg); err == nil {
					c.send <- n
				} else {
					log.Printf("Failed to unmarshal notification for user %s: %v", c.userID, err)
				}
			case <-c.closeChan:
				return
			}
		}
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error for user %s: %v", c.userID, err)
			}
			break
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("write error for user %s: %v", c.userID, err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
