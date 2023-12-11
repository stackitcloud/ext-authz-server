# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

VERSION                                := $(shell cat VERSION)
REGISTRY                               := eu.gcr.io/gardener-project/gardener
PREFIX                                 := ext-authz-server
EXTERNAL_AUTHZ_SERVER_IMAGE_REPOSITORY := $(REGISTRY)/$(PREFIX)
EXTERNAL_AUTHZ_SERVER_IMAGE_TAG        := $(VERSION)
REPO_ROOT                              := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

#########################################
# Tools                                 #
#########################################

TOOLS_DIR := hack/tools
include vendor/github.com/gardener/gardener/hack/tools.mk

#################################################################
# Rules related to binary build, Docker image build and release #
#################################################################

.PHONY: ext-authz-server-docker-image
ext-authz-server-docker-image:
	@docker build -t $(EXTERNAL_AUTHZ_SERVER_IMAGE_REPOSITORY):$(EXTERNAL_AUTHZ_SERVER_IMAGE_TAG) -f Dockerfile --rm .
.PHONY: docker-images
docker-images: ext-authz-server-docker-image

.PHONY: release
release: docker-images docker-login docker-push

.PHONY: docker-login
docker-login:
	@gcloud auth activate-service-account --key-file .kube-secrets/gcr/gcr-readwrite.json

.PHONY: docker-push
docker-push:
	@if ! docker images $(EXTERNAL_AUTHZ_SERVER_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(EXTERNAL_AUTHZ_SERVER_IMAGE_TAG); then echo "$(EXTERNAL_AUTHZ_SERVER_IMAGE_REPOSITORY) version $(EXTERNAL_AUTHZ_SERVER_IMAGE_TAG) is not yet built. Please run 'ext-authz-server-docker-image'"; false; fi
	@docker -- push $(EXTERNAL_AUTHZ_SERVER_IMAGE_REPOSITORY):$(EXTERNAL_AUTHZ_SERVER_IMAGE_TAG)

#####################################################################
# Rules for verification, formatting, linting, testing and cleaning #
#####################################################################

.PHONY: revendor
revendor:
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod vendor
	@chmod +x $(REPO_ROOT)/vendor/github.com/gardener/gardener/hack/*

.PHONY: check
check: $(GOIMPORTS) $(GOLANGCI_LINT)
	@$(REPO_ROOT)/vendor/github.com/gardener/gardener/hack/check.sh ./cmd/... ./pkg/...

.PHONY: format
format: $(GOIMPORTS) $(GOIMPORTSREVISER)
	@$(REPO_ROOT)/vendor/github.com/gardener/gardener/hack/format.sh ./cmd ./pkg

.PHONY: test
test:
	@$(REPO_ROOT)/vendor/github.com/gardener/gardener/hack/test.sh ./cmd/... ./pkg/...

.PHONY: test-cov
test-cov:
	@$(REPO_ROOT)/vendor/github.com/gardener/gardener/hack/test-cover.sh ./cmd/... ./pkg/...

.PHONY: test-cov-clean
test-cov-clean:
	@$(REPO_ROOT)/vendor/github.com/gardener/gardener/hack/test-cover-clean.sh

.PHONY: verify
verify: check format test

.PHONY: verify-extended
verify-extended: check format test-cov test-cov-clean
