build:
	go build -o bin/cdex

run: build
	./bin/dex

test:
	go test -v ./...
