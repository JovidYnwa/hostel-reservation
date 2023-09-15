run:
	./bin/api

build:
	go build -o bin/api

test:
	got test -v ./...