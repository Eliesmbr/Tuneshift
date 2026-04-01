.PHONY: dev build docker run clean

# Run backend in development mode
dev:
	go run ./cmd/server

# Build Go binary
build:
	CGO_ENABLED=0 go build -o bin/tuneshift ./cmd/server

# Build Docker image
docker:
	docker compose build

# Run with Docker
run:
	docker compose up

# Run with Docker (detached)
run-d:
	docker compose up -d

# Stop Docker
stop:
	docker compose down

# Frontend development (inside web/)
frontend-dev:
	cd web && npm run dev

# Frontend build (inside web/)
frontend-build:
	cd web && npm run build

# Clean build artifacts
clean:
	rm -rf bin/ web/dist/ web/node_modules/
