package domain

import "time"

// Notification represents a notification record with its delivery state.
type Notification struct {
	ID        string     `json:"id"`
	Channel   string     `json:"channel"`
	Recipient string     `json:"recipient"`
	Subject   string     `json:"subject"`
	Body      string     `json:"body"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
}

// CreateRequest holds the fields required to create a new notification.
type CreateRequest struct {
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}
