#!/usr/bin/env bash

BASE_DIR=$(cd "$(dirname ${BASH_SOURCE[0]})" && pwd)

COVERAGE_DIR="${BASE_DIR}/.coverage"
OUT_TXT="out.txt"
OUT_HTML="out.html"

# -- 0.
if [ ! -d "${COVERAGE_DIR}" ]; then
  mkdir "${COVERAGE_DIR}"
fi
# -- 1. coverage.txt
go test -race -coverprofile="${COVERAGE_DIR}/${OUT_TXT}" -covermode=atomic
# -- 2. coverage.html
go tool cover -html "${COVERAGE_DIR}/${OUT_TXT}" -o "${COVERAGE_DIR}/${OUT_HTML}"
