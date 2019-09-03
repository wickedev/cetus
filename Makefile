all: cli

cli:
	go build -o dist/cetus pkg/cli/main.go

install:
	go mod vendor

run: install
	go run main.go

clean:
	rm -rf dist vendor