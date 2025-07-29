.PHONY: localdev-start
localdev-start:
	@echo "Starting local development environment"
	docker compose -f docker-compose.yml up -d


# DB COMMANDS

.PHONY: migration
migration:
	echo "---> Creating a new migration"
	echo "---> Input args... $(filename)"
	@migrate create -ext sql -dir ./db/migrations/ $(filename)

.PHONY: db-migrate
db-migrate:
	@echo "Running db migrations"
	go mod tidy
	@go run db/migrations/main.go up

.PHONY: db-migrate-down
db-migrate-down:
	@echo "Running db migrations down"
	go mod tidy
	@go run db/migrations/main.go down

.PHONY: db-rollback
db-rollback:
	@echo "Rolling back db migrations"
	go mod tidy
	@go run db/migrations/main.go down
	
	