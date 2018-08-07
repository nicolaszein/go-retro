# COMMANDS

help:  ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

run: ## Run server
	@go run main.go

create-env: ## Create .env file
	@cp -n .env.sample .env

makemigration: ## Generate migration
	@buffalo db generate sql $(name)

migrate: ## Run migrations
	@buffalo db migrate up -e $(env)

migrate-down: ## Run migrations down
	@buffalo db migrate down -e $(env)

dependencies: ## Install dependencies
	@dep ensure

tests: ## Run tests
	@go test ./... -v

test-coverage: ## Run tests with coverage
	@go test ./... -cover


.PHONY: help
