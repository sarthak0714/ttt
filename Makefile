
# Templ Gen
generate:
	@templ generate

# Build the server
build: generate
	@mkdir -p bin
	@go build -o bin/app server/main.go

# Run the server
run:build
	@bin/app

# Clean 
clean:
	@rm -f bin/server
	@rm -rf templates/generated

# Install dependencies
deps:
	@go mod tidy

test:
	@go test -v ./...
