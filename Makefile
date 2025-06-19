.PHONY: build up down logs restart build-app lint swagger

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

restart: down up

build-app:
	CGO_ENABLED=0 go build -o app ./cmd

lint:
	golangci-lint run

swagger:
	swag init -g ./cmd/main.go