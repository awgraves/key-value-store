.DEFAULT_GOAL := help

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development targets
.PHONY: dev-kvs-up dev-kvs-down
dev-kvs-up: ## Start the key-value service in development mode
	cd kv_service && docker-compose -f docker-compose.dev.yaml up --build -d

dev-kvs-down: ## Stop the key-value service in development mode
	cd kv_service && docker-compose -f docker-compose.dev.yaml down

# Testing targets
.PHONY: test-kvs
test-kvs: ## Run all unit tests for the key-value service
	cd kv_service && go test ./...