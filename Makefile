# Makefile for Thrive Track Workout

.PHONY: help docker-up docker-down run migrate-up migrate-down test

help:
	@echo "Available commands:"
	@echo "  docker-up     - Start Docker containers"
	@echo "  docker-down   - Stop Docker containers"
	@echo "  run           - Run the Go application"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  test          - Runs the test suite"

docker-up:
	docker compose up --build

docker-down:
	docker compose down

run:
	go run main.go

migrate-up:
	goose -dir ./migration postgres "host=localhost user=root password=postgres dbname=postgres port=5432 sslmode=disable" up

migrate-down:
	goose -dir ./migration postgres "host=localhost user=root password=postgres dbname=postgres port=5432 sslmode=disable" down

test:
	go test -v ./...
