# Changes

All notable changes to **mailx** are documented in this file.

## v0.5.20231221

#### Added

- Can manually set email message header field `FROM`.
    * `func (m *Message) SetSender(address string)`
    * `func (m *Message) SetFrom(sender *mail.Address)`

## v0.4.20231124

#### Fixed

- Unable to display unicode attachment name correctly.

## v0.3.20220208

#### Changed

- Adding unit tests for coverage.
- Hide manually set email sender address.

## v0.2.20211115

#### Changed

- Export a struct: `Message`.

## v0.1.20211112

#### Added

- Release first version.
