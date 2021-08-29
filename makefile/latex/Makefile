help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help

MAIN_FILE:=cv.tex
OUTPUT_FORMAT:=pdf
PDF_VIEWER:=zathura

DOCKER_IMAGE:=latex-builder
USER:=$$(id -u -n)
UID:=$$(id -u)

docker-image: req-docker ## build the docker image to build the CV
	@docker build --build-arg USER_NAME=$(USER) --build-arg UID=$(UID) . -t $(DOCKER_IMAGE)


docker-compile: req-docker ## compile to pdf format using the Docker image
	@docker run -it --rm -v $(PWD):/$(HOME) -w $(HOME) $(DOCKER_IMAGE) ./build.sh $(MAIN_FILE)

compile: req-tectonic ## compile using tectonic
	@tectonic -c minimal --outfmt $(OUTPUT_FORMAT) $(MAIN_FILE)

watch: req-entr ## watch tex files for changes and automatically compile
	@(find . -type f -name "*.tex" | entr -cr make compile 2&>1 /dev/null &)

unwatch: ## stop watching tex files
	kill $$(ps aux | grep entr | grep -v grep | awk '{print $$2}')

open: req-viewer ## open cv.pdf
	@$(shell $(PDF_VIEWER) $(basename $(MAIN_FILE)).$(OUTPUT_FORMAT) > /dev/null 2>&1 &)

req-docker:
	@docker info > /dev/null 2>&1 || { echo >&2 "require running docker"; exit 1; }
req-tectonic:
	@command -v tectonic >/dev/null 2>&1 || { echo >&2 "require tectonic"; exit 1; }
req-entr:
	@command -v entr >/dev/null 2>&1 || { echo >&2 "require entr"; exit 1; }
req-viewer:
	@command -v $(PDF_VIEWER) > /dev/null 2>&1 || { echo >&2 "require $(PDF_VIEWER)"; exit 1; }


