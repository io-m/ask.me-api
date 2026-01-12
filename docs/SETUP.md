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
- `reacted` - users → messages (emoji reactions)

## Seed Data

The seeder **truncates all collections first** to ensure a clean state, then creates data aligned with the frontend mocks.

**Users (24):**

| ID | Username | Description |
|----|----------|-------------|
| `u-johndoe` | johndoe | Main logged-in user (default for testing) |
| `u-sandro` | sandro_r | Relationship question author |
| `u-maria` | maria_dev | Tech developer, imposter syndrome responder |
| `u-jordan` | jordantech | Career advice seeker |
| `u-taylor` | taylor_wanderlust | Travel enthusiast |
| `u-emma` | emma_wellness | Work-life balance author |
| `u-lucas` | lucasmartinez | 20s mistakes question author |
| `u-alex` | alex_innovator | Tech stack enthusiast |
| `u-olivia` | olivia_health | Health & wellness |
| `u-alice` | alice_wfh | Remote work expert |
| `u-bob` | bob_coder | Software developer |
| `u-david` | david_tech | Tech stack question author |
| `u-oliver` | oliver_senior | Senior developer |
| `u-sarah` | sarah_japan | Japan travel expert |
| `u-mike` | mike_adventures | Travel enthusiast |
| `u-tom` | tom_reader | Book club member |
| `u-lisa` | lisa_books | Book club member |
| `u-rachel` | rachel_lit | Book club member |
| `u-maya` | maya_coffee | Coffee vs tea author |
| `u-chris` | chris_beans | Coffee enthusiast |
| `u-anna` | anna_matcha | Tea enthusiast |
| `u-nina` | nina_fullstack | Full-stack developer |
| `u-james` | james_growth | Personal growth |
| `u-sophie` | sophie_balance | Life balance advocate |

**Posts (18):**

Posts johndoe can respond to:
- `p1` Sandro's relationship post
- `p2` Maria's imposter syndrome question
- `p3` Taylor's travel poll (poll)
- `p4` Jordan's 20-year-old advice question
- `p5` Taylor's destination poll (poll)
- `p6` Emma's work-life balance question
- `p7` Lucas's biggest 20s mistake question

Johndoe's own questions:
- `q-mine-1` Remote work focus question
- `q-mine-2` Imposter syndrome question
- `q-mine-3` Japan trip recommendations (group chat)
- `q-mine-4` Book club recommendations (group chat)

Group chat posts:
- `p-group-1` David's tech stack question (group)
- `p-group-2` Maya's coffee vs tea (group)

Plus additional posts for variety.

**Chats (12):**

Direct chats (johndoe responding to others):
- `c1` p1 - johndoe ↔ sandro
- `c4` p4 - johndoe ↔ jordan
- `c6` p6 - johndoe ↔ emma
- `c7` p7 - johndoe ↔ lucas

Direct chats (others responding to johndoe's questions):
- `chat-1` q-mine-1 - johndoe ↔ alice
- `chat-2` q-mine-1 - johndoe ↔ bob
- `chat-7` q-mine-2 - johndoe ↔ maria
- `chat-8` q-mine-2 - johndoe ↔ oliver

Group chats:
- `chat-group-1` q-mine-3 - Japan trip (johndoe, sarah, mike, emma)
- `chat-group-2` p-group-1 - Tech stack (david, johndoe, alex, nina)
- `chat-group-3` q-mine-4 - Book club (johndoe, lisa, tom, rachel)
- `chat-group-4` p-group-2 - Coffee vs Tea (maya, johndoe, chris, anna)

**Tags (10):**

career, tech, relationships, travel, wellness, personal-growth, productivity, books, food, remote-work

## Testing Endpoints

Use the `api.http` file with VS Code REST Client extension or similar.

Example:
```http
### Get user
GET http://localhost:8080/users/u-johndoe

### Get feed for authenticated user
GET http://localhost:8080/me/feed
X-User-ID: u-johndoe

### Get chats for authenticated user
GET http://localhost:8080/me/chats
X-User-ID: u-johndoe
```
