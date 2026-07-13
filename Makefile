ifneq ($(wildcard .env),)
include .env
export
endif

PORT ?= 8080
POSTGRES_USER ?= admin
POSTGRES_PASSWORD ?= admin
POSTGRES_DB ?= orders_service
POSTGRES_PORT ?= 5432
POSTGRES_HOST ?= localhost
POSTGRES_SSLMODE ?= disable
DATABASE_URL ?= "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)"


.PHONY: help up down db-up wait-db migrate-up db-down backend


help:
	@echo "Available commands:"
	@echo "  make help - Show this help message"
	@echo "  make up - Start the database, wait for it to be ready, run migrations, and start the API"
	@echo "  make down - Drop the database container"
	@echo "  make db-up - Start the database container"
	@echo "  make wait-db - Wait for the database to be ready"
	@echo "  make migrate-up - Run migrations"
	@echo "  make db-down - Stop the database container"
	@echo "  make backend - Start the API"

up: db-up wait-db migrate-up backend

down:
	@echo "Dropping database container..."
	@docker compose down

db-up:
	@echo "Starting the PostgreSQL container in the background..."
	@docker compose up -d

wait-db:
	@echo "Waiting for the database to be ready..."
	@until docker exec -t order_service_db pg_isready -u $(POSTGRES_USER) -d $(POSTGRES_DB) > /dev/null 2>&1; do \
		i=$$((i+1)); \
		if [ $$i -gt 30 ]; then \
			echo "\nError: The database did not become ready in time."; \
			exit 1; \
		fi; \
		echo -n "."; sleep 1; \
	done
	@echo "\nDatabase is ready!"

migrate-up:
	@command -v migrate >/dev/null 2>&1 || { \
		@echo >&2 "The command 'migrate' is not installed. \
		Please install it by running 'go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest'."; \
		exit 1; \
	}
	@echo "Executing migrations..."
	@migrate -database "$(DATABASE_URL)" -path migrations up

migrate-down:
	@command -v migrate >/dev/null 2>&1 || { \
		@echo >&2 "The command 'migrate' is not installed. \
		Please install it by running 'go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest'."; \
		exit 1; \
	}
	@echo "Reverting migrations..."
	@migrate -database "$(DATABASE_URL)" -path migrations down

db-down:
	@echo "Stopping the PostgreSQL container..."
	@docker compose stop order_service_db	

backend:
	@echo "Starting the backend API..."
	@go run ./cmd/api/main.go