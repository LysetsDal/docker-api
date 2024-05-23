app := bin/docker-api

build:
	@go build -o $(app) cmd/*.go

run: build
	@./$(app)

