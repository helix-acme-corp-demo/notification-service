# Design: Rate Limiting for Outbound Notification Sends

## Overview

Add a `RateLimitedSender` decorator in the `sender` package that wraps the existing `Sender` interface with per-channel rate limiting using the `github.com/helix-acme-corp-demo/ratelimit` library. Wired up in `main.go` — no other packages change.

## Architecture

### Decorator Pattern

```
NotificationHandler
       │
       ▼
RateLimitedSender   ← new, wraps:
       │
       ▼
  logSender (or any Sender)
```

The `RateLimitedSender` sits between the handler and the real sender. It calls `limiter.Allow(ctx, n.Channel)` before delegating to the inner sender. If the decision is denied, it returns an error immediately (with the `RetryAfter` duration in the message) without calling the inner sender.

### New File: `internal/sender/ratelimited.go`

```go
type RateLimitedSender struct {
    inner   Sender
    limiter ratelimit.Limiter
}

func NewRateLimited(inner Sender, limiter ratelimit.Limiter) Sender { ... }

func (s *RateLimitedSender) Send(ctx context.Context, n *domain.Notification) error {
    decision := s.limiter.Allow(ctx, n.Channel)
    if !decision.Allowed {
        return fmt.Errorf("rate limit exceeded for channel %q, retry after %s", n.Channel, decision.RetryAfter)
    }
    return s.inner.Send(ctx, n)
}
```

### Wiring in `main.go`

```go
limiter := ratelimit.TokenBucket(
    ratelimit.WithRate(100, time.Minute),
    ratelimit.WithBurst(10),
)
notifSender := sender.NewRateLimited(sender.NewLog(logger), limiter)
```

The `go.mod` gains one new dependency: `github.com/helix-acme-corp-demo/ratelimit`.

## Key Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Limiter algorithm | `TokenBucket` | Tolerates short bursts (e.g. a batch of notifications at startup) while smoothing steady-state throughput. `SlidingWindow` is available but stricter. |
| Rate key | `n.Channel` | Channels (email, sms, webhook) map to distinct downstream providers with their own limits. Per-recipient limiting is out of scope. |
| Error handling | Return error from `Send` | The existing `retryx` retry loop in `NotificationHandler.Create()` already handles send errors and marks status `failed`. No new error path needed. |
| State | In-memory only | Matches the existing in-memory `Store`. Cross-restart persistence is out of scope. |
| Config location | `main.go` hardcoded constants | Simple. No config file or env-var parsing added unless the operator asks for it. |

## Codebase Notes for Implementors

- **`sender.Sender` interface** is defined in `internal/sender/sender.go` — just `Send(ctx, *domain.Notification) error`.
- **Retry logic** lives in `handler/notification.go` via `retryx.Do` (3 attempts, 100ms base delay). Rate limit errors will be retried — this is acceptable since `RetryAfter` is usually sub-second for token bucket.
- **`ratelimit` package** exposes `Limiter` interface, `TokenBucket(opts...)`, `SlidingWindow(opts...)`, `WithRate(count, window)`, `WithBurst(n)`, and a `Decision{Allowed, Remaining, RetryAfter}` struct.
- **Import path** for the ratelimit library: `github.com/helix-acme-corp-demo/ratelimit`.
- **No changes needed** to `domain`, `store`, `handler`, or `health` packages.