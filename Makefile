.PHONY: run install clean

all: build

build:
	go build -o ../dist/cetus main.go

run:
	go run cli/main.go

test:
	go test **/*_test.go

install:
ifeq ($(GOPATH),)
	@echo GOPATH is empty
else
	go build -o ${GOPATH}/bin/cetus main.go
	@echo cetus cli installed on ${GOPATH}/bin/cetus
endif

clean:
	rm -rf ../dist/cetus