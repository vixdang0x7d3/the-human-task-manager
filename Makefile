
include .env

SHELL := /bin/bash

.PHONY: build
build:
	@go build -o bin/thtm

.PHONY: run
run: build
	@./bin/thtm

.PHONY: migrate-up
migrate-up:
	@GOOSE_MIGRATION_DIR=$(DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

.PHONY: migrate-down
migrate-down:
	@GOOSE_MIGRATION_DIR=$(DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down

.PHONY: migrate-upto
migrate-upto:
	@GOOSE_MIGRATION_DIR=$(DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up-to $(VERSION)

.PHONY: migrate-downto
migrate-downto:
	@GOOSE_MIGRATION_DIR=$(DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down-to $(VERSION)

.PHONY: live
live:
	@air

.PHONY: clean
clean: unset
	rm -rf bin/* tmp/*
	
