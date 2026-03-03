package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/helix-acme-corp-demo/envelope"
	"github.com/helix-acme-corp-demo/logpipe"
	"github.com/helix-acme-corp-demo/retryx"

	"github.com/helix-acme-corp-demo/notification-service/internal/domain"
	"github.com/helix-acme-corp-demo/notification-service/internal/sender"
	"github.com/helix-acme-corp-demo/notification-service/internal/store"
)

// NotificationHandler handles notification CRUD and delivery operations.
type NotificationHandler struct {
	store  *store.Store
	sender sender.Sender
	logger logpipe.Logger
}

// New creates a NotificationHandler with the given dependencies.
func New(store *store.Store, sender sender.Sender, logger logpipe.Logger) *NotificationHandler {
	return &NotificationHandler{
		store:  store,
		sender: sender,
		logger: logger,
	}
}

// Create returns an HTTP handler that accepts a notification creation request,
// persists the notification, and attempts delivery with automatic retry.
func (h *NotificationHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			envelope.Write(w, envelope.BadRequest("invalid_json", "could not decode request body"))
			return
		}

		now := time.Now()
		n := &domain.Notification{
			ID:        logpipe.GenerateID(),
			Channel:   req.Channel,
			Recipient: req.Recipient,
			Subject:   req.Subject,
			Body:      req.Body,
			Status:    "pending",
			CreatedAt: now,
		}

		h.store.Save(n)

		ctx := r.Context()
		err := retryx.Do(ctx, func(ctx context.Context) error {
			return h.sender.Send(ctx, n)
		}, retryx.WithMaxAttempts(3), retryx.WithBaseDelay(100*time.Millisecond))

		if err != nil {
			n.Status = "failed"
			h.logger.Error("notification delivery failed",
				logpipe.String("id", n.ID),
				logpipe.Err(err),
			)
		} else {
			sentAt := time.Now()
			n.Status = "sent"
			n.SentAt = &sentAt
		}

		h.store.Save(n)

		envelope.Write(w, envelope.Created(n))
	}
}

// Get returns an HTTP handler that retrieves a single notification by ID.
func (h *NotificationHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		n, ok := h.store.Find(id)
		if !ok {
			envelope.Write(w, envelope.NotFound("notification not found"))
			return
		}
		envelope.Write(w, envelope.OK(n))
	}
}

// List returns an HTTP handler that retrieves all notifications.
func (h *NotificationHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		envelope.Write(w, envelope.OK(h.store.All()))
	}
}
