#!/usr/bin/env bash

# standard bash error handling
set -o nounset # treat unset variables as an error and exit immediately.
set -o errexit # exit immediately when a command fails.
set -E         # needs to be set if we want the ERR trap

# Currently we are using the newest go (1.13) with go modules
export GO111MODULE=on

readonly CURRENT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
readonly ROOT_PATH=${CURRENT_DIR}/../..

source "${CURRENT_DIR}/utilities.sh" || { echo 'Cannot load CI utilities.'; exit 1; }

pushd ${ROOT_PATH} > /dev/null

# Exit handler. This function is called anytime an EXIT signal is received.
# This function should never be explicitly called.
function _trap_exit () {
    popd > /dev/null
}
trap _trap_exit EXIT

echo -e "${INVERTED}"
echo "USER: " + ${USER}
echo "PATH: " + ${PATH}
echo "GOPATH:" + ${GOPATH}
echo -e "${NC}"

##
# GO TEST INTEGRATION
##
shout "? go test integration"
go test ./tests/integration/... -v -tags=integration
# Check if tests passed
if [[ $? != 0 ]]; then
	echo -e "${RED}✗ go test integration\n${NC}"
	exit 1
else echo -e "${GREEN}√ go test integration${NC}"
fi
