
.PHONY: help run test db-up db-down migrate

help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Run the application
	@echo "Running application..."
	go run main.go

dependencies: 
	@echo "Adding dependencies..."
	go mod tidy
	go mod vendor

test: ## Run tests
	@echo "Running tests..."
	go test ./internal -v 

db-up: ## Start PostgreSQL container
	@echo "Starting PostgreSQL..."
	docker-compose up -d postgres
	@echo "PostgreSQL is ready!"
	@echo "Connection details:"
	@echo "  Host: localhost:5432"
	@echo "  Database: taller_challenge"
	@echo "  User: postgres"
	@echo "  Password: postgres123"

db-down: ## Stop PostgreSQL container
	@echo "Stopping PostgreSQL..."
	docker-compose down
	@echo "PostgreSQL stopped"

migrate: ## Run database migrations
	@echo "Running migrations..."
	@if [ -d "./migrations" ]; then \
		for migration in ./migrations/*.sql; do \
			if [ -f "$$migration" ]; then \
				echo "Running: $$(basename "$$migration")"; \
				docker exec -i taller_challenge-postgres-1 psql -U postgres -d taller_challenge < "$$migration"; \
				if [ $$? -eq 0 ]; then \
					echo "$$(basename "$$migration") completed successfully"; \
				else \
					echo "$$(basename "$$migration") failed"; \
				fi; \
			fi; \
		done; \
		echo "All migrations completed"; \
	else \
		echo "No migrations directory found"; \
	fi
