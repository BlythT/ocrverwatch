# Makefile

.PHONY: lint test

# Linting target
lint:
	@echo "Running linter..."
	golangci-lint run -v ./...

test:
	go test -v ./...

