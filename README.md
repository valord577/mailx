Mailx
======

[![Go Report Card](https://goreportcard.com/badge/github.com/valord577/mailx)](https://goreportcard.com/report/github.com/valord577/mailx)
[![Go Reference](https://pkg.go.dev/badge/github.com/valord577/mailx.svg)](https://pkg.go.dev/github.com/valord577/mailx)
[![GitHub License](https://img.shields.io/github/license/valord577/mailx)](LICENSE)

Mailx is a library that makes it easier to send email via SMTP. It is an enhancement of the golang standard library `net/smtp`.

Compatibility
------

- Compatible with go 1.16+

Features
------

Gomail supports:

- Attachments
- Embedded files
- HTML and text templates
- TLS connection and STARTTLS extension
- Sending multiple emails with the same SMTP connection

Installing
------

go mod:

```shell
go get github.com/valord577/mailx
```

Example
------

- See [example](example_test.go)

Changes
------

See the [CHANGES](CHANGE.md) for changes.

License
------

See the [LICENSE](LICENSE) for Rights and Limitations (MIT).
