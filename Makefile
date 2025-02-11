BINARY_NAME=updock
BUILD_DIR=bin
ENTRY_POINT=src/main.go

#all: build
all: run

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(ENTRY_POINT)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)