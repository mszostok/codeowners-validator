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
# Go modules
##
shout "? go mod tidy"
go mod tidy
STATUS=$( git status --porcelain go.mod go.sum )
if [ ! -z "$STATUS" ]; then
    echo "${RED}✗ go mod tidy modified go.mod and/or go.sum${NC}"
    exit 1
else echo -e "${GREEN}√ go mod tidy${NC}"
fi

##
# GO BUILD
##
buildEnv=""
if [[ "${RUN_ON_CI:-x}" == "true" ]]; then
	# build binary statically
	buildEnv="env CGO_ENABLED=0"
fi
shout "? go build"
${buildEnv} go build -o codeowners-validator ./main.go
goBuildResult=$?
if [[ ${goBuildResult} != 0 ]]; then
    echo -e "${RED}✗ go build ${NC}\n $goBuildResult${NC}"
    exit 1
else echo -e "${GREEN}√ go build ${NC}"
fi

##
# GO TEST
##
shout "? go test"
go test -race ./...
# Check if tests passed
if [[ $? != 0 ]]; then
	echo -e "${RED}✗ go test\n${NC}"
	exit 1
else echo -e "${GREEN}√ go test${NC}"
fi
