# notification-service

Internal notification microservice for Acme Corp. Sends notifications via email, Slack, and webhooks with automatic retry.

## Architecture

This service uses several internal Acme libraries:
- **[retryx](https://github.com/helix-acme-corp-demo/retryx)** — Automatic retry with exponential backoff for delivery
- **[logpipe](https://github.com/helix-acme-corp-demo/logpipe)** — Structured JSON logging with correlation IDs
- **[envelope](https://github.com/helix-acme-corp-demo/envelope)** — Standardized API response formatting

## API Endpoints

### Health Check
```
GET /health
```

### Create Notification
```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{"channel":"email","recipient":"user@example.com","subject":"Hello","body":"World"}'
```

### Get Notification
```bash
curl http://localhost:8080/notifications/{id}
```

### List Notifications
```bash
curl http://localhost:8080/notifications
```

## Running

```bash
go run ./cmd/server
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | HTTP server port |

## License

MIT
