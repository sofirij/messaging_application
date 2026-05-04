# variables
BINARY_NAME=app.exe
BUILD_DIR=./bin
MAIN_PATH=./cmd/server/main.go

# default target — runs when you just type "make"
all: format tidy build run

# format code
format:
	gofmt -w .

# format and fix imports
imports:
	goimports -w .

# build the binary
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# run the binary
run:
	-$(BUILD_DIR)/$(BINARY_NAME)

# tidy dependencies
tidy:
	go mod tidy

# clean build artifacts
clean:
	rd /s /q $(BUILD_DIR)

# show available commands
help:
	@echo Available commands:
	@echo   make            - format all code, build binary and run
	@echo   make format     - format all code
	@echo   make build      - build binary
	@echo   make tidy       - tidy go modules
	@echo   make clean      - remove build artifacts
	@echo   make run        - run the compiled executable