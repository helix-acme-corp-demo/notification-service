# Design for Rate Limiting on Outbound API Calls

## Architecture

The notification service uses a modular design with a `Sender` interface responsible for delivering notifications. To add rate limiting, integrate the ratelimit library (from the `/home/retro/work/ratelimit` repo) into the sender implementation. This library provides a token bucket algorithm for rate limiting outbound HTTP calls.

Modify the sender (e.g., for email/SMS providers) to wrap API requests with rate limiting. For example, a `RateLimitedSender` struct composes an existing `Sender` and uses the ratelimit package to throttle calls before delegating to the underlying sender.

## Key Decisions

- **Location for Rate Limiting:** Apply in the `Sender.Send` method just before making HTTP requests to external providers, ensuring no internal operations are affected.
- **Rate Limiting Mechanism:** Use token bucket from the existing ratelimit library, as it's already available and battle-tested.
- **Configuration:** Read rate (requests per time unit) and burst amount from environment variables (e.g., `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`) for runtime configurability without code changes.
- **Fallback Behavior:** If rate limit is exceeded, wait for available tokens (synchronous blocking); if burst is full, drop or delay the request gracefully.
- **Thread Safety:** The ratelimit library handles concurrency for Go's goroutine-based service.

Rationale: This matches the codebase's Go patterns and leverages an existing library instead of reinventing. Keeps implementation simple and focused on outbound calls.
