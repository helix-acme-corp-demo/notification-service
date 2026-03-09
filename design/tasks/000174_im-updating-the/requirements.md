# Requirements: Rate Limiting for Outbound API Calls

## Background

The notification service makes outbound delivery calls via the `Sender` interface. Currently there is no rate limiting on these calls, which risks overwhelming downstream channels (email providers, webhook endpoints, etc.). A shared `ratelimit` library already exists at `github.com/helix-acme-corp-demo/ratelimit` with both `TokenBucket` and `SlidingWindow` implementations.

## User Stories

### US-1: Per-channel rate limiting
As an operator, I want outbound notification sends to be rate-limited per channel (e.g. `email`, `sms`, `webhook`), so that we don't flood downstream providers and trigger throttling or bans.

**Acceptance criteria:**
- Each distinct `channel` value has its own rate limit bucket.
- When the limit is exceeded, the notification is marked `failed` and the error is logged with the `RetryAfter` duration.
- Allowed requests proceed to the underlying `Sender` without observable latency overhead.

### US-2: Configurable rate
As a developer, I want the rate limit (requests per window) and burst size to be configurable at startup, so that limits can be tuned per environment without redeployment of new code.

**Acceptance criteria:**
- Rate and burst values can be set when wiring up the sender in `main.go`.
- Defaults (100 req/min, burst 10) are used when no explicit configuration is provided.

### US-3: Transparent `Sender` wrapper
As a developer, I want rate limiting to be added as a decorator around the existing `Sender` interface, so that existing code in `NotificationHandler` does not need to change.

**Acceptance criteria:**
- A new `RateLimitedSender` type satisfies the `sender.Sender` interface.
- It wraps any existing `Sender` implementation.
- `NotificationHandler`, `Store`, and `domain` packages require zero modifications.

## Out of Scope

- Persistent (cross-restart) rate limit state.
- Per-recipient rate limiting.
- Dynamic runtime reconfiguration of limits.