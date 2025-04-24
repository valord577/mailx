#!/usr/bin/env bash
set -e

PROJ_ROOT=$(cd "$(dirname ${BASH_SOURCE[0]})"; pwd)

# -- 1. qa_coverage.txt
go test -race -coverprofile="${PROJ_ROOT}/qa_coverage.txt" -covermode=atomic ${PROJ_ROOT}/...
# -- 2. qa_coverage.html
go tool cover -html "${PROJ_ROOT}/qa_coverage.txt" -o "${PROJ_ROOT}/qa_coverage.html"
