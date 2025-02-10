BINARY_NAME=assessment

# Go source files
SRC_FILES=main.go json.go

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(SRC_FILES)

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# Phony targets (avoids conflicts with files)
.PHONY: all build clean

