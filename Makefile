.PHONY: build test lint clean

build:
	mkdir -p bin
	go build -o bin/soch ./cmd/soch

test:
	go test -shuffle=on -v ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin
