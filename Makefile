# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=wio
BINARY_UNIX=$(BINARY_NAME)_unix

all: build run

get:
	@echo Getting Required tools to build this project
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/kardianos/govendor
	@echo Tools Downloaded and Built!!

build:
	@echo Building $(BINARY_NAME) project:
	@cd "$(CURDIR)/cmd/$(BINARY_NAME)/utils/io" && go-bindata -nomemcopy -pkg io -prefix ../../../../ ../../../../assets/...
	@cd "$(CURDIR)/cmd/$(BINARY_NAME)" && $(GOBUILD) -o ../../bin/$(BINARY_NAME) -v
	@echo Project built!!

clean:
	@echo Cleaning $(BINARY_NAME) project files:
	@$(GOCLEAN)
	@rm -f bin/$(BINARY_NAME)
	@rm -f bin/$(BINARY_UNIX)
	@echo Cleaning Finished!!

run:
	@./bin/$(BINARY_NAME) ${ARGS}