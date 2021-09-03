.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

DIFFCMD:=nvim -d

EDITOR:=$(shell type -p nvim || echo vim)
ifneq (,$(findstring vim,$(EDITOR)))
	EDITOR += -c 'set syntax=yaml' -
else
	EDITOR:=less
endif

BUILD:=kustomize build

check-stg: req-kustomize ## Generate staging manifests - do not apply them
	@$(BUILD) staging/ | $(EDITOR)
check-prod: req-kustomize ## Generate production manifests - do not apply them
	@$(BUILD) production/ | $(EDITOR)

val-stg: req-kubeval req-kustomize ## Validate staging Kubernetes manifests - do not apply them
	@$(BUILD) staging/ | kubeval --strict --ignore-missing-schemas
val-prod: req-kubeval req-kustomize ## Validate staging Kubernetes manifests - do not apply them
	@$(BUILD) production/ | kubeval --strict --ignore-missing-schemas

analyze-stg: req-kustomize req-kube-linter ## Static analysis on staging manifests - do not apply them
	@$(BUILD) staging/ > stg.yaml
	@kube-linter lint stg.yaml || true
	@rm stg.yaml
analyze-prod: req-kustomize req-kube-linter ## Static analysis on production manifests - do not apply them
	@$(BUILD) production/ > prod.yaml
	@kube-linter lint prod.yaml || true
	@rm prod.yaml

diff: ## Diff prod and stg manifests
	@$(BUILD) staging/ > stg.yaml
	@$(BUILD) production/ > prod.yaml
	@$(DIFFCMD) stg.yaml prod.yaml
	@rm stg.yaml prod.yaml

req-kustomize:
	@command -v kustomize >/dev/null 2>&1 || { echo >&2 "require kustomize"; exit 1; }
req-kubeval:
	@command -v kubeval >/dev/null 2>&1 || { echo >&2 "require kubeval"; exit 1; }
req-kube-linter:
	@command -v kube-linter >/dev/null 2>&1 || { echo >&2 "require kube-linter"; exit 1; }

