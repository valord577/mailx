package mailx

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"net/mail"
	"runtime"
	"time"
)

// @author valor.

const (
	charset = "utf-8"

	headerEncoder = mime.BEncoding

	multipartEncoding = "base64"
)

var multipartWriter = func(w io.Writer) io.WriteCloser {
	return base64.NewEncoder(base64.StdEncoding, &multipartBase64Writer{w: w})
}

// CopyFunc is the function that runs when the message is sent.
// It should copy the content of the emails to the io.Writer(SMTP).
type CopyFunc func(io.Writer) (int, error)

// Message represents an email message.
type Message struct {
	header *header
	parts  []*part
	files  []*file
}

func (m *Message) sender() (string, error) {
	if m.header == nil || m.header.from == nil {
		return "", errors.New("empty email sender")
	}
	sender := m.header.from.Address
	if sender == "" {
		return "", errors.New("empty email sender")
	}
	return sender, nil
}

func (m *Message) rcpt() ([]string, error) {
	if m.header == nil {
		return nil, errors.New("empty email rcpt")
	}

	lenTo := len(m.header.to)
	lenCc := len(m.header.cc)
	lenBcc := len(m.header.bcc)
	total := lenTo + lenCc + lenBcc
	if total == 0 {
		return nil, errors.New("empty email rcpt")
	}

	rcpt := make([]string, 0, total)
	if lenTo > 0 {
		for _, address := range m.header.to {
			rcpt = append(rcpt, address.Address)
		}
	}
	if lenCc > 0 {
		for _, address := range m.header.cc {
			rcpt = append(rcpt, address.Address)
		}
	}
	if lenBcc > 0 {
		for _, address := range m.header.bcc {
			rcpt = append(rcpt, address.Address)
		}
	}
	return rcpt, nil
}

type part struct {
	ctype  string // Content-Type
	copier CopyFunc
}

func (p *part) contentType() string {
	return p.ctype + "; charset=" + charset
}

type header struct {
	from *mail.Address
	to   []*mail.Address
	cc   []*mail.Address
	bcc  []*mail.Address

	subject string
	datefmt string

	ua string

	extra map[string][]string
}

func (h *header) presets() []string {
	return []string{"FROM", "TO", "CC", "BCC", "SUBJECT", "DATE", "MIME-VERSION", "USER-AGENT", "MESSAGE-ID"}
}

// date returns a valid RFC 5322 date.
func (h *header) date() string {
	if h.datefmt == "" {
		return time.Now().Format(time.RFC1123Z)
	}
	return h.datefmt
}

// mimeVersion returns MIME-VERSION
func (h *header) mimeVersion() string {
	return "1.0 (Produced by Mailx)"
}

func (h *header) messageId() (string, error) {
	var buf [32]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return "", err
	}
	return "<--" + hex.EncodeToString(buf[:]) + "@GolangMailxMessageID>", nil
}

func (h *header) userAgent() string {
	if h.ua == "" {
		return "github/valord577/mailx " + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
	}
	return h.ua
}

func (h *header) writeTo(w io.Writer) (int, error) {
	if h.from == nil {
		return 0, errors.New("empty email header: 'FROM'")
	}
	if len(h.to) == 0 {
		return 0, errors.New("empty email header: 'TO'")
	}
	if h.subject == "" {
		return 0, errors.New("empty email header: 'SUBJECT'")
	}

	mid, err := h.messageId()
	if err != nil {
		return 0, errors.New("failed to generate 'MESSAGE-ID': " + err.Error())
	}

	b := &bytes.Buffer{}

	// MESSAGE-ID
	b.WriteString("MESSAGE-ID: ")
	b.WriteString(headerEncoder.Encode(charset, mid))
	b.WriteString("\r\n")

	// FROM
	b.WriteString("FROM: ")
	b.WriteString(h.from.String())
	b.WriteString("\r\n")

	// TO
	for _, to := range h.to {
		b.WriteString("TO: ")
		b.WriteString(to.String())
		b.WriteString("\r\n")
	}

	// CC
	length := len(h.cc)
	if length > 0 {
		for _, cc := range h.cc {
			b.WriteString("CC: ")
			b.WriteString(cc.String())
			b.WriteString("\r\n")
		}
	}

	// SUBJECT
	b.WriteString("SUBJECT: ")
	b.WriteString(headerEncoder.Encode(charset, h.subject))
	b.WriteString("\r\n")

	// DATE
	b.WriteString("DATE: ")
	b.WriteString(headerEncoder.Encode(charset, h.date()))
	b.WriteString("\r\n")

	// MIME-VERSION
	b.WriteString("MIME-VERSION: ")
	b.WriteString(headerEncoder.Encode(charset, h.mimeVersion()))
	b.WriteString("\r\n")

	// USER-AGENT
	b.WriteString("USER-AGENT: ")
	b.WriteString(headerEncoder.Encode(charset, h.userAgent()))
	b.WriteString("\r\n")

	// extra headers
	length = len(h.extra)
	if length > 0 {
		for _, key := range h.presets() {
			delete(h.extra, key)
		}
	}

	length = len(h.extra)
	if length > 0 {
		for k, vs := range h.extra {
			for _, v := range vs {
				b.WriteString(k)
				b.WriteString(": ")
				b.WriteString(headerEncoder.Encode(charset, v))
				b.WriteString("\r\n")
			}
		}
	}

	return w.Write(b.Bytes())
}
