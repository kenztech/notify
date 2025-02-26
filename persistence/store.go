package persistence

import "github.com/kenztech/notify/models"

type Store interface {
	SaveNotification(n models.Notification) error
	GetNotifications(targetID string, groupIDs []string) ([]models.Notification, error)
}
