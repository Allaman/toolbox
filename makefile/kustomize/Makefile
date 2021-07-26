.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

STAGING_PATH:=staging/
PRODUCTION_PATH:=production/
EDITOR:=$(shell type -p nvim || echo vim)
ifneq (,$(findstring vim,$(EDITOR)))
	EDITOR += -c 'set syntax=yaml' -
else
	EDITOR:=less
endif

BUILD:=kustomize build

check-stg: req-kustomize ## Generate staging manifests - do not apply them
	@$(BUILD) ${STAGING_PATH} | $(EDITOR)
check-prod: req-kustomize ## Generate production manifests - do not apply them
	@$(BUILD) ${PRODUCTION_PATH} | $(EDITOR)

val-stg: req-kubeval req-kustomize ## Validate staging Kubernetes manifests - do not apply them
	@$(BUILD) ${STAGING_PATH} | kubeval --strict --ignore-missing-schemas
val-prod: req-kubeval req-kustomize ## Validate staging Kubernetes manifests - do not apply them
	@$(BUILD) ${PRODUCTION_PATH} | kubeval --strict --ignore-missing-schemas

analyze-stg: req-kustomize req-kube-linter ## Static analysis on staging manifests - do not apply them
	@$(BUILD) ${STAGING_PATH} > stg.yaml
	@kube-linter lint stg.yaml
analyze-prod: req-kustomize req-kube-linter ## Static analysis on production manifests - do not apply them
	@$(BUILD) ${PRODUCTION_PATH} > prod.yaml
	@kube-linter lint prod.yaml

req-kustomize:
	@command -v kustomize >/dev/null 2>&1 || { echo >&2 "require kustomize"; exit 1; }
req-kubeval:
	@command -v kubeval >/dev/null 2>&1 || { echo >&2 "require kubeval"; exit 1; }
req-kube-linter:
	@command -v kube-linter >/dev/null 2>&1 || { echo >&2 "require kube-linter"; exit 1; }
