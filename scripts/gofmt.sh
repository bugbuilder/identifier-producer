#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_PATH=$(dirname "$(readlink -f "$BASH_SOURCE")")
cd $SCRIPT_PATH/..

find_files() {
  find . -not \( \
      \( \
        -wholename './.git' \
        -o -wholename '*/vendor/*' \
      \) -prune \
    \) -name '*.go'
}

GOFMT="gofmt -s -l -w"
bad_files=$(find_files | xargs $GOFMT)
if [[ -n "${bad_files}" ]]; then
  echo "!!! '$GOFMT' run on the following files: "
  echo "${bad_files}"
fi