# For local development

# Build the executable for current platform with whaletui as the name
build:
	go build -o whaletui .

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

imports:
	goimports -w .

vet:
	go vet ./...

test-all: vet imports fmt lint test
