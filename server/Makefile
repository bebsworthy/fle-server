# Makefile for FLE Server
# Backend-core feature development tools

# Application name
APP_NAME := fle-server
BINARY_NAME := $(APP_NAME)
BINARY_PATH := bin/$(BINARY_NAME)

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet

# Main application path
MAIN_PATH := ./cmd/server

# Test parameters
TEST_TIMEOUT := 30s
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html

# Linting
GOLANGCI_LINT := golangci-lint
GOLANGCI_LINT_VERSION := v1.55.2

# Build flags
BUILD_FLAGS := -v
LDFLAGS := -w -s

# Color output
NOCOLOR=\033[0m
RED=\033[0;31m
GREEN=\033[0;32m
ORANGE=\033[0;33m
BLUE=\033[0;34m
PURPLE=\033[0;35m
CYAN=\033[0;36m
LIGHTGRAY=\033[0;37m
DARKGRAY=\033[1;30m
LIGHTRED=\033[1;31m
LIGHTGREEN=\033[1;32m
YELLOW=\033[1;33m
LIGHTBLUE=\033[1;34m
LIGHTPURPLE=\033[1;35m
LIGHTCYAN=\033[1;36m
WHITE=\033[1;37m

##@ General

.PHONY: help
help: ## Display this help message
	@echo "$(LIGHTBLUE)FLE Server - Development Commands$(NOCOLOR)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(LIGHTCYAN)<target>$(NOCOLOR)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(LIGHTCYAN)%-15s$(NOCOLOR) %s\n", $$1, $$2 } /^##@/ { printf "\n$(LIGHTBLUE)%s$(NOCOLOR)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: run
run: ## Run the server in development mode
	@echo "$(LIGHTGREEN)Running server...$(NOCOLOR)"
	$(GORUN) $(MAIN_PATH)/main.go

.PHONY: run-watch
run-watch: ## Run the server with hot reload (requires entr)
	@echo "$(LIGHTGREEN)Running server with hot reload...$(NOCOLOR)"
	@command -v entr >/dev/null 2>&1 || { echo "$(RED)Error: entr is not installed. Install with: brew install entr$(NOCOLOR)"; exit 1; }
	find . -name "*.go" | entr -r make run

.PHONY: deps
deps: ## Download and verify dependencies
	@echo "$(LIGHTGREEN)Downloading dependencies...$(NOCOLOR)"
	$(GOMOD) download
	$(GOMOD) verify

.PHONY: tidy
tidy: ## Tidy up go.mod and go.sum
	@echo "$(LIGHTGREEN)Tidying modules...$(NOCOLOR)"
	$(GOMOD) tidy

##@ Building

.PHONY: build
build: ## Build the server binary
	@echo "$(LIGHTGREEN)Building $(BINARY_NAME)...$(NOCOLOR)"
	@mkdir -p bin
	$(GOBUILD) $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)Built: $(BINARY_PATH)$(NOCOLOR)"

.PHONY: build-debug
build-debug: ## Build the server binary with debug symbols
	@echo "$(LIGHTGREEN)Building $(BINARY_NAME) with debug symbols...$(NOCOLOR)"
	@mkdir -p bin
	$(GOBUILD) $(BUILD_FLAGS) -gcflags="all=-N -l" -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)Built: $(BINARY_PATH) (debug)$(NOCOLOR)"

.PHONY: build-all
build-all: ## Build for all platforms
	@echo "$(LIGHTGREEN)Building for all platforms...$(NOCOLOR)"
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" -o bin/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "$(GREEN)Built binaries for all platforms in bin/$(NOCOLOR)"

##@ Testing

.PHONY: test
test: ## Run all tests
	@echo "$(LIGHTGREEN)Running tests...$(NOCOLOR)"
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./...

.PHONY: test-short
test-short: ## Run short tests only
	@echo "$(LIGHTGREEN)Running short tests...$(NOCOLOR)"
	$(GOTEST) -v -short -timeout $(TEST_TIMEOUT) ./...

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "$(LIGHTGREEN)Running tests with race detection...$(NOCOLOR)"
	$(GOTEST) -v -race -timeout $(TEST_TIMEOUT) ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage and generate HTML report
	@echo "$(LIGHTGREEN)Running tests with coverage...$(NOCOLOR)"
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) -coverprofile=$(COVERAGE_OUT) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "$(GREEN)Coverage report: $(COVERAGE_HTML)$(NOCOLOR)"
	@echo "$(CYAN)Opening coverage report in browser...$(NOCOLOR)"
	@command -v open >/dev/null 2>&1 && open $(COVERAGE_HTML) || echo "Open $(COVERAGE_HTML) in your browser"

.PHONY: test-coverage-func
test-coverage-func: ## Show test coverage by function
	@echo "$(LIGHTGREEN)Generating function coverage...$(NOCOLOR)"
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) -coverprofile=$(COVERAGE_OUT) ./...
	$(GOCMD) tool cover -func=$(COVERAGE_OUT)

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "$(LIGHTGREEN)Running benchmarks...$(NOCOLOR)"
	$(GOTEST) -v -bench=. -benchmem ./...

##@ Code Quality

.PHONY: fmt
fmt: ## Format all Go code
	@echo "$(LIGHTGREEN)Formatting code...$(NOCOLOR)"
	$(GOFMT) ./...

.PHONY: vet
vet: ## Run go vet
	@echo "$(LIGHTGREEN)Running go vet...$(NOCOLOR)"
	$(GOVET) ./...

.PHONY: lint
lint: ## Run golangci-lint
	@echo "$(LIGHTGREEN)Running golangci-lint...$(NOCOLOR)"
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "$(RED)Error: golangci-lint is not installed. Run 'make install-lint'$(NOCOLOR)"; exit 1; }
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with --fix
	@echo "$(LIGHTGREEN)Running golangci-lint with auto-fix...$(NOCOLOR)"
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "$(RED)Error: golangci-lint is not installed. Run 'make install-lint'$(NOCOLOR)"; exit 1; }
	$(GOLANGCI_LINT) run --fix

.PHONY: install-lint
install-lint: ## Install golangci-lint
	@echo "$(LIGHTGREEN)Installing golangci-lint $(GOLANGCI_LINT_VERSION)...$(NOCOLOR)"
	@command -v curl >/dev/null 2>&1 || { echo "$(RED)Error: curl is required$(NOCOLOR)"; exit 1; }
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

.PHONY: check
check: fmt vet lint test ## Run all code quality checks

##@ Security

.PHONY: security
security: ## Run security checks with gosec
	@echo "$(LIGHTGREEN)Running security checks...$(NOCOLOR)"
	@command -v gosec >/dev/null 2>&1 || { echo "$(YELLOW)Installing gosec...$(NOCOLOR)"; go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; }
	gosec ./...

.PHONY: audit
audit: ## Audit dependencies for vulnerabilities
	@echo "$(LIGHTGREEN)Auditing dependencies...$(NOCOLOR)"
	@command -v govulncheck >/dev/null 2>&1 || { echo "$(YELLOW)Installing govulncheck...$(NOCOLOR)"; go install golang.org/x/vuln/cmd/govulncheck@latest; }
	govulncheck ./...

##@ Docker

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(LIGHTGREEN)Building Docker image...$(NOCOLOR)"
	docker build -t $(APP_NAME):latest .

.PHONY: docker-run
docker-run: docker-build ## Run Docker container
	@echo "$(LIGHTGREEN)Running Docker container...$(NOCOLOR)"
	docker run --rm -p 8080:8080 --env-file .env $(APP_NAME):latest

##@ Database

.PHONY: db-up
db-up: ## Start database containers (when implemented)
	@echo "$(YELLOW)Database support will be added in future features$(NOCOLOR)"

.PHONY: db-down
db-down: ## Stop database containers (when implemented)
	@echo "$(YELLOW)Database support will be added in future features$(NOCOLOR)"

##@ Utilities

.PHONY: clean
clean: ## Clean build artifacts and temporary files
	@echo "$(LIGHTGREEN)Cleaning...$(NOCOLOR)"
	$(GOCLEAN)
	rm -rf bin/
	rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)
	rm -f *.prof
	@echo "$(GREEN)Cleaned build artifacts$(NOCOLOR)"

.PHONY: clean-cache
clean-cache: ## Clean Go module cache
	@echo "$(LIGHTGREEN)Cleaning module cache...$(NOCOLOR)"
	$(GOCMD) clean -modcache

.PHONY: version
version: ## Show Go version and module info
	@echo "$(LIGHTGREEN)Go version:$(NOCOLOR)"
	@$(GOCMD) version
	@echo "$(LIGHTGREEN)Module info:$(NOCOLOR)"
	@$(GOMOD) list -m

.PHONY: deps-update
deps-update: ## Update all dependencies to latest versions
	@echo "$(LIGHTGREEN)Updating dependencies...$(NOCOLOR)"
	$(GOGET) -u ./...
	$(GOMOD) tidy

.PHONY: deps-graph
deps-graph: ## Generate dependency graph (requires graphviz)
	@echo "$(LIGHTGREEN)Generating dependency graph...$(NOCOLOR)"
	@command -v dot >/dev/null 2>&1 || { echo "$(RED)Error: graphviz is not installed. Install with: brew install graphviz$(NOCOLOR)"; exit 1; }
	$(GOMOD) graph | dot -T png -o deps.png
	@echo "$(GREEN)Dependency graph: deps.png$(NOCOLOR)"

##@ Development Tools

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(LIGHTGREEN)Installing development tools...$(NOCOLOR)"
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/kisielk/errcheck@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	@make install-lint
	@echo "$(GREEN)Development tools installed$(NOCOLOR)"

.PHONY: workspace-init
workspace-init: deps install-tools ## Initialize development workspace
	@echo "$(LIGHTGREEN)Initializing development workspace...$(NOCOLOR)"
	@cp -n .env.example .env 2>/dev/null || echo "$(YELLOW).env file already exists$(NOCOLOR)"
	@make check
	@echo "$(GREEN)Workspace initialized! Run 'make run' to start the server$(NOCOLOR)"

##@ CI/CD

.PHONY: ci-test
ci-test: ## Run tests for CI environment
	@echo "$(LIGHTGREEN)Running CI tests...$(NOCOLOR)"
	$(GOTEST) -v -race -timeout $(TEST_TIMEOUT) -coverprofile=$(COVERAGE_OUT) ./...

.PHONY: ci-lint
ci-lint: ## Run linting for CI environment
	@echo "$(LIGHTGREEN)Running CI linting...$(NOCOLOR)"
	$(GOLANGCI_LINT) run --out-format=github-actions

.PHONY: ci-build
ci-build: ## Build for CI environment
	@echo "$(LIGHTGREEN)Running CI build...$(NOCOLOR)"
	$(GOBUILD) $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" -o $(BINARY_PATH) $(MAIN_PATH)

.PHONY: ci
ci: ci-lint ci-test ci-build ## Run full CI pipeline

# Default target
.DEFAULT_GOAL := help