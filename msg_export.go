package mailx

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/mail"
	"strings"
)

// @author valor.

// NewMessage creates a new email message.
func NewMessage() *Message {
	return &Message{
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

// SetTo sets the header of email message: 'TO'.
func (m *Message) SetTo(address ...string) {
	to := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		to = append(to, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.to = to
}

// SetTo adds the header of email message: 'TO'.
func (m *Message) AddTo(address ...string) {
	to := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		to = append(to, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.to = append(m.header.to, to...)
}

// SetRcptTo sets the header of email message: 'TO'.
func (m *Message) SetRcptTo(to ...*mail.Address) {
	m.header.to = to
}

// AddRcptTo adds the header of email message: 'TO'.
func (m *Message) AddRcptTo(to ...*mail.Address) {
	m.header.to = append(m.header.to, to...)
}

// SetCc sets the header of email message: 'CC'.
func (m *Message) SetCc(address ...string) {
	cc := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		cc = append(cc, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.cc = cc
}

// AddCc adds the header of email message: 'CC'.
func (m *Message) AddCc(address ...string) {
	cc := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		cc = append(cc, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.cc = append(m.header.cc, cc...)
}

// SetRcptCc sets the header of email message: 'CC'.
func (m *Message) SetRcptCc(cc ...*mail.Address) {
	m.header.cc = cc
}

// AddRcptCc adds the header of email message: 'CC'.
func (m *Message) AddRcptCc(cc ...*mail.Address) {
	m.header.cc = append(m.header.cc, cc...)
}

// SetBcc sets the header of email message: 'BCC'.
func (m *Message) SetBcc(address ...string) {
	bcc := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		bcc = append(bcc, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.bcc = bcc
}

// AddBcc adds the header of email message: 'BCC'.
func (m *Message) AddBcc(address ...string) {
	bcc := make([]*mail.Address, 0, len(address))
	for _, addr := range address {
		bcc = append(bcc, &mail.Address{
			Name:    "",
			Address: addr,
		})
	}
	m.header.bcc = append(m.header.bcc, bcc...)
}

// SetRcptBcc sets the header of email message: 'BCC'.
func (m *Message) SetRcptBcc(bcc ...*mail.Address) {
	m.header.bcc = bcc
}

// AddRcptBcc adds the header of email message: 'BCC'.
func (m *Message) AddRcptBcc(bcc ...*mail.Address) {
	m.header.bcc = append(m.header.bcc, bcc...)
}

// SetSubject sets the header of email message: 'SUBJECT'.
func (m *Message) SetSubject(subject string) {
	m.header.subject = subject
}

// SetDate sets the header of email message: 'DATE'.
func (m *Message) SetDate(datefmt string) {
	m.header.datefmt = datefmt
}

// SetUserAgent sets the header of email message: 'USER-AGENT'.
func (m *Message) SetUserAgent(ua string) {
	m.header.ua = ua
}

// AddHeader adds other headers of email message.
func (m *Message) AddHeader(key string, value ...string) {
	k := strings.ToUpper(key)
	m.header.extra[k] = value
}

// SetCopierBody sets a custom part of the body of email message.
func (m *Message) SetCopierBody(contentType string, copier CopyFunc) {
	m.parts = []*part{
		{
			ctype:  contentType,
			copier: copier,
		},
	}
}

// AddCopierBody adds a custom part of the body of email message.
func (m *Message) AddCopierBody(contentType string, copier CopyFunc) {
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

// SetPlainBody sets a text part of the body of email message.
func (m *Message) SetPlainBody(text string) {
	m.parts = []*part{
		{
			ctype:  "text/plain",
			copier: newTextCopier(text),
		},
	}
}

// AddPlainBody adds a text part of the body of email message.
func (m *Message) AddPlainBody(text string) {
	m.parts = append(m.parts,
		&part{
			ctype:  "text/plain",
			copier: newTextCopier(text),
		},
	)
}

// SetHtmlBody sets a html part of the body of email message.
func (m *Message) SetHtmlBody(html string) {
	m.parts = []*part{
		{
			ctype:  "text/html",
			copier: newTextCopier(html),
		},
	}
}

// AddHtmlBody adds a html part of the body of email message.
func (m *Message) AddHtmlBody(html string) {
	m.parts = append(m.parts,
		&part{
			ctype:  "text/html",
			copier: newTextCopier(html),
		},
	)
}

// Attach adds a attachment of email message.
func (m *Message) Attach(filename string, copier CopyFunc) {
	f := &file{
		filename:   filename,
		attachment: true,
		copier:     copier,
	}
	m.files = append(m.files, f)
}

// Embed adds a embedded file of email message.
func (m *Message) Embed(cid string, copier CopyFunc) {
	f := &file{
		filename:   cid,
		attachment: false,
		copier:     copier,
	}
	m.files = append(m.files, f)
}

// WriteTo implements io.WriterTo.
// It dumps the whole message to SMTP server.
func (m *Message) WriteTo(w io.Writer) (int64, error) {
	var (
		s int = 0
		n int

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
		n int

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
		n int

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
