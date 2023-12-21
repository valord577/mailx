#!/usr/bin/env bash

PROJ_ROOT=$(cd "$(dirname ${BASH_SOURCE[0]})"; pwd)

COVERAGE="${PROJ_ROOT}/.coverage"

# -- 0.
if [ ! -d "${COVERAGE}" ]; then (mkdir -p "${COVERAGE}") fi
# -- 1. coverage.txt
go test -race -coverprofile="${COVERAGE}/coverage.txt" -covermode=atomic ${PROJ_ROOT}/...
# -- 2. coverage.html
go tool cover -html "${COVERAGE}/coverage.txt" -o "${COVERAGE}/coverage.html"
