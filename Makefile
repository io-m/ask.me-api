.PHONY: docker-up docker-down docker-logs db-setup seed run build test

# Docker commands
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f arangodb

# Database setup (creates database and collections)
db-setup:
	@echo "Setting up ArangoDB database..."
	@curl -u root:rootpassword -X POST http://localhost:8529/_api/database \
		-H "Content-Type: application/json" \
		-d '{"name": "askme"}' || true
	@echo "\nCreating collections..."
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "users"}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "posts"}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "tags"}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "chats"}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "messages"}' || true
	@echo "\nCreating edge collections..."
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "created", "type": 3}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "responded", "type": 3}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "post_has_tag", "type": 3}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "follows", "type": 3}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "participates_in", "type": 3}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "tagged", "type": 3}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "voted", "type": 3}' || true
	@curl -u root:rootpassword -X POST http://localhost:8529/_db/askme/_api/collection \
		-H "Content-Type: application/json" -d '{"name": "reacted", "type": 3}' || true
	@echo "\nDatabase setup complete!"

# Seed mock data
seed:
	ARANGO_DATABASE=askme ARANGO_USERNAME=root ARANGO_PASSWORD=rootpassword go run ./cmd/seed

# Run the API server
run:
	ARANGO_DATABASE=askme ARANGO_USERNAME=root ARANGO_PASSWORD=rootpassword go run ./cmd/api

# Build the binary
build:
	go build -o bin/api ./cmd/api

# Run tests
test:
	go test -v ./...

# Full setup: start docker, wait, setup db, seed
setup: docker-up
	@echo "Waiting for ArangoDB to start..."
	@sleep 5
	@$(MAKE) db-setup
	@$(MAKE) seed