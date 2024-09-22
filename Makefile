build:
	@go build -o bin/thtm

run: build
	@./bin/thtm
