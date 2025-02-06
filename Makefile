BINARY_NAME=updock
BUILD_DIR=bin

all: build

build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go