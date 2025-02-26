package models

import (
	"encoding/json"
	"time"
)

type Notification struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	Message   string   `json:"message"`
	TargetID  string   `json:"targetId,omitempty"`
	GroupIDs  []string `json:"groupIds,omitempty"`
	Timestamp int64    `json:"timestamp"`
}

func NewNotification(id, typ, message, targetID string, groupIDs []string, timestamp int64) Notification {
	if id == "" {
		id = time.Now().String()
	}
	return Notification{
		ID:        id,
		Type:      typ,
		Message:   message,
		TargetID:  targetID,
		GroupIDs:  groupIDs,
		Timestamp: timestamp,
	}
}

// Marshal serializes the Notification to JSON
func (n Notification) Marshal() ([]byte, error) {
	return json.Marshal(n)
}

// Unmarshal deserializes JSON into a Notification
func (n *Notification) Unmarshal(data []byte) error {
	return json.Unmarshal(data, n)
}
