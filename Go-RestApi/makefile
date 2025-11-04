include .env
export

MIGRATION_DIR=./pkg/db/migrations

migrate-up:
	goose -dir ${MIGRATION_DIR} postgres "${DB_URL}" up

migrate-down:
	goose -dir ${MIGRATION_DIR} postgres "${DB_URL}" down

migrate-status:
	goose -dir $(MIGRATION_DIR) postgres "${DB_URL}" status

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=add_users_table"; \
	else \
		goose -dir $(MIGRATION_DIR) create $(name) sql; \
	fi
