# Current addon version
VERSION ?= $$(git describe --abbrev=0 --tags --match "v*")

# Image URL to use all building/pushing image targets
IMG ?= clastix/capsule-rancher-addon:$(VERSION)

# Get information about git current status
GIT_HEAD_COMMIT ?= $$(git rev-parse --short HEAD)
GIT_TAG_COMMIT  ?= $$(git rev-parse --short $(VERSION))
GIT_MODIFIED_1  ?= $$(git diff $(GIT_HEAD_COMMIT) $(GIT_TAG_COMMIT) --quiet && echo "" || echo ".dev")
GIT_MODIFIED_2  ?= $$(git diff --quiet && echo "" || echo ".dirty")
GIT_MODIFIED    ?= $$(echo "$(GIT_MODIFIED_1)$(GIT_MODIFIED_2)")
GIT_REPO        ?= $$(git config --get remote.origin.url)
BUILD_DATE      ?= $$(git log -1 --format="%at" | xargs -I{} date -d @{} +%Y-%m-%dT%H:%M:%S)

# Build the container image
image:
	docker build . -t $(IMG) --build-arg GIT_HEAD_COMMIT=$(GIT_HEAD_COMMIT) \
		--build-arg GIT_TAG_COMMIT=$(GIT_TAG_COMMIT) \
		--build-arg GIT_MODIFIED=$(GIT_MODIFIED) \
		--build-arg GIT_REPO=$(GIT_REPO) \
		--build-arg GIT_LAST_TAG=$(VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE)

image/dlv: VERSION := dlv
image/dlv:
	docker build . --build-arg "GCFLAGS=all=-N -l" --tag $(IMG) --target dlv

# Push the container image
image/push: image
	docker push $(IMG)

image/dlv/push: VERSION := dlv
image/dlv/push: image/dlv
	docker push $(IMG)

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Generate code
CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.10.0)

# go-install-tool will 'go install' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-install-tool
@[ -f $(1) ] || { \
set -e ;\
echo "Installing $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
}
endef