BINARY := guacamole-pre-commit

.PHONY: build test lint

build:
	go build -o $(BINARY) .

test:
	go test ./...

lint:
	gofmt -l .
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed: https://golangci-lint.run/welcome/install/" && exit 1)
	golangci-lint run
