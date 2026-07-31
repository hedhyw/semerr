FILES_DIR?=$(PWD)/pkg/v1
GOLANGCI_LINT_VER:=v2.12.2

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
	GOBIN=$(PWD)/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VER)
