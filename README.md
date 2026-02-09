# Notification System

Scalable notification system implemented in Go (Fiber + MySQL).

## Features

- **Notification Management**: Create, list, cancel, and view status of notifications.
- **Processing Engine**: Asynchronous worker pool with priority queue and rate limiting.
- **Delivery**: Integration with external webhook providers.
- **Observability**: Metrics and health check endpoints.

## Setup

1. **Prerequisites**: Docker and Docker Compose.
2. **Run**:
   ```bash
   docker-compose up --build
   ```

## Configuration

Environment variables in `docker-compose.yml`:
- `WEBHOOK_URL`: URL for the external provider (e.g., webhook.site).
- `MYSQL_*`: Database credentials.

## API Endpoints

### Notifications

- **Create**: `POST /notification`
  ```json
  {
    "recipient": "+905551234567",
    "channel": "sms",
    "content": "Hello World",
    "priority": "high"
  }
  ```
- **List**: `GET /notification?status=1&page=1`
- **Details**: `GET /notification/:id`
- **Cancel**: `PUT /notification/cancel/:id`

### Observability

- **Metrics**: `GET /metrics`
- **Health**: `GET /health`

## Architecture

- **Clean Architecture**: Domain, Ports, Adapters (Usecase, Repository, Handler).
- **Worker Pool**: Polls DB for pending tasks, respects priority and rate limits.