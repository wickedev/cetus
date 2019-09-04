.PHONY: install clean

all: cli

cli:
	go build -o dist/cetus pkg/cli/main.go

install:
ifeq ($(GOPATH),)
	@echo GOPATH is empty
else
	go build -o ${GOPATH}/bin/cetus pkg/cli/main.go
	@echo cetus cli installed on ${GOPATH}/bin/cetus
endif

clean:
	rm -rf dist