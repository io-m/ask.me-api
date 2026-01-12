# ask.me Backend Setup Guide

## Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Make

## Quick Start

```bash
# 1. Start ArangoDB and seed data
make setup

# 2. Run the API server
make run
```

The API will be available at `http://localhost:8080`

ArangoDB UI: `http://localhost:8529` (user: `root`, password: `rootpassword`)

## Available Commands

```bash
make docker-up    # Start ArangoDB container
make docker-down  # Stop ArangoDB container
make docker-logs  # View ArangoDB logs
make db-setup     # Create database and collections
make seed         # Seed mock data
make run          # Run API server
make build        # Build binary to bin/api
make test         # Run tests
make setup        # Full setup (docker + db + seed)
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | API server port |
| `ARANGO_ENDPOINT` | `http://localhost:8529` | ArangoDB endpoint |
| `ARANGO_DATABASE` | (required) | Database name |
| `ARANGO_USERNAME` | (required) | ArangoDB username |
| `ARANGO_PASSWORD` | (required) | ArangoDB password |

## Project Structure

```
api/
├── cmd/
│   ├── api/           # Main API server
│   │   ├── main.go    # Entry point
│   │   ├── app.go     # DI wiring
│   │   └── routes.go  # Route registration
│   └── seed/          # Database seeder
│       └── main.go
├── internal/          # Private application code
│   ├── config/        # Configuration
│   ├── domain/        # Shared types and errors
│   ├── user/          # User feature module
│   ├── post/          # Post feature module
│   ├── chat/          # Chat feature module
│   ├── feed/          # Feed feature module
│   └── tag/           # Tag feature module
├── pkg/               # Public packages
│   ├── arango/        # ArangoDB client wrapper
│   └── httputil/      # HTTP utilities
├── docs/              # Documentation
├── docker-compose.yml # ArangoDB container
├── Makefile           # Build commands
├── api.http           # HTTP client file for testing
└── go.mod
```

## Feature Module Structure

Each feature follows this pattern:

```
internal/user/
├── interfaces.go   # Repository, Service, Handler interfaces
├── entity.go       # Domain types (User, CreateUserRequest, etc.)
├── aql.go          # AQL queries as constants
├── repository.go   # Repository implementation
├── service.go      # Business logic implementation
└── handler.go      # HTTP handler implementation
```

## Database Collections

### Document Collections
- `users` - User profiles
- `posts` - Posts (text and polls)
- `tags` - Canonical tags
- `chats` - Chat containers
- `messages` - Chat messages

### Edge Collections
- `created` - users → posts (authorship)
- `responded` - users → posts (responses)
- `post_has_tag` - posts → tags
- `follows` - users → users
- `participates_in` - users → chats
- `tagged` - posts → users
- `voted` - users → posts (poll votes)

## Seed Data

The seeder creates:

**Users (5):**
- `u1` alex_dev - tech/career enthusiast
- `u2` maria_chen - health/lifestyle
- `u3` john_doe - finance/career
- `u4` sarah_k - relationships/fun
- `u5` dev_master - tech/education

**Posts (7):**
- `p1` Frontend→Backend question (career, serious)
- `p2` Morning routine question (lifestyle, casual)
- `p3` Framework poll - React/Vue/Angular/Svelte (tech)
- `p4` Remote job success story (career)
- `p5` Partner finances question (relationships, serious)
- `p6` Learning Go question (education)
- `p7` Sleep hours poll (health)

**Chats (3):**
- `c1` about p1: alex_dev ↔ john_doe
- `c2` about p2: maria_chen ↔ dev_master
- `c3` about p4: john_doe ↔ alex_dev

**Tags (10):**
career-change, backend-dev, frontend, remote-work, productivity, relationships, health-tips, investing, frameworks, learning

## Testing Endpoints

Use the `api.http` file with VS Code REST Client extension or similar.

Example:
```http
### Get user
GET http://localhost:8080/users/u1

### Get feed for user
GET http://localhost:8080/feed?userId=u1&limit=10
```
