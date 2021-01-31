.DEFAULT_GOAL = all

ROOT_DIR:=$(dir $(abspath $(firstword $(MAKEFILE_LIST))))

# enable module support across all go commands.
export GO111MODULE = on
# enable consistent Go 1.12/1.13 GOPROXY behavior.
export GOPROXY = https://proxy.golang.org

all: build-race test-unit test-integration test-lint
.PHONY: all

# When running integration tests on windows machine
# it cannot execute binary without extension.
# It needs to be parametrized, so we can override it on CI.
export BINARY_PATH = $(ROOT_DIR)/codeowners-validator$(BINARY_EXT)

############
# Building #
############

build:
	go build -o $(BINARY_PATH) ./main.go
.PHONY: build

build-race:
	go build -race -o codeowners-validator ./main.go
.PHONY: build-race

###########
# Testing #
###########

test-unit:
	./hack/run-test-unit.sh
.PHONY: test-unit

test-integration: build
	./hack/run-test-integration.sh
.PHONY: test-integration

test-lint:
	./hack/run-lint.sh
.PHONY: test-lint

test-hammer:
	go test -count=100 ./...
.PHONY: test-hammer

test-unit-cover-html: test-unit
	go tool cover -html=./coverage.txt
.PHONY: cover-html

###############
# Development #
###############

fix-lint-issues:
	LINT_FORCE_FIX=true ./hack/run-lint.sh
.PHONY: fix-lint
