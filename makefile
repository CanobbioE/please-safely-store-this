.PHONY: test
GO111MODULE=on
PKG_NAME=.

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))

USER_GROUP := $(if $(USER), $$(id -u ${USER}):$$(id -g ${USER}) , 0:0)

default:

fmt:
	@echo "Fixing code with gofmt..."
	gofmt -s -w $$(go list -f "{{.Dir}}" ./...)

fieldalignment:
	@echo "Analyzer structs and rearranged to use less memory with fieldalignment..."
	fieldalignment -fix -test=false ./...

lint:
	@echo "Checking code against linters..."
	@golangci-lint run --new-from-rev=$$(git merge-base HEAD master) --timeout 6m0s ./...

gci:
	@echo "Fixing imports with gci..."
	@gci write -s standard -s default -s "prefix(github.com/CanobbioE/please-safely-store-this)" -s blank -s dot ./cmd ./internal

install-tools:
	@echo "Installing tools..."
	@go install github.com/daixiang0/gci@latest
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	@echo "Installation completed!"