.PHONY: help build build-local up down logs ps test
.DEDAULT_GOAL := help

DOCKER_TAB := latest
build: ## Build docker image to deploy
	docker build -t budougumi0617/gotodo:${DOCKER_TAB} \
		-- target deploy ./

build-local: ## Build docker image to local development
	docker compose build --no-cache

up: ## Do docker compose up with hot reload
	docker compose up -d

down: ## Do docker compose down
	docker compose down

logs: ## Tail docker compose logs
	docker compose logs -f

ps: # Check container status
	docker compose ps

test: ## Execute tests
	go test -race -shuffle=on ./...

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$$' ${MAKEFILE_LIST} | \
		awk `BEGIN {FS = ":.*?## "}; {prontf "\033[36m%-20s\033[0m %s\n", $$1, $$2}`