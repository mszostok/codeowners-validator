.DEFAULT_GOAL = build

# enable module support across all go commands.
export GO111MODULE = on
# enable consistent Go 1.12/1.13 GOPROXY behavior.
export GOPROXY = https://proxy.golang.org

# Build

build:
	go build -o codeowners-validator ./main.go
.PHONY: build

build-race:
	go build -race -o codeowners-validator ./main.go
.PHONY: build-race

# Test
test-unit:
	./hack/ci/run-test-unit.sh
.PHONY: test-unit

test-integration: build
	env BINARY_PATH=$(PWD)/codeowners-validator ./hack/ci/run-test-integration.sh
.PHONY: test-integration

test-lint:
	./hack/ci/run-lint.sh
.PHONY: test-lint

test-hammer:
	go test -count=100 ./...
.PHONY: test-hammer

cover-html:
	go test -v -coverprofile=./tmp/coverage.out ./...
	go tool cover -html=./tmp/coverage.out
.PHONY: cover-html
