.PHONY: build test clean

build:
	mkdir -p bin
	go build -o bin/soch ./cmd/soch

test:
	go test -shuffle=on -v ./...

clean:
	rm -rf bin
