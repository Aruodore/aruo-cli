SHELL := /bin/sh
.DEFAULT_GOAL := help
.PHONY: help install dev check test build security
help: ## Show supported tasks.
	@awk 'BEGIN {FS = ":.*## "; printf "Development tasks:\n\n"} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
install: ## Install dependencies and create the lockfile.
	npm install
dev: ## Start the local development server.
	npm run dev
check: ## Run formatting, linting, types, tests, and the production build.
	npm run check
test: ## Run tests once.
	npm test
build: ## Create the production artifact.
	npm run build
security: ## Query the npm advisory service.
	npm audit --audit-level=high
