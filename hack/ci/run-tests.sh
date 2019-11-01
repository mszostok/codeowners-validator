#!/usr/bin/env bash

# standard bash error handling
set -o nounset # treat unset variables as an error and exit immediately.
set -o errexit # exit immediately when a command fails.
set -E         # needs to be set if we want the ERR trap

# Currently we are using the newest go (1.13) but project is still not switched to go modules
export GO111MODULE=off

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
# DEP STATUS
##
shout "? dep status"
depResult=$(dep status -v)
if [[ $? != 0 ]]; then
	echo -e "${RED}✗ dep status\n$depResult${NC}"
	exit 1
else echo -e "${GREEN}√ dep status${NC}"
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
go test ./...
# Check if tests passed
if [[ $? != 0 ]]; then
	echo -e "${RED}✗ go test\n${NC}"
	exit 1
else echo -e "${GREEN}√ go test${NC}"
fi

goFilesToCheck=$(find . -type f -name "*.go" | egrep -v "\/vendor\/|_*/automock/|_*/testdata/|/pkg\/|_*export_test.go")

##
# GO IMPORTS & FMT
##
go build -o goimports-vendored ./vendor/golang.org/x/tools/cmd/goimports
buildGoImportResult=$?
if [[ ${buildGoImportResult} != 0 ]]; then
	echo -e "${RED}✗ go build goimports${NC}\n$buildGoImportResult${NC}"
	exit 1
fi

shout "? goimports"
goImportsResult=$(echo "${goFilesToCheck}" | xargs -L1 ./goimports-vendored -w -l)
rm goimports-vendored

if [[ $(echo ${#goImportsResult}) != 0 ]]; then
	echo -e "${RED}✗ goimports and fmt ${NC}\n$goImportsResult${NC}"
	exit 1
else echo -e "${GREEN}√ goimports and fmt ${NC}"
fi
