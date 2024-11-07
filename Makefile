build:
	@go build -o bin/nodeBox

run: build
	@./bin/nodeBox

test:
	@go test ./...
