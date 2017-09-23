NAME     := beacon
VERSION  := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS  := -ldflags="-X \"main.version=$(VERSION)\" -X \"main.revision=$(REVISION)\""

bin/$(NAME): format deps
	go build $(LDFLAGS) -o bin/$(NAME)

linux: format
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(NAME)

format:
	go fmt

clean:
	rm -rf bin/*

install:
	go install $(LDFLAGS)

deps:
	glide install

update:
	glide update

.PHONY: format, clean, install, deps, update
