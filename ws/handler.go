package ws

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/kenztech/notify/messaging"
	"github.com/kenztech/notify/models"
	"github.com/kenztech/notify/persistence"
)

type Handler struct {
	hub      *Hub
	store    persistence.Store
	broker   messaging.Broker
	upgrader websocket.Upgrader
}

func NewHandler(hub *Hub, store persistence.Store, broker messaging.Broker) *Handler {
	return &Handler{
		hub:    hub,
		store:  store,
		broker: broker,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Handler) ServeWs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId query parameter required", http.StatusBadRequest)
		return
	}
	groupIDs := r.URL.Query()["groupId"]

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	client := NewClient(h.hub, conn, userID, groupIDs, h.broker)
	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) SendNotification(w http.ResponseWriter, r *http.Request) {
	notification := models.NewNotification(
		time.Now().String(),
		chi.URLParam(r, "type"),
		chi.URLParam(r, "message"),
		r.URL.Query().Get("targetId"),
		r.URL.Query()["groupId"],
		time.Now().Unix(),
	)

	if err := h.store.SaveNotification(notification); err != nil {
		http.Error(w, "Failed to save notification", http.StatusInternalServerError)
		return
	}

	h.hub.broadcast <- notification
	w.WriteHeader(http.StatusAccepted)
}
