#!/usr/bin/env bash
# Inspired by https://liam.sh/post/makefiles-for-go-projects

# standard bash error handling
set -o nounset # treat unset variables as an error and exit immediately.
set -o errexit # exit immediately when a command fails.
set -E         # needs to be set if we want the ERR trap

readonly CURRENT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
readonly ROOT_PATH=$( cd "${CURRENT_DIR}/../.." && pwd )
readonly GOLANGCI_LINT_VERSION="v1.23.8"

source "${CURRENT_DIR}/utilities.sh" || { echo 'Cannot load CI utilities.'; exit 1; }

# This will find all files (not symlinks) with the executable bit set:
# https://apple.stackexchange.com/a/116371
binariesToCompress=$(find "${ROOT_PATH}/dist" -perm +111 -type f)

shout "Staring compression for: \n$binariesToCompress"

command -v upx > /dev/null || { echo 'UPX binary not found, skipping compression.'; exit 1; }

# I just do not like playing with xargs ¯\_(ツ)_/¯
for i in $binariesToCompress
do
  upx --brute "$i"
done
