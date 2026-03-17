# Architecture

The test task requires a simple notification sender using the existing notification-service architecture. It will follow the clean architecture pattern with domain, handler, sender, and store layers as seen in the codebase.

Key components:
- Handler: REST API endpoints for test notifications
- Sender: Mock sender for testing (no actual emails/Slack sends)
- Store: In-memory storage for notifications (extend existing store module)

# Key Decisions

- Use in-memory storage instead of external DB for simplicity, as this is a test implementation.
- Add a new /test endpoint to the server, reusing the existing health check pattern.
- No external libraries needed; leverage Go standard library for mocking.
- Rationale: Minimal complexity for a test task, avoids over-engineering.

# Dependencies

None new; all logic built on existing internal modules.