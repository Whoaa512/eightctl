.PHONY: fmt lint test

fmt:
	gofumpt -w ./

lint:
	golangci-lint run ./...

test:
	go test ./...
