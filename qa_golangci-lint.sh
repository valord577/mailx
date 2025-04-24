#!/usr/bin/env bash
set -e

PROJ_ROOT=$(cd "$(dirname ${BASH_SOURCE[0]})"; pwd)

if command -v golangci-lint >/dev/null 2>&1 ; then
  golangci-lint version
  golangci-lint config path
  golangci-lint run ${PROJ_ROOT}/...
else
  go install -ldflags '-s -w' github.com/golangci/golangci-lint/cmd/golangci-lint@latest

  $(go env GOPATH)/bin/golangci-lint version
  $(go env GOPATH)/bin/golangci-lint config path
  $(go env GOPATH)/bin/golangci-lint run ${PROJ_ROOT}/...
fi
