## ----------------------------------------------------------------------
## This makefile can be used to execute common functions to interact with
## the source code, these functions ease local development and can also be
## used in CI/CD pipelines.
## ----------------------------------------------------------------------

linter_version=v2.12.2

.PHONY: help check-lint lint serve-godoc test

# REFERENCE: https://stackoverflow.com/questions/16931770/makefile4-missing-separator-stop
help: ## - Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

check-lint: ## - validate/install golangci-lint installation
	@which golangci-lint &> /dev/null || (curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin ${linter_version})

lint: check-lint ## - lint the source
	@golangci-lint run --verbose

serve-godoc: ## - serve (web) the godocs
	go doc -http

test: ## - test the source with verbose output
	@go test -cover -v --count=1 ./...
