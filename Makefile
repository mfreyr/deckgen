BINARY_NAME=deckgen
BIN_DIR=bin

.PHONY: build dev clean deps

build: deps
	@echo "Tidying go module dependencies..."
	@go mod tidy
	@echo "Generating Go code from templ files..."
	templ generate
	@echo "Building production binary..."
	go build -ldflags="-w -s" -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/server/main.go

dev: deps
	@echo "Starting development server with live reload..."
	templ generate --watch --cmd="go run ./cmd/server/main.go"

deps:
	@echo "Installing development dependencies..."
	@go install github.com/a-h/templ/cmd/templ@latest

clean:
	@echo "Cleaning up..."
	@rm -f $(BIN_DIR)/$(BINARY_NAME)
	@go clean
