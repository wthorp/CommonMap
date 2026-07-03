.DEFAULT_GOAL := help

GOLANGCI_LINT_VERSION ?= $(shell cat .golangci-lint-version)
GOLANGCI_LINT := $(CURDIR)/bin/golangci-lint
LEFTHOOK_VERSION ?= $(shell cat .lefthook-version)
LEFTHOOK := $(CURDIR)/bin/lefthook
GOVULNCHECK_VERSION ?= $(shell cat .govulncheck-version)
GOVULNCHECK := $(CURDIR)/bin/govulncheck

.PHONY: help
help: ## Show available targets
	@awk 'BEGIN { \
		FS = ":.*##"; \
		printf "Usage:\n  make \033[36m<target>\033[0m\n" \
	} \
	/^[a-zA-Z0-9_\/.-]+:.*?##/ { \
		printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 \
	} ' $(MAKEFILE_LIST)

.PHONY: tools
tools: $(GOLANGCI_LINT) $(LEFTHOOK) $(GOVULNCHECK) ## Install local development tools into ./bin

.PHONY: setup
setup: tools hooks ## Bootstrap local tools and install Git hooks

$(GOLANGCI_LINT):
	GOBIN=$(CURDIR)/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$(LEFTHOOK):
	GOBIN=$(CURDIR)/bin go install github.com/evilmartians/lefthook/v2@$(LEFTHOOK_VERSION)

$(GOVULNCHECK):
	GOBIN=$(CURDIR)/bin go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

.PHONY: fmt
fmt: ## Format Go source with gofmt
	gofmt -w $$(find . -type f -name '*.go' -not -path './bin/*')

.PHONY: fmt-check
fmt-check: ## Check whether Go source is gofmt-clean
	@test -z "$$(gofmt -l $$(find . -type f -name '*.go' -not -path './bin/*'))" || \
		(echo "gofmt reported unformatted files"; gofmt -l $$(find . -type f -name '*.go' -not -path './bin/*'); exit 1)

.PHONY: tidy
tidy: ## Normalize go.mod/go.sum
	go mod tidy

.PHONY: vet
vet: ## Run go vet across all packages
	go vet ./...

.PHONY: test
test: ## Run Go tests across all packages
	go test ./...

.PHONY: lint
lint: tools ## Run golangci-lint with the repo config
	$(GOLANGCI_LINT) run ./...

.PHONY: lint-fix
lint-fix: tools ## Apply supported golangci-lint fixes
	$(GOLANGCI_LINT) run --fix ./...

.PHONY: vuln
vuln: tools ## Run govulncheck against the module
	$(GOVULNCHECK) ./...

.PHONY: hooks
hooks: tools ## Install repo-managed Git hooks via lefthook
	$(LEFTHOOK) install

.PHONY: hooks-run
hooks-run: tools ## Run installed lefthook hooks manually
	$(LEFTHOOK) run pre-commit

.PHONY: check
check: fmt-check vet test lint ## Run the full local quality gate
