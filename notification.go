package notify

import (
	"github.com/go-chi/chi/v5"
	"github.com/kenztech/notify/messaging"
	"github.com/kenztech/notify/persistence"
	"github.com/kenztech/notify/ws"
)

// System encapsulates the notification system
type Notify struct {
	Hub     *ws.Hub
	Handler *ws.Handler
	Store   persistence.Store
	Broker  messaging.Broker
}

// Config holds configuration for the notification system
type Config struct {
	Store  persistence.Store
	Broker messaging.Broker
}

// NewSystem initializes a new notification system
func NewNotify(cfg Config) *Notify {
	hub := ws.NewHub(cfg.Broker)
	handler := ws.NewHandler(hub, cfg.Store, cfg.Broker)

	go hub.Run()

	return &Notify{
		Hub:     hub,
		Handler: handler,
		Store:   cfg.Store,
		Broker:  cfg.Broker,
	}
}

// RegisterRoutes adds notification routes to a Chi router
func (s *Notify) RegisterRoutes(r chi.Router) {
	r.Get("/ws", s.Handler.ServeWs)
	r.Post("/notify/{type}/{message}", s.Handler.SendNotification)
}
