#!/usr/bin/env bash
set -e

PROJ_ROOT=$(cd "$(dirname ${BASH_SOURCE[0]})"; pwd)

export GO111MODULE="on"
export CGO_ENABLED="1"
export GOPROXY="https://goproxy.cn,direct"
export GOSUMDB="sum.golang.google.cn"

# -- 1. coverage.txt
go test -race -coverprofile="${PROJ_ROOT}/coverage.txt" -covermode=atomic ${PROJ_ROOT}/...
# -- 2. coverage.html
go tool cover -html "${PROJ_ROOT}/coverage.txt" -o "${PROJ_ROOT}/coverage.html"
