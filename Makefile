-include .envrc

.PHONY: build
build: build_frontend
	go build -o bin/web ./cmd/web

.PHONY: run_backend
run_backend:
	go run ./cmd/web -dsn ${GLOBECHAT_DB_DSN} -gclientid ${PUBLIC_GOOGLE_CLIENT_ID} -mediadir ./media

.PHONY: run_frontend
run_frontend:
	cd ui && PUBLIC_GOOGLE_CLIENT_ID=${PUBLIC_GOOGLE_CLIENT_ID} npm run dev

.PHONY: psql
psql:
	psql ${GLOBECHAT_DB_DSN}

.PHONY: migrate_up
migrate_up:
	@echo "Running up migrations"
	migrate -path ./migrations -database ${GLOBECHAT_DB_DSN} up

.PHONY: migrate_down
migrate_down:
	@echo "Running down migration"
	migrate -path ./migrations -database ${GLOBECHAT_DB_DSN} down 1

.PHONY: create_migration
create_migration:
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name is required. Usage: make migration name=your_migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: force_migration
force_migration:
	@if [ -z "$(version)" ]; then \
		echo "Error: Migration version is required. Usage: make migration version=your_migration_version"; \
		exit 1; \
	fi
	@echo "Forcing migration version: $(version)"
	migrate -path ./migrations -database ${GLOBECHAT_DB_DSN} force ${version}

.PHONY: migration_version
migration_version:
	migrate -path ./migrations -database ${GLOBECHAT_DB_DSN} version

.PHONY: build_frontend
build_frontend:
	cd ui && $(if $(PUBLIC_GOOGLE_CLIENT_ID),PUBLIC_GOOGLE_CLIENT_ID=$(PUBLIC_GOOGLE_CLIENT_ID)) npm run build
	mkdir -p ./frontend/static
	rm -rf ./frontend/static/*
	cp -r ./ui/build/* ./frontend/static/