SSO_DB_DSN=postgres://myuser:mypassword@localhost:5555/postgres?sslmode=disable
MIGRATIONS_DIR=sql/

.PHONY: migrate-up migrate-down migrate-status migrate-create

# Создать новую миграцию: make migrate-create name=name_table
migrate-create:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

