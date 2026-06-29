.PHONY: all lint test cover build clean

all: lint test build

lint:
	golangci-lint run ./...

test:
	go test -count=1 -race ./...

cover:
	go test -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

build:
	go build -trimpath -ldflags="-s -w" -o scraper ./cmd/scraper/

clean:
	rm -f scraper coverage.out coverage.html
