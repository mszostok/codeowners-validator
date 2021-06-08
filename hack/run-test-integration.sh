#!/usr/bin/env bash

# standard bash error handling
set -o nounset # treat unset variables as an error and exit immediately.
set -o errexit # exit immediately when a command fails.
set -E         # needs to be set if we want the ERR trap

CURRENT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ROOT_PATH=$( cd "${CURRENT_DIR}/.." && pwd )
readonly CURRENT_DIR
readonly ROOT_PATH

# shellcheck source=./hack/lib/utilities.sh
source "${CURRENT_DIR}/lib/utilities.sh" || { echo 'Cannot load CI utilities.'; exit 1; }

pushd "${ROOT_PATH}" > /dev/null

# Exit handler. This function is called anytime an EXIT signal is received.
# This function should never be explicitly called.
function _trap_exit () {
    popd > /dev/null
}
trap _trap_exit EXIT

function print_info() {
  echo -e "${INVERTED}"
  echo "USER: ${USER:-"unknown"}"
  echo "PATH: ${PATH:-"unknown"}"
  echo "GOPATH: ${GOPATH:-"unknown"}"
  echo -e "${NC}"
}

function test::integration() {
  shout "? go test integration"

  # Check if tests passed
  # shellcheck disable=SC2046
  if ! go test -v -tags=integration ./tests/integration/... $(test::update_golden);
  then
    echo -e "${RED}✗ go test integration\n${NC}"
    exit 1
  else
    echo -e "${GREEN}√ go test integration${NC}"
  fi
}

function test::update_golden() {
    if [[ "${UPDATE_GOLDEN:-"false"}" == "true" ]]; then
      echo "-update"
    fi
}

function main() {
    print_info

    test::integration
}

main
