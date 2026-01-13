.PHONY: run test build clean deps help swagger swagger-serve dev docker-build docker-up docker-down docker-logs docker-restart docker-clean docker-rebuild docker-ps docker-logs-api docker-logs-localstack docker-up-build docker-init-resources docker-up-init docker-check-resources

# Variables
BINARY_NAME=checker-api
MAIN_PATH=cmd/server/main.go
BUILD_DIR=bin

# Detect Docker Compose command (try docker compose first, then docker-compose)
DOCKER_CMD := $(shell command -v docker 2> /dev/null)
DOCKER_COMPOSE_CMD := $(shell if docker compose version >/dev/null 2>&1; then echo "docker compose"; elif command -v docker-compose >/dev/null 2>&1; then echo "docker-compose"; else echo ""; fi)

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
RED=\033[0;31m
NC=\033[0m # No Color

# Check if Docker is available
check-docker:
	@if [ -z "$(DOCKER_CMD)" ]; then \
		echo "$(RED)Error: Docker is not installed or not in PATH$(NC)"; \
		echo "$(YELLOW)Please install Docker Desktop for macOS:$(NC)"; \
		echo "  Visit: https://www.docker.com/products/docker-desktop/"; \
		echo "  Or install via Homebrew: brew install --cask docker"; \
		exit 1; \
	fi
	@if [ -z "$(DOCKER_COMPOSE_CMD)" ]; then \
		echo "$(RED)Error: Docker Compose is not available$(NC)"; \
		echo "$(YELLOW)Please install docker-compose using one of these methods:$(NC)"; \
		echo "  1. Install via Homebrew: $(YELLOW)brew install docker-compose$(NC)"; \
		echo "  2. Or install Docker Desktop which includes Docker Compose"; \
		echo "  3. Or install via pip: $(YELLOW)pip install docker-compose$(NC)"; \
		exit 1; \
	fi
	@if ! docker info >/dev/null 2>&1; then \
		echo "$(RED)Error: Cannot connect to Docker daemon$(NC)"; \
		echo "$(YELLOW)Docker is installed but the daemon is not running$(NC)"; \
		echo "$(YELLOW)Please start Docker Desktop:$(NC)"; \
		echo "  1. Open Docker Desktop from Applications"; \
		echo "  2. Wait for Docker to start (whale icon in menu bar)"; \
		echo "  3. Then try again"; \
		exit 1; \
	fi

## help: Show this help message
help:
	@echo "$(GREEN)Available commands:$(NC)"
	@echo ""
	@echo "$(BLUE)Application:$(NC)"
	@echo "  $(YELLOW)make run$(NC)        - Run the application"
	@echo "  $(YELLOW)make dev$(NC)        - Generate Swagger, run app and open Swagger UI"
	@echo "  $(YELLOW)make test$(NC)       - Run tests"
	@echo "  $(YELLOW)make build$(NC)      - Build the application"
	@echo "  $(YELLOW)make clean$(NC)     - Clean generated files"
	@echo "  $(YELLOW)make deps$(NC)      - Install/update dependencies"
	@echo "  $(YELLOW)make swagger$(NC)   - Generate Swagger documentation"
	@echo ""
	@echo "$(BLUE)Docker:$(NC)"
	@echo "  $(YELLOW)make docker-build$(NC)          - Build Docker images"
	@echo "  $(YELLOW)make docker-up$(NC)            - Start Docker containers"
	@echo "  $(YELLOW)make docker-up-build$(NC)      - Build and start Docker containers"
	@echo "  $(YELLOW)make docker-down$(NC)          - Stop Docker containers"
	@echo "  $(YELLOW)make docker-logs$(NC)          - Show Docker logs (all services)"
	@echo "  $(YELLOW)make docker-logs-api$(NC)      - Show logs for checker-api service"
	@echo "  $(YELLOW)make docker-logs-localstack$(NC) - Show logs for LocalStack service"
	@echo "  $(YELLOW)make docker-restart$(NC)       - Restart Docker containers"
	@echo "  $(YELLOW)make docker-clean$(NC)         - Stop and remove containers, volumes, and images"
	@echo "  $(YELLOW)make docker-rebuild$(NC)       - Rebuild and restart containers"
	@echo "  $(YELLOW)make docker-ps$(NC)            - Show running Docker containers"
	@echo "  $(YELLOW)make docker-init-resources$(NC) - Create DynamoDB table and SQS queue in LocalStack"
	@echo "  $(YELLOW)make docker-up-init$(NC)       - Start containers and initialize resources"
	@echo "  $(YELLOW)make docker-check-resources$(NC) - Check if DynamoDB table and SQS queue exist"

## run: Run the application
run:
	@echo "$(GREEN)Running application...$(NC)"
	@go run $(MAIN_PATH)

## test: Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	@go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage generated in coverage.html$(NC)"

## build: Build the application
build:
	@echo "$(GREEN)Building application...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Binary generated in $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## clean: Clean generated files
clean:
	@echo "$(GREEN)Cleaning generated files...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -rf docs
	@echo "$(GREEN)Cleanup completed$(NC)"

## deps: Install/update dependencies
deps:
	@echo "$(GREEN)Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

## fmt: Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Formatting completed$(NC)"

## vet: Run go vet to detect common errors
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)Analysis completed$(NC)"

## lint: Ejecuta fmt y vet
lint: fmt vet

## swagger: Generate Swagger documentation
swagger:
	@echo "$(GREEN)Checking for swag CLI...$(NC)"
	@if ! command -v swag > /dev/null 2>&1; then \
		echo "$(YELLOW)swag not found, installing...$(NC)"; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@echo "$(GREEN)Generating Swagger documentation...$(NC)"
	@if command -v swag > /dev/null 2>&1; then \
		swag init -g cmd/server/main.go -o docs; \
	else \
		$(shell go env GOPATH)/bin/swag init -g cmd/server/main.go -o docs; \
	fi
	@echo "$(GREEN)Swagger documentation generated in docs/$(NC)"

## swagger-serve: Generate and show Swagger documentation URL
swagger-serve: swagger
	@echo "$(GREEN)Swagger documentation available at http://localhost:8080/swagger/index.html$(NC)"

## dev: Generate Swagger, run application and open Swagger UI in browser
dev: swagger
	@echo "$(GREEN)Starting application with Swagger...$(NC)"
	@echo "$(YELLOW)Swagger UI will open automatically in your browser$(NC)"
	@(sleep 3 && open http://localhost:$${SERVER_PORT:-8080}/swagger/index.html 2>/dev/null || echo "$(YELLOW)Could not open browser automatically. Please visit: http://localhost:$${SERVER_PORT:-8080}/swagger/index.html$(NC)") &
	@go run $(MAIN_PATH)

## docker-build: Build Docker images
docker-build: check-docker
	@echo "$(GREEN)Building Docker images...$(NC)"
	@$(DOCKER_COMPOSE_CMD) build

## docker-up: Start Docker containers
docker-up: check-docker build-init-script
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	@$(DOCKER_COMPOSE_CMD) up -d
	@echo "$(GREEN)Containers started. API available at http://localhost:8080$(NC)"
	@echo "$(GREEN)LocalStack available at http://localhost:4566$(NC)"

## docker-up-build: Build and start Docker containers
docker-up-build: check-docker build-init-script docker-build docker-up

## docker-down: Stop Docker containers
docker-down: check-docker
	@echo "$(GREEN)Stopping Docker containers...$(NC)"
	@$(DOCKER_COMPOSE_CMD) down
	@echo "$(GREEN)Containers stopped$(NC)"

## docker-logs: Show Docker logs
docker-logs: check-docker
	@echo "$(GREEN)Showing Docker logs...$(NC)"
	@$(DOCKER_COMPOSE_CMD) logs -f

## docker-logs-api: Show logs for checker-api service only
docker-logs-api: check-docker
	@echo "$(GREEN)Showing checker-api logs...$(NC)"
	@$(DOCKER_COMPOSE_CMD) logs -f checker-api

## docker-logs-localstack: Show logs for LocalStack service only
docker-logs-localstack: check-docker
	@echo "$(GREEN)Showing LocalStack logs...$(NC)"
	@$(DOCKER_COMPOSE_CMD) logs -f localstack

## docker-restart: Restart Docker containers
docker-restart: check-docker
	@echo "$(GREEN)Restarting Docker containers...$(NC)"
	@$(DOCKER_COMPOSE_CMD) restart
	@echo "$(GREEN)Containers restarted$(NC)"

## docker-clean: Stop and remove containers, volumes, and images
docker-clean: check-docker
	@echo "$(YELLOW)This will remove containers, volumes, and images. Are you sure?$(NC)"
	@echo "$(GREEN)Stopping and removing Docker containers, volumes, and images...$(NC)"
	@$(DOCKER_COMPOSE_CMD) down -v --rmi local
	@echo "$(GREEN)Cleanup completed$(NC)"

## docker-rebuild: Rebuild and restart containers
docker-rebuild: check-docker
	@echo "$(GREEN)Rebuilding and restarting Docker containers...$(NC)"
	@$(DOCKER_COMPOSE_CMD) up -d --build
	@echo "$(GREEN)Containers rebuilt and restarted$(NC)"

## docker-ps: Show running Docker containers
docker-ps: check-docker
	@echo "$(GREEN)Running Docker containers:$(NC)"
	@$(DOCKER_COMPOSE_CMD) ps

## build-init-script: Build the LocalStack initialization script
build-init-script:
	@echo "$(GREEN)Compilando script de inicialización...$(NC)"
	@cd scripts && go build -o localstack-init localstack-init.go
	@echo "$(GREEN)✅ Script compilado: scripts/localstack-init$(NC)"

## docker-init-resources: Create DynamoDB table and SQS queue in LocalStack
docker-init-resources: check-docker
	@echo "$(GREEN)Creando recursos en LocalStack...$(NC)"
	@echo "$(YELLOW)Esperando a que LocalStack esté listo...$(NC)"
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if curl -s http://localhost:4566/_localstack/health >/dev/null 2>&1; then \
			echo "$(GREEN)LocalStack está listo$(NC)"; \
			break; \
		fi; \
		if [ $$i -eq 10 ]; then \
			echo "$(RED)Error: LocalStack no está disponible$(NC)"; \
			echo "$(YELLOW)Por favor ejecuta: make docker-up$(NC)"; \
			exit 1; \
		fi; \
		sleep 2; \
	done
	@echo "$(GREEN)Creando tabla DynamoDB: comparisons$(NC)"
	@curl -s -X POST http://localhost:4566/ \
		-H "Content-Type: application/x-amz-json-1.0" \
		-H "X-Amz-Target: DynamoDB_20120810.CreateTable" \
		-d '{"TableName":"comparisons","AttributeDefinitions":[{"AttributeName":"id","AttributeType":"S"}],"KeySchema":[{"AttributeName":"id","KeyType":"HASH"}],"BillingMode":"PAY_PER_REQUEST"}' \
		>/dev/null 2>&1 && echo "$(GREEN)✅ Tabla DynamoDB creada$(NC)" || echo "$(YELLOW)ℹ️  La tabla ya existe o hubo un error$(NC)"
	@echo "$(GREEN)Creando cola SQS: comparison-queue$(NC)"
	@GET_RESPONSE=$$(curl -s -X POST http://localhost:4566/ \
		-H "Content-Type: application/x-amz-json-1.0" \
		-H "X-Amz-Target: AWSSimpleQueueServiceV20121105.GetQueueUrl" \
		-d '{"QueueName":"comparison-queue"}' 2>&1); \
	if echo "$$GET_RESPONSE" | grep -q "QueueUrl"; then \
		QUEUE_URL=$$(echo "$$GET_RESPONSE" | grep -o '"QueueUrl":"[^"]*"' | cut -d'"' -f4); \
		echo "$(GREEN)✅ Cola SQS ya existe: $$QUEUE_URL$(NC)"; \
	else \
		echo "$(YELLOW)La cola no existe, creándola...$(NC)"; \
		CREATE_RESPONSE=$$(curl -s -w "\n%{http_code}" -X POST http://localhost:4566/ \
			-H "Content-Type: application/x-amz-json-1.0" \
			-H "X-Amz-Target: AWSSimpleQueueServiceV20121105.CreateQueue" \
			-d '{"QueueName":"comparison-queue"}' 2>&1); \
		HTTP_CODE=$$(echo "$$CREATE_RESPONSE" | tail -n1); \
		BODY=$$(echo "$$CREATE_RESPONSE" | sed '$$d'); \
		if [ "$$HTTP_CODE" = "200" ] || echo "$$BODY" | grep -q "QueueUrl"; then \
			echo "$(GREEN)✅ Cola SQS creada exitosamente$(NC)"; \
			if echo "$$BODY" | grep -q "QueueUrl"; then \
				QUEUE_URL=$$(echo "$$BODY" | grep -o '"QueueUrl":"[^"]*"' | cut -d'"' -f4); \
				echo "$(GREEN)   URL: $$QUEUE_URL$(NC)"; \
			fi; \
		else \
			if echo "$$BODY" | grep -q "QueueAlreadyExists" || echo "$$BODY" | grep -q "QueueAlreadyExistsException"; then \
				echo "$(YELLOW)ℹ️  La cola ya existe (creada por otro proceso)$(NC)"; \
			else \
				echo "$(YELLOW)⚠️  HTTP Code: $$HTTP_CODE$(NC)"; \
				echo "$(YELLOW)⚠️  Respuesta: $$BODY$(NC)"; \
			fi; \
		fi; \
		echo "$(YELLOW)Esperando a que la cola esté disponible...$(NC)"; \
		QUEUE_VERIFIED=0; \
		for i in 1 2 3 4 5; do \
			sleep 3; \
			GET_RESPONSE=$$(curl -s -X POST http://localhost:4566/ \
				-H "Content-Type: application/x-amz-json-1.0" \
				-H "X-Amz-Target: AWSSimpleQueueServiceV20121105.GetQueueUrl" \
				-d '{"QueueName":"comparison-queue"}' 2>&1); \
			if [ -n "$$GET_RESPONSE" ] && echo "$$GET_RESPONSE" | grep -q "QueueUrl"; then \
				QUEUE_URL=$$(echo "$$GET_RESPONSE" | grep -o '"QueueUrl":"[^"]*"' | cut -d'"' -f4); \
				echo "$(GREEN)✅ Cola SQS verificada (intento $$i): $$QUEUE_URL$(NC)"; \
				QUEUE_VERIFIED=1; \
				break; \
			fi; \
		done; \
		if [ $$QUEUE_VERIFIED -eq 0 ]; then \
			echo "$(YELLOW)⚠️  No se pudo verificar inmediatamente, pero la cola puede existir$(NC)"; \
			echo "$(YELLOW)   Verifica con: make docker-check-resources$(NC)"; \
		fi; \
	fi
	@echo "$(GREEN)✅ Recursos inicializados$(NC)"

## docker-up-init: Start containers and initialize resources
docker-up-init: build-init-script docker-up docker-init-resources
	@echo "$(GREEN)✅ Aplicación lista para usar$(NC)"
	@echo "$(YELLOW)API disponible en: http://localhost:8080$(NC)"
	@echo "$(YELLOW)Swagger UI: http://localhost:8080/swagger/index.html$(NC)"

## docker-check-resources: Check if DynamoDB table and SQS queue exist
docker-check-resources: check-docker
	@echo "$(GREEN)Verificando recursos en LocalStack...$(NC)"
	@echo ""
	@echo "$(YELLOW)Esperando a que LocalStack esté listo...$(NC)"
	@for i in 1 2 3 4 5; do \
		if curl -s http://localhost:4566/_localstack/health >/dev/null 2>&1; then \
			echo "$(GREEN)LocalStack está listo$(NC)"; \
			break; \
		fi; \
		if [ $$i -eq 5 ]; then \
			echo "$(RED)Error: LocalStack no está disponible$(NC)"; \
			echo "$(YELLOW)Por favor ejecuta: make docker-up$(NC)"; \
			exit 1; \
		fi; \
		sleep 2; \
	done
	@echo ""
	@echo "$(BLUE)════════════════════════════════════════$(NC)"
	@echo "$(YELLOW)Tablas DynamoDB:$(NC)"
	@echo "$(BLUE)════════════════════════════════════════$(NC)"
	@TABLES=$$(curl -s -X POST http://localhost:4566/ \
		-H "Content-Type: application/x-amz-json-1.0" \
		-H "X-Amz-Target: DynamoDB_20120810.ListTables" \
		-d '{}' 2>/dev/null); \
	if echo "$$TABLES" | grep -q "comparisons"; then \
		echo "$(GREEN)✅ Tabla 'comparisons' existe$(NC)"; \
		echo "$$TABLES" | python3 -m json.tool 2>/dev/null || echo "$$TABLES"; \
	else \
		echo "$(RED)❌ Tabla 'comparisons' NO existe$(NC)"; \
		if [ -n "$$TABLES" ]; then \
			echo "Tablas disponibles:"; \
			echo "$$TABLES" | python3 -m json.tool 2>/dev/null || echo "$$TABLES"; \
		fi; \
	fi
	@echo ""
	@echo "$(BLUE)════════════════════════════════════════$(NC)"
	@echo "$(YELLOW)Colas SQS:$(NC)"
	@echo "$(BLUE)════════════════════════════════════════$(NC)"
	@GET_QUEUE_RESPONSE=$$(curl -s -X POST http://localhost:4566/ \
		-H "Content-Type: application/x-amz-json-1.0" \
		-H "X-Amz-Target: AWSSimpleQueueServiceV20121105.GetQueueUrl" \
		-d '{"QueueName":"comparison-queue"}' 2>/dev/null); \
	if [ -n "$$GET_QUEUE_RESPONSE" ] && echo "$$GET_QUEUE_RESPONSE" | grep -q "QueueUrl"; then \
		QUEUE_URL=$$(echo "$$GET_QUEUE_RESPONSE" | grep -o '"QueueUrl":"[^"]*"' | cut -d'"' -f4); \
		echo "$(GREEN)✅ Cola 'comparison-queue' existe$(NC)"; \
		echo "$(GREEN)   URL: $$QUEUE_URL$(NC)"; \
		QUEUES=$$(curl -s -X POST http://localhost:4566/ \
			-H "Content-Type: application/x-amz-json-1.0" \
			-H "X-Amz-Target: AWSSimpleQueueServiceV20121105.ListQueues" \
			-d '{}' 2>/dev/null); \
		if [ -n "$$QUEUES" ]; then \
			echo "Todas las colas:"; \
			echo "$$QUEUES" | python3 -m json.tool 2>/dev/null || echo "$$QUEUES"; \
		fi; \
	else \
		echo "$(RED)❌ Cola 'comparison-queue' NO existe$(NC)"; \
		QUEUES=$$(curl -s -X POST http://localhost:4566/ \
			-H "Content-Type: application/x-amz-json-1.0" \
			-H "X-Amz-Target: AWSSimpleQueueServiceV20121105.ListQueues" \
			-d '{}' 2>/dev/null); \
		if [ -n "$$QUEUES" ]; then \
			echo "Colas disponibles:"; \
			echo "$$QUEUES" | python3 -m json.tool 2>/dev/null || echo "$$QUEUES"; \
		fi; \
		echo "$(YELLOW)Ejecuta: make docker-init-resources$(NC)"; \
	fi
	@echo ""
	@echo "$(GREEN)✅ Verificación completada$(NC)"
