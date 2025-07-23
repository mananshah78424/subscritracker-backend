.PHONY: db-migrate
db-migrate:
	@echo "Running db migrations"
	go mod tidy
	@go run cmd/migrations/main.go up