.DEFAULT_GOAL := help

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development targets
.PHONY: dev-up dev-down dev-kvs-up dev-kvs-down
dev-up: ## Start both services in development mode
	docker-compose -f docker-compose.dev.yaml up --build -d

dev-down: ## Stop both services in development mode
	docker-compose -f docker-compose.dev.yaml down

dev-kvs-up: ## Start only the key-value service in development mode
	docker-compose -f docker-compose.dev.yaml up --build -d kv-service

dev-kvs-down: ## Stop only the key-value service in development mode
	docker-compose -f docker-compose.dev.yaml stop kv-service

# Testing targets
.PHONY: test-kvs
test-kvs: ## Run all unit tests for the key-value service
	cd kv_service && go test ./...

test-client: ## Run all unit tests for the test client
	cd test_client && go test ./...

test: ## Run all unit tests for both services
	make test-kvs
	make test-client