include .envrc

.PHONY: run
run:
	go run ./cmd/api --port

.PHONY: psql
psql:
	psql ${GLOBECHAT_DB_DSN}

.PHONY: up
up:
	@echo "Running up migrations"
	migrate -path ./migrations -database ${GLOBECHAT_DB_DSN} up