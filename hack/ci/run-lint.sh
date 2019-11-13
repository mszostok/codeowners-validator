#!/usr/bin/env bash

# standard bash error handling
set -o nounset # treat unset variables as an error and exit immediately.
set -o errexit # exit immediately when a command fails.
set -E         # needs to be set if we want the ERR trap

readonly CURRENT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
readonly ROOT_PATH=${CURRENT_DIR}/../..
readonly GOLANGCI_LINT_VERSION="v1.21.0"

source "${CURRENT_DIR}/utilities.sh" || { echo 'Cannot load CI utilities.'; exit 1; }

golangci::install() {
  shout "Install the golangci-lint in version ${GOLANGCI_LINT_VERSION}"
  curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b "$GOPATH"/bin ${GOLANGCI_LINT_VERSION}

  echo -e "${GREEN}√ install golangci-lint${NC}"
}

golangci::run_checks() {
  shout "Run golangci-lint checks"
  LINTS=(
    # default golangci-lint lints
    deadcode errcheck gosimple govet ineffassign staticcheck \
    structcheck typecheck unused varcheck \
    # additional lints
    golint gofmt misspell gochecknoinits unparam scopelint gosec
  )

  ENABLE=$(sed 's/ /,/g' <<< "${LINTS[@]}")

  golangci-lint --disable-all --enable="${ENABLE}" run ./...

  echo -e "${GREEN}√ run golangci-lint${NC}"
}

main() {
  if [[ "${RUN_ON_CI:-x}" == "true" ]]; then
    golangci::install
  fi

  golangci::run_checks
}

main
