# HTTP Example

A REST API for event management built with Go, PostgreSQL, and Docker.

## Quick Start

```bash
# 1. Start PostgreSQL
make db-up

# 2. Dependencies
make dependencies

# 2. Copy environment configuration
cp .env.local .env

# 3. Run migrations
make migrate

# 4. Run tests
make test

# 5. Start application
make run

# 6. Postman collection
Import the postman collection: Events.postman_collection.json
```



The server will start on `http://localhost:8080`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/events` | Create new event |
| GET    | `/events` | List all events |
| GET    | `/events/{id}` | Get event by ID |
| PUT    | `/events/{id}` | Update event |
| DELETE | `/events/{id}` | Delete event |

### Example Request

```bash
# Create event
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Go Conference",
    "description": "A conference about Go programming",
    "start_time": "2025-08-22T10:00:00Z",
    "end_time": "2025-08-22T12:00:00Z"
  }'

# List events
curl http://localhost:8080/events
```

## Database

- Server: `postgres`
- Username: `postgres`
- Password: `postgres123`
- Database: `taller_challenge`

## Commands

```bash
make help      # Show available commands
make run       # Run the application  
make test      # Run tests
make db-up     # Start PostgreSQL container
make db-down   # Stop PostgreSQL container
make migrate   # Run database migrations
```

## Project Structure

```
taller_challenge/
├── main.go                     # Application entry point
├── Makefile                    # Basic commands
├── docker-compose.yml          # PostgreSQL
├── migrations/                 # Database migrations
│   └── 001_create_events_table.sql
├── api/
│   └── eventController.go      # HTTP handlers
└── internal/
    ├── config.go               # Database connection
    ├── db.go                   # Repository implementation
    └── interfaces.go           # Repository interface
```

## Configuration

Create `.env` file:

```bash
# Database Configuration
DATABASE_URL=postgres://postgres:postgres123@localhost:5432/taller_challenge?sslmode=disable

```
