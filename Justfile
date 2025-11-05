# chronogo development automation

# Import apptainer module
mod apptainer

# List all recipes
default:
    @just --list

# Run all tests
test:
    go test ./...

# Run linter with all checks enabled
lint:
    golangci-lint run --enable-all ./...
