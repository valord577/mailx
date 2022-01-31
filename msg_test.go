package mailx

import (
	"io"
	"net/mail"
	"strings"
	"testing"
)

func TestMessage(t *testing.T) {
	m := NewMessage()
	m.SetSender("alex@example.com")
	m.SetTo("bob@example.com", "cora@example.com")
	m.SetRcptCc(&mail.Address{Name: "Dan", Address: "dan@example.com"})
	m.SetSubject("This is a subject of email.")

	m.SetPlainBody(strings.Repeat("This is a text/plain body.", 4))
	m.Attach("attach.txt", func(w io.Writer) (int, error) {
		return io.WriteString(w, "this is a txt attachment.")
	})

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Fatalf("write message, err: %s", err.Error())
	}
}

func TestEmbeded(t *testing.T) {
	m := NewMessage()
	m.SetFrom(&mail.Address{Name: "alex", Address: "alex@example.com"})
	m.AddTo("bob@example.com")
	m.AddCc("cora@example.com")
	m.SetBcc("dan@example.com")
	m.SetSubject("This is a subject of email.")

	m.SetHtmlBody(`This is a text/html body. <img src="cid:CID0"/>`)
	m.Embed("CID0", func(w io.Writer) (int, error) {
		return io.WriteString(w, "this is a embeded attachment.")
	})

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Fatalf("write message, err: %s", err.Error())
	}
}
