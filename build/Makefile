BINARY_NAME=assessment

# Go source files
SRC_FILES=../src/main.go ../src/json.go

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o ../bin/$(BINARY_NAME) $(SRC_FILES)

clean:
	@echo "Cleaning up..."
	@rm -f ../bin/$(BINARY_NAME)

# Phony targets (avoids conflicts with files)
.PHONY: all build clean

