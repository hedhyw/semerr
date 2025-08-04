FILES_DIR?=$(PWD)/pkg/v1
GOLANGCI_LINT_VER:=v2.3.1

all: generate lint test
.PHONY: all

generate:
	go generate ./semerr.go
.PHONY: generate

test:
	go test \
		-coverpkg=github.com/hedhyw/semerr/pkg/... \
		-covermode=count \
		-coverprofile=coverage.out \
		./...
	go tool cover -func=coverage.out
.PHONY: test

lint: bin/golangci-lint
	./bin/golangci-lint run --timeout=10m ./...
.PHONY: lint

bin/golangci-lint:
	curl \
		-sSfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s $(GOLANG_CI_LINT_VER)
