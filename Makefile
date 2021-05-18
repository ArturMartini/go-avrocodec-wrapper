.DEFAULT_GOAL := test

.PHONY: help clean cover cover-func cover-brower test doc fmt

DOC_ADDR             := :8081
BIN_PATH             := dist
COVER_FILE_PATH      := $(BIN_PATH)/coverage.out
TESTS_PATH           := ./...

help:  ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Clean up build files
	rm -rf dist

$(BIN_PATH):
	mkdir -p $(BIN_PATH)

$(COVER_FILE_PATH): $(BIN_PATH)
	go test -coverprofile=$(COVER_FILE_PATH) $(TESTS_PATH)

cover:  ## Run coverage tests
	go test -cover $(TESTS_PATH)
	
cover-func: $(COVER_FILE_PATH)  ## Run coverage tests by function
	go tool cover -func=$(COVER_FILE_PATH)

cover-file: $(COVER_FILE_PATH) ## Create a file with coverage test

cover-browser: cover-file ## Show coverage test in a browser
	go tool cover -html=$(COVER_FILE_PATH)

test: ## Run all unit tests
	go test $(TESTS_PATH)

doc: ## Start a go doc server, need to have installed go tools: go get -u golang.org/x/tools/...
	godoc -http $(DOC_ADDR)

fmt: ## Format all code with gofmt
	gofmt -s -w .
