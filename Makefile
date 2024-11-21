include .env

.PHONY: db-up \
	db-down \
	goose-up \
	goose-down \
	goose-upto \
	goose-downto \
	goose status \
	sqlc \
	templ-build \
	tailwind-build \
	live/tailwind \
	build \
	run \
	test \
	live \
	clean

db-up:
	@docker compose up -d

db-down:
	@docker compose down

goose-up:
	@GOOSE_MIGRATION_DIR=$(MIGRATION_DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

goose-down:
	@GOOSE_MIGRATION_DIR=$(MIGRATION_DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down

goose-upto:
	@GOOSE_MIGRATION_DIR=$(MIGRATION_DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up-to $(VERSION)

goose-downto:
	@GOOSE_MIGRATION_DIR=$(MIGRATION_DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down-to $(VERSION)

goose-status:
	@GOOSE_MIGRATION_DIR=$(MIGRATION_DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) status

sqlc:
	@sqlc generate

tailwind-build:
	@tailwindcss -i static/css/input.css -o static/css/output.css --minify

templ-build:
	@templ generate -path internal/template/

build:
	@go build -o tmp/thtm

run: build
	@./tmp/thtm

test:
	@go test ./...

live/tailwind:
	@tailwindcss -i static/css/input.css -o static/css/output.css --watch --minify

live: 
	@air

clean: db-down goose-down
	rm -rf bin/* tmp/*
	
