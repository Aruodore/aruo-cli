SHELL := /bin/sh

GO ?= go
GOLANGCI_LINT ?= golangci-lint
GORELEASER ?= goreleaser
SYFT ?= syft
COVERAGE_DIR ?= coverage

.DEFAULT_GOAL := help

.PHONY: help bootstrap tools fmt fmt-check lint vet test test-race coverage benchmark tidy tidy-check check release-check clean

help: ## Show the supported development tasks.
	@awk 'BEGIN {FS = ":.*## "; printf "Aruo development tasks:\n\n"} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

bootstrap: ## Download Go modules and verify required external tools.
	$(GO) mod download
	@$(MAKE) tools

tools: ## Verify pinned tools are installed; versions are defined in .tool-versions.
	@command -v $(GO) >/dev/null || { echo "error: Go is required (see .go-version)" >&2; exit 1; }
	@command -v $(GOLANGCI_LINT) >/dev/null || { echo "error: golangci-lint 2.11.4 is required (see docs/development)" >&2; exit 1; }
	@command -v $(GORELEASER) >/dev/null || { echo "warning: GoReleaser 2.17.1 is required only for release work" >&2; }
	@command -v $(SYFT) >/dev/null || { echo "warning: Syft 1.44.0 is required only for release work" >&2; }

fmt: ## Format Go source with the configured golangci-lint formatters.
	$(GOLANGCI_LINT) fmt

fmt-check: ## Fail when Go source is not canonically formatted.
	@test -z "$$($(GOLANGCI_LINT) fmt --diff)" || { $(GOLANGCI_LINT) fmt --diff; exit 1; }

lint: ## Run the curated Go lint policy.
	$(GOLANGCI_LINT) run

vet: ## Run Go's standard static analysis.
	$(GO) vet ./...

test: ## Run deterministic unit and integration tests.
	$(GO) test -shuffle=on ./...

test-race: ## Run tests with the race detector.
	$(GO) test -race -shuffle=on ./...

coverage: ## Write an atomic-mode coverage profile.
	@mkdir -p $(COVERAGE_DIR)
	$(GO) test -race -covermode=atomic -coverprofile=$(COVERAGE_DIR)/coverage.out ./...

benchmark: ## Run Go benchmarks without using cached test results.
	$(GO) test -run='^$$' -bench=. -benchmem -count=5 ./...

tidy: ## Normalize Go module metadata.
	$(GO) mod tidy

tidy-check: ## Fail if go.mod or go.sum is not tidy.
	$(GO) mod tidy -diff

check: fmt-check tidy-check vet lint test-race ## Run the full local pull-request gate.

release-check: ## Validate release configuration without publishing.
	@command -v $(SYFT) >/dev/null || { echo "error: Syft 1.44.0 is required for SBOM generation" >&2; exit 1; }
	$(GORELEASER) check
	$(GORELEASER) release --snapshot --clean

clean: ## Remove local build, coverage, and release output.
	$(GO) clean -testcache
	@rm -rf -- $(COVERAGE_DIR) dist
