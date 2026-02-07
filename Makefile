APP_NAME=prodify-api
GO_CMD=go
ATLAS_CMD=atlas
SWAG_CMD=swag
DOCKER_COMPOSE_CMD=docker-compose

server:
	air \
	--build.cmd "go build -o tmp/main ./cmd/hercules" \
	--build.bin "tmp/main" \
	--build.delay "100" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

run:
	$(GO_CMD) run cmd/api.go

build:
	$(GO_CMD) build -o bin/$(APP_NAME) cmd/api.go

test:
	$(GO_CMD) test ./... -coverprofile=coverage.out

fmt:
	$(GO_CMD) fmt ./...

docs:
	$(SWAG_CMD) init -g cmd/api.go -o internal/docs

migrate.diff:
	$(ATLAS_CMD) migrate diff --env development

migrate:
	$(ATLAS_CMD) migrate apply --env development

migrate.status:
	$(ATLAS_CMD) migrate status --env development

migrate.lint:
	$(ATLAS) migrate lint --env development

docker.up:
	$(DOCKER_COMPOSE_CMD) up -d

docker.down:
	$(DOCKER_COMPOSE_CMD) down

docker.build:
	docker build -t prodify-api .

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  server - Run the server"
	@echo "  run - Run the application"
	@echo "  build - Build the application"
	@echo "  test - Run the tests"
	@echo "  fmt - Format the code"
	@echo "  docs - Generate the documentation"
	@echo "  migrate.diff - Generate the migration diff"
	@echo "  migrate - Apply the migrations"
	@echo "  migrate.status - Show the migration status"
	@echo "  migrate.lint - Lint the migrations"
	@echo "  docker.up - Start the docker containers"
	@echo "  docker.down - Stop the docker containers"
	@echo "  docker.build - Build the docker image"
	@echo "  help - Show this help message"