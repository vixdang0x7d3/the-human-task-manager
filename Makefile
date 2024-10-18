
include .env

SHELL := /bin/bash

.PHONY: build
build:
	@go build -o bin/thtm

.PHONY: run
run: build
	@./bin/thtm

.PHONY: goose-up
goose-up:
	@GOOSE_MIGRATION_DIR=$(DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

.PHONY: goose-down
goose-down:
	@GOOSE_MIGRATION_DIR=$(DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down

.PHONY: goose-upto
goose-upto:
	@GOOSE_MIGRATION_DIR=$(DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up-to $(VERSION)

.PHONY: goose-downto
goose-downto:
	@GOOSE_MIGRATION_DIR=$(DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down-to $(VERSION)


.PHONY: goose-status
goose-status:
	@GOOSE_MIGRATION_DIR=$(DIR) goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) status

.PHONY: live
live:
	@air

.PHONY: clean
clean: unset
	rm -rf bin/* tmp/*
	
