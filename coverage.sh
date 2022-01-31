#!/usr/bin/env bash

COVERAGE_TXT="coverage.txt"
COVERAGE_HTML="coverage.html"
GO_MODULE="github.com/valord577/mailx"

TEST_SMTP_HOST_K="testSmtpHost"
TEST_SMTP_HOST_V=""
TEST_SMTP_PORT_K="testSmtpPort"
TEST_SMTP_PORT_V=""
TEST_SMTP_USER_K="testSmtpUser"
TEST_SMTP_USER_V=""
TEST_SMTP_PASS_K="testSmtpPass"
TEST_SMTP_PASS_V=""
TEST_SMTP_SSL_K="testSmtpSSL"
TEST_SMTP_SSL_V=""
TEST_MAIL_RECV_TO_K="testMailRecvTo"
TEST_MAIL_RECV_TO_V=""

# -- 1. coverage.txt
go test \
  -ldflags "\
    -X '${GO_MODULE}.${TEST_SMTP_HOST_K}=${TEST_SMTP_HOST_V}' \
    -X '${GO_MODULE}.${TEST_SMTP_PORT_K}=${TEST_SMTP_PORT_V}' \
    -X '${GO_MODULE}.${TEST_SMTP_USER_K}=${TEST_SMTP_USER_V}' \
    -X '${GO_MODULE}.${TEST_SMTP_PASS_K}=${TEST_SMTP_PASS_V}' \
    -X '${GO_MODULE}.${TEST_SMTP_SSL_K}=${TEST_SMTP_SSL_V}' \
    -X '${GO_MODULE}.${TEST_MAIL_RECV_TO_K}=${TEST_MAIL_RECV_TO_V}' \
  " \
  "${GO_MODULE}" \
  -race -coverprofile=${COVERAGE_TXT} -covermode=atomic
# -- 2. coverage.html
go tool cover -html ${COVERAGE_TXT}  -o ${COVERAGE_HTML}
