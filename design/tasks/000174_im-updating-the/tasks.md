# Implementation Tasks

- [ ] Add `github.com/helix-acme-corp-demo/ratelimit` to `go.mod` and `go.sum` via `go get`
- [ ] Create `internal/sender/ratelimited.go` implementing `RateLimitedSender` that wraps `Sender` and calls `limiter.Allow(ctx, n.Channel)` before delegating
- [ ] Return a descriptive error from `RateLimitedSender.Send` when the decision is denied, including the `RetryAfter` duration
- [ ] Export a `NewRateLimited(inner Sender, limiter ratelimit.Limiter) Sender` constructor
- [ ] Update `cmd/server/main.go` to instantiate a `ratelimit.TokenBucket` limiter (100 req/min, burst 10) and wrap `sender.NewLog` with `sender.NewRateLimited`
- [ ] Write unit tests in `internal/sender/ratelimited_test.go` covering: allowed request delegates to inner sender, denied request returns error and does not call inner sender, separate channels have independent buckets