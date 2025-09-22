.DEFAULT_GOAL := help

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development targets
.PHONY: dev-up dev-down dev-kvs-up dev-kvs-down dev-logs
dev-up: ## Start both services in development mode
	docker-compose -f docker-compose.dev.yaml up --build -d

dev-down: ## Stop both services in development mode
	docker-compose -f docker-compose.dev.yaml down

dev-kvs-up: ## Start only the key-value service in development mode
	docker-compose -f docker-compose.dev.yaml up --build -d kv-service

dev-kvs-down: ## Stop only the key-value service in development mode
	docker-compose -f docker-compose.dev.yaml stop kv-service

dev-logs: ## View logs from development services
	docker-compose -f docker-compose.dev.yaml logs -f

# Production targets
.PHONY: prod-up prod-down prod-build prod-logs
prod-up: ## Start both services in production mode
	docker-compose up --build -d

prod-down: ## Stop both services in production mode
	docker-compose down

prod-build: ## Build production images without starting services
	docker-compose build

prod-logs: ## View logs from production services
	docker-compose logs -f

# Testing targets
.PHONY: test-kvs test-client test build-test-image
build-test-image: ## Build test image with pre-cached dependencies
	@echo "Building test image..."
	docker build -f Dockerfile.test -t kv-test-runner .

test-kvs: ## Run all unit tests for the key-value service
	@docker image inspect kv-test-runner >/dev/null 2>&1 || make build-test-image
	docker run --rm -v $(PWD)/kv_service:/app kv-test-runner

test-client: ## Run all unit tests for the test client
	@docker image inspect kv-test-runner >/dev/null 2>&1 || make build-test-image
	docker run --rm -v $(PWD)/test_client:/app kv-test-runner

test: ## Run all unit tests for both services
	make test-kvs
	make test-client