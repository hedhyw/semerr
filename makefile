FILES_DIR?=$(PWD)/pkg/v1
GOLANGCI_LINT_VER:=v1.42.1

all: generate prepare lint test
.PHONY: all

generate:
	go generate ./semerr.go
.PHONY: generate

test:
	go test ./...
.PHONY: test

lint:
	./bin/golangci-lint run --timeout=10m ./...
.PHONY: lint

prepare:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b ./bin $(GOLANGCI_LINT_VER)
.PHONY: install
