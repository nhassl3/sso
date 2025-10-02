.PHONY: build run runb clean

BINARY_NAME := ssoapp
BUILD_DIR := build
MAIN_PACKAGE := ./cmd/ssoapp # Adjust if your main package is elsewhere

# Target to build the Go application
build:
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Target to run the application
build-run: build
	@$(BUILD_DIR)/$(BINARY_NAME) --config="./config/local.yaml"

run:
	@$(BUILD_DIR)/$(BINARY_NAME) --config="./config/local.yaml"

clean:
	@rm -rf $(BUILD_DIR)

-DEFAULT-GOAL: run