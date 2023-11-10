BINARY_NAME=genz
SUPPORTED_OS=linux darwin windows
SUPPORTED_ARCH=amd64 arm64

install:
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.1
	go install .

lint:
	golangci-lint run

test:
	go test -v ./...

build-all:
	for os in $(SUPPORTED_OS); do \
		for arch in $(SUPPORTED_ARCH); do \
			GOOS=$$os GOARCH=$$arch go build -o bin/$(BINARY_NAME)-$$os-$$arch; \
		done; \
	done

build:
	go build -o bin/$(BINARY_NAME)