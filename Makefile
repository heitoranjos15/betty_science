# Variables
DOCKER_COMPOSE = docker-compose
GO = go
APP_NAME = betty-science

# Default target
.PHONY: all
all: build

# Build the Docker images
.PHONY: build
build:
	$(DOCKER_COMPOSE) build

# Start the services
.PHONY: up
up:
	$(DOCKER_COMPOSE) up -d

# Stop the services
.PHONY: down
down:
	$(DOCKER_COMPOSE) down

# Run the Go application locally (without Docker)
.PHONY: run
run:
	$(GO) run main.go

# Run the bot located in cmd/bot/riot.go
.PHONY: load_frames
load_frames:
	$(GO) run cmd/bot/riot.go load_frames

load_games:
	$(GO) run cmd/bot/riot.go load_games

.PHONY: schedule
schedule:
	$(GO) run cmd/league_of_legends/command.go CMD=schedule

game:
	$(GO) run cmd/league_of_legends/command.go CMD=game

frame:
	$(GO) run cmd/league_of_legends/command.go CMD=frame

.PHONY: update_matches
update_matches:
	$(GO) run cmd/bot/riot.go CMD=update_matches

.PHONY: complete_games
complete_games:
	$(GO) run cmd/bot/riot.go CMD=complete_games

# Build the Go application
.PHONY: build-go
build-go:
	$(GO) build -o $(APP_NAME) .

# Clean up Docker resources
.PHONY: clean
clean:
	$(DOCKER_COMPOSE) down -v --rmi all --remove-orphans
	rm -f $(APP_NAME)

# View logs
.PHONY: logs
logs:
	$(DOCKER_COMPOSE) logs -f

# Mongo-shell
.PHONY: mongo-shell
mongo-shell:
	$(DOCKER_COMPOSE) exec mongo mongosh

# Run tests
.PHONY: test
test:
	$(GO) test ./... -v

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# Check for linting issues
.PHONY: lint
lint:
	golangci-lint run
