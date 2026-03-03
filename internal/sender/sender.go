package sender

import (
	"context"

	"github.com/helix-acme-corp-demo/logpipe"
	"github.com/helix-acme-corp-demo/notification-service/internal/domain"
)

// Sender delivers a notification through its configured channel.
type Sender interface {
	Send(ctx context.Context, n *domain.Notification) error
}

type logSender struct {
	logger logpipe.Logger
}

// NewLog returns a Sender that logs notification details instead of
// performing real delivery.
func NewLog(logger logpipe.Logger) Sender {
	return &logSender{logger: logger}
}

func (s *logSender) Send(_ context.Context, n *domain.Notification) error {
	s.logger.Info("sending notification",
		logpipe.String("id", n.ID),
		logpipe.String("channel", n.Channel),
		logpipe.String("recipient", n.Recipient),
		logpipe.String("subject", n.Subject),
	)
	return nil
}
