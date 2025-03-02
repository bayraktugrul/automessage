<pre><code>                _        __  __                                
     /\        | |      |  \/  |                               
    /  \  _   _| |_ ___ | \  / | ___  ___ ___  __ _  __ _  ___ 
   / /\ \| | | | __/ _ \| |\/| |/ _ \/ __/ __|/ _` |/ _` |/ _ \
  / ____ \ |_| | || (_) | |  | |  __/\__ \__ \ (_| | (_| |  __/
 /_/    \_\__,_|\__\___/|_|  |_|\___||___/___/\__,_|\__, |\___|
                                                     __/ |     
                                                    |___/
</code></pre>

[![][workflow-badge]][workflow-actions]

AutoMessage is automatic message sending system designed to automatically send messages at specified intervals. It pushes every Y messages in X period seconds, making it ideal for scheduled notifications, alerts, and automated communication workflows.

## Installation

### Prerequisites

- Go 1.23.3 or higher
- Docker and Docker Compose (for containerized deployment & testing)

### Local Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/bayraktugrul/automessage.git
   cd automessage
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

### Dockerize

For containerized deployment or to test on local, Docker Compose can be used:

```bash
# Build and start containers
docker compose up -d

# Rebuild containers (for code changes)
docker compose down -v && docker compose build --no-cache && docker compose up
```

### Swagger
```
http://localhost:8080/swagger/index.html
```

## Configuration

Configuration can be done through environment variables or a configuration file:

### Environment Variables

```
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=automsg
      - DB_SSL_MODE=disable
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=
      - WEBHOOK_URL=https://webhook.site/your-id
      - PORT=8080
      - MESSAGE_INITIAL_BATCH_SIZE=10
      - MESSAGE_PERIODIC_BATCH_SIZE=2
      - MESSAGE_INTERVAL_SECONDS=120
```
Configuration values can be changed on docker-compose.yml while testing.

## Usage

### API Endpoints

- `PUT /send`:  quaryParam: operation, values: START,STOP
- `GET /messages`: List sent messages

### Example requests

```
curl -X 'GET' \
  'http://localhost:8080/messages?page=1&pageSize=10' \
  -H 'accept: application/json'
```

```
curl -X 'PUT' \
  'http://localhost:8080/send' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "operation": "START"
}'
```
## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Example output

<p align="center">
<img src="https://github.com/user-attachments/assets/2c24d41e-7871-4fb7-b284-717094d10db8">
</p>

[workflow-actions]: https://github.com/bayraktugrul/automessage/actions
[workflow-badge]: https://github.com/bayraktugrul/automessage/workflows/build/badge.svg
