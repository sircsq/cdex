build:
	go build -o bin/cdex

run: build
	./bin/cdex

test:
	go test -v ./...
