SHELL := /bin/sh

.DEFAULT_GOAL := help

.PHONY: help install dev check test build db-up db-down db-migrate security

help: ## Show supported development tasks.
	@awk 'BEGIN {FS = ":.*## "; printf "Development tasks:\n\n"} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install the pinned direct dependency graph.
	npm install

dev: ## Start the development server.
	npm run dev

check: ## Run the complete deterministic local quality gate.
	npm run check

test: ## Run unit tests once.
	npm test

build: ## Build the production server bundle.
	npm run build

db-up: ## Start the local PostgreSQL service.
	docker compose up -d db

db-down: ## Stop local services without deleting data.
	docker compose down

db-migrate: ## Apply committed migrations to the configured database.
	npm run db:migrate

security: ## Query the npm advisory service; requires network access.
	npm run check:security
