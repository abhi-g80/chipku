BIN="./bin"
SRC=$(shell find . -name "*.go")

ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

ifeq (, $(shell which richgo))
$(warning "could not find richgo in $(PATH), run: go get github.com/kyoh86/richgo")
endif

BUILD = build
BINARY = chipku
LDFLAGS = "-s -w"

.PHONY: fmt lint test install_deps clean run

default: all

all: fmt test

fmt:
	$(info **************** checking formatting **************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info **************** running lint tools ***************)
	golangci-lint run -v

richtest: install_deps
	$(info **************** running tests - rich ***********)
	richgo test -v ./...

test: install_deps
	$(info **************** running tests ********************)
	go test -v ./...

install_deps:
	$(info **************** downloading dependencies *********)
	go get -v ./...

build: clean
	$(info **************** building binaries ****************)
	mkdir $(BUILD)
	go build -v -ldflags=$(LDFLAGS) -o $(BUILD)/$(BINARY)

clean:
	$(info **************** house keeping ********************)
	rm -rf $(BIN)
	rm -rf $(BUILD)

up:
	$(info **************** docker build + up ****************)
	docker-compose up --build --remove-orphans --detach

run:
	$(info **************** run  *****************************)
	go run . serve
