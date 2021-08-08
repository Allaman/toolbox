.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

REVEALJS_URL:=https://unpkg.com/reveal.js@4.1.0/
MARKDOWN_FILE:=presentation.md
OUTPUT_FILE:=index.html

build: ## Build the revealjs presentation
ifneq ("$(wildcard $(./custom.css))","")
	@echo Building slides with custom.css
	pandoc -t revealjs -s -o $(OUTPUT_FILE) $(MARKDOWN_FILE) -V revealjs-url=$(REVEALJS_URL) --css=custom.css --slide-level=2 --self-contained
else
	@echo Building slides without custom.css
	pandoc -t revealjs -s -o $(OUTPUT_FILE) $(MARKDOWN_FILE) -V revealjs-url=$(REVEALJS_URL) --slide-level=2 --self-contained
endif

req-pandoc:
	@command -v pandoc >/dev/null 2>&1 || { echo >&2 "require pandoc"; exit 1; }


