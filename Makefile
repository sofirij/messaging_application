env ?= .env
include $(env)

BINARY_NAME=app.exe
BUILD_DIR=bin
MAIN_PATH=./cmd/server/main.go
MIGRATE=migrate -path migrations -database "$(DB_URL)"

# default target — runs when you just type "make"
all: format tidy build run

# format code
format:
	gofmt -w .

# format and fix imports
imports:
	go mod tidy

# build the binary
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# run the binary
run:
	-$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test ./internal/test -v -failfast -count=5

# tidy dependencies
tidy:
	go mod tidy

# clean build artifacts
clean:
	rd /s /q $(BUILD_DIR)

migrate-force:
	$(MIGRATE) force $(version)

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down -all

migrate-version:
	$(MIGRATE) version

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# show available commands
help:
	@echo   Available commands:
	@echo   make            			- format all code, build binary and run
	@echo   make format     			- format all code
	@echo   make build      			- build binary
	@echo   make test                   - run all tests
	@echo   make tidy       			- tidy go modules
	@echo   make clean      			- remove build artifacts
	@echo   make run        			- run the compiled executable
	@echo   make migrate-force          - clear dirty database version
	@echo   make migrate-up             - apply all pending migrations
	@echo   make migrate-down           - roll back the most recent migration
	@echo   make migrate-version        - show current migration version
	@echo   make migrate-create name=x  - create a new migration pair
