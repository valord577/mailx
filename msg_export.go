package mailx

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net/mail"
	"strings"
)

// @author valor.

func NewMessage() *message {
	return &message{
		header: &header{
			to:  make([]*mail.Address, 0),
			cc:  make([]*mail.Address, 0),
			bcc: make([]*mail.Address, 0),

			extra: make(map[string][]string),
		},
		parts: make([]*part, 0),
		files: make([]*file, 0),
	}
}

func (m *message) SetSender(address string) {
	m.header.from = &mail.Address{
		Name:    "",
		Address: address,
	}
}

func (m *message) SetFrom(sender *mail.Address) {
	m.header.from = sender
}

func (m *message) SetTo(address ...string) {
	to := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		to = append(to, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.to = to
}

func (m *message) AddTo(address ...string) {
	to := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		to = append(to, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.to = append(m.header.to, to...)
}

func (m *message) SetRcptTo(to ...*mail.Address) {
	m.header.to = to
}

func (m *message) AddRcptTo(to ...*mail.Address) {
	m.header.to = append(m.header.to, to...)
}

func (m *message) SetCc(address ...string) {
	cc := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		cc = append(cc, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.cc = cc
}

func (m *message) AddCc(address ...string) {
	cc := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		cc = append(cc, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.cc = append(m.header.cc, cc...)
}

func (m *message) SetRcptCc(cc ...*mail.Address) {
	m.header.cc = cc
}

func (m *message) AddRcptCc(cc ...*mail.Address) {
	m.header.cc = append(m.header.cc, cc...)
}

func (m *message) SetBcc(address ...string) {
	bcc := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		bcc = append(bcc, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.bcc = bcc
}

func (m *message) AddBcc(address ...string) {
	bcc := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		bcc = append(bcc, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.bcc = append(m.header.bcc, bcc...)
}

func (m *message) SetRcptBcc(bcc ...*mail.Address) {
	m.header.bcc = bcc
}

func (m *message) AddRcptBcc(bcc ...*mail.Address) {
	m.header.bcc = append(m.header.bcc, bcc...)
}

func (m *message) SetSubject(subject string) {
	m.header.subject = subject
}

func (m *message) SetDate(datefmt string) {
	m.header.datefmt = datefmt
}

func (m *message) SetUserAgent(ua string) {
	m.header.ua = ua
}

func (m *message) AddHeader(key string, value ...string) {
	k := strings.ToUpper(key)
	m.header.extra[k] = value
}

func (m *message) SetCopierBody(contentType string, copier CopyFunc) {
	m.parts = []*part{
		{
			ctype:  contentType,
			copier: copier,
		},
	}
}

func (m *message) AddCopierBody(contentType string, copier CopyFunc) {
	m.parts = append(m.parts,
		&part{
			ctype:  contentType,
			copier: copier,
		},
	)
}

func newTextCopier(s string) CopyFunc {
	return func(w io.Writer) (int, error) {
		return io.WriteString(w, s)
	}
}

func (m *message) SetPlainBody(text string) {
	m.parts = []*part{
		{
			ctype:  "text/plain",
			copier: newTextCopier(text),
		},
	}
}

func (m *message) AddPlainBody(text string) {
	m.parts = append(m.parts,
		&part{
			ctype:  "text/plain",
			copier: newTextCopier(text),
		},
	)
}

func (m *message) SetHtmlBody(html string) {
	m.parts = []*part{
		{
			ctype:  "text/html",
			copier: newTextCopier(html),
		},
	}
}

func (m *message) AddHtmlBody(html string) {
	m.parts = append(m.parts,
		&part{
			ctype:  "text/html",
			copier: newTextCopier(html),
		},
	)
}

func (m *message) Attach(filename string, saved string, copier CopyFunc) {
	f := &file{
		filename:   filename,
		saved:      saved,
		attachment: true,
		copier:     copier,
	}
	m.files = append(m.files, f)
}

func (m *message) Embed(filename string, copier CopyFunc) {
	f := &file{
		filename:   filename,
		attachment: false,
		copier:     copier,
	}
	m.files = append(m.files, f)
}

func (m *message) WriteTo(w io.Writer) (int64, error) {
	if m.header == nil {
		return 0, errors.New("empty email header")
	}

	var (
		s int = 0
		n int = 0

		err error
	)

	var buf [30]byte
	_, err = rand.Read(buf[:])
	if err != nil {
		return 0, err
	}
	boundary := "--GolangMailxBoundary" + hex.EncodeToString(buf[:])

	partStart := "--" + boundary
	partClose := "--" + boundary + "--"

	n, err = m.header.writeTo(w)
	if err != nil {
		return 0, err
	}
	s += n

	n, err = io.WriteString(w, "Content-Type: multipart/mixed;\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	n, err = io.WriteString(w, " boundary="+boundary+"\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	n, err = io.WriteString(w, "\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	if len(m.parts) > 0 {
		for _, part := range m.parts {
			n, err = writePart(partStart, part, w)
			if err != nil {
				return 0, err
			}
			s += n
		}
	}

	if len(m.files) > 0 {
		for _, file := range m.files {
			n, err = writeFile(partStart, file, w)
			if err != nil {
				return 0, err
			}
			s += n
		}
	}

	n, err = io.WriteString(w, partClose+"\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	return int64(s), nil
}

func writePart(partStart string, part *part, out io.Writer) (int, error) {
	var (
		s int = 0
		n int = 0

		err error
	)

	n, err = io.WriteString(out, partStart+"\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	n, err = io.WriteString(out, "Content-Type: "+part.contentType()+"\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	n, err = io.WriteString(out, "Content-Transfer-Encoding: "+multipartEncoding+"\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	n, err = io.WriteString(out, "\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	// Headers ended, write the body of part
	partWriter := multipartWriter(out)
	n, err = part.copier(partWriter)
	if err != nil {
		return 0, err
	}
	partWriter.Close()
	s += n

	n, err = io.WriteString(out, "\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	return s, nil
}

func writeFile(partStart string, file *file, out io.Writer) (int, error) {
	var (
		s int = 0
		n int = 0

		err error
	)

	n, err = io.WriteString(out, partStart+"\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	// The name of file may contain non-ascii characters.
	n, err = io.WriteString(out, "Content-Type: "+headerEncoder.Encode(charset, file.contentType())+"\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	// The name of file may contain non-ascii characters.
	n, err = io.WriteString(out, "Content-Disposition: "+headerEncoder.Encode(charset, file.disposition())+"\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	if !file.attachment {
		// The name of file may contain non-ascii characters.
		n, err = io.WriteString(out, "Content-ID: <"+headerEncoder.Encode(charset, file.filename)+">\r\n")
		if err != nil {
			return 0, err
		}
		s += n
	}

	n, err = io.WriteString(out, "Content-Transfer-Encoding: "+multipartEncoding+"\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	n, err = io.WriteString(out, "\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	// Headers ended, write the body of file
	partWriter := multipartWriter(out)
	n, err = file.copier(partWriter)
	if err != nil {
		return 0, err
	}
	partWriter.Close()
	s += n

	n, err = io.WriteString(out, "\r\n")
	if err != nil {
		return 0, err
	}
	s += n

	return s, nil
}
