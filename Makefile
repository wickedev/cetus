all: install
	go build -o dist/main

install:
	go mod vendor

run: install
	go run main.go

clean:
	rm -rf dist vendor