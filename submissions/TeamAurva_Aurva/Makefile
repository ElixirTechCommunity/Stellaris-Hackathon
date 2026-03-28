.PHONY: help build run-control run-scanner run-dashboard docker-up docker-down test clean init-db deps

help:
	@echo "Aurva - DPDP Act 2023 Compliance Engine"
	@echo ""
	@echo "Available commands:"
	@echo "  make docker-up       - Start PostgreSQL via Docker Compose"
	@echo "  make docker-down     - Stop Docker services"
	@echo "  make run-control     - Run control plane"
	@echo "  make run-dashboard   - Run Next.js dashboard"
	@echo "  make build           - Build all Go binaries"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Clean build artifacts"

docker-up:
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	@echo "PostgreSQL is ready!"

docker-down:
	docker-compose down

run-control:
	@echo "Starting control plane..."
	cd control-plane && go run main.go

run-dashboard:
	@echo "Starting dashboard..."
	cd dashboard && npm run dev

build:
	@echo "Building control plane..."
	cd control-plane && go build -o ../bin/control-plane main.go
	@echo "Building cloud scanner..."
	cd cloud-scanner && go build -o ../bin/cloud-scanner main.go
	@echo "Build complete! Binaries in ./bin/"

test:
	@echo "Running Go tests..."
	go test ./... -v
	@echo "Testing PII classifier..."
	cd shared/pii && go test -v

clean:
	rm -rf bin/
	rm -rf dashboard/.next/
	rm -rf dashboard/node_modules/
	go clean

init-db:
	@echo "Initializing database schema..."
	PGPASSWORD=aurva_dev_password psql -h localhost -U aurva -d aurva -f schema.sql
	@echo "Database initialized!"

deps:
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing Node.js dependencies..."
	cd dashboard && npm install
	@echo "Dependencies installed!"
