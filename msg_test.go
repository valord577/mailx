package mailx

import (
	"io"
	"net/mail"
	"strings"
	"testing"
	"time"
)

func TestMessage1(t *testing.T) {
	m := NewMessage()
	m.SetSender("alex@example.com")
	m.SetTo("aaa-1@example.com")
	m.AddTo("aaa-2@example.com")
	m.SetCc("bbb-1@example.com")
	m.AddCc("bbb-2@example.com")
	m.SetBcc("ccc-1@example.com")
	m.AddBcc("ccc-2@example.com")

	m.SetSubject("This is a subject of email.")

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Fatalf("write message, err: %s", err.Error())
	}
}

func TestMessage2(t *testing.T) {
	m := NewMessage()
	m.SetFrom(&mail.Address{Name: "alex", Address: "alex@example.com"})
	m.SetRcptTo(&mail.Address{Name: "aaa-1", Address: "aaa-1@example.com"})
	m.AddRcptTo(&mail.Address{Name: "aaa-2", Address: "aaa-2@example.com"})
	m.SetRcptCc(&mail.Address{Name: "bbb-1", Address: "bbb-1@example.com"})
	m.AddRcptCc(&mail.Address{Name: "bbb-2", Address: "bbb-2@example.com"})
	m.SetRcptBcc(&mail.Address{Name: "ccc-1", Address: "ccc-1@example.com"})
	m.AddRcptBcc(&mail.Address{Name: "ccc-2", Address: "ccc-2@example.com"})

	m.SetSubject("This is a subject of email.")

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Fatalf("write message, err: %s", err.Error())
	}
}

func TestMessage3(t *testing.T) {
	m := NewMessage()
	m.SetSender("alex@example.com")
	m.SetTo("aaaaa@example.com")

	m.SetSubject("This is a subject of email.")
	m.SetPlainBody("This is a text/plain body.This is a text/plain body.")
	m.AddPlainBody("This is a text/plain body.This is a text/plain body.")

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Fatalf("write message, err: %s", err.Error())
	}
}

func TestMessage4(t *testing.T) {
	m := NewMessage()
	m.SetSender("alex@example.com")
	m.SetTo("aaaaa@example.com")

	m.SetSubject("This is a subject of email.")
	m.SetHtmlBody("<p>This is a text/html body.This is a text/html body.</p>")
	m.AddHtmlBody("<p>This is a text/html body.This is a text/html body.</p>")

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Fatalf("write message, err: %s", err.Error())
	}
}

func TestMessage5(t *testing.T) {
	m := NewMessage()
	m.SetSender("alex@example.com")
	m.SetTo("aaaaa@example.com")

	m.SetSubject("This is a subject of email.")
	m.SetCopierBody("text/plain", func(w io.Writer) (int, error) {
		body := strings.Repeat("This is a text/plain body.", 4)
		return io.WriteString(w, body)
	})
	m.AddCopierBody("text/html", func(w io.Writer) (int, error) {
		body := strings.Repeat("<p>This is a text/html body.</p>", 4)
		return io.WriteString(w, body)
	})

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Fatalf("write message, err: %s", err.Error())
	}
}

func TestMessage6(t *testing.T) {
	m := NewMessage()
	m.SetSender("alex@example.com")
	m.SetTo("aaaaa@example.com")

	m.SetSubject("This is a subject of email.")
	m.SetHtmlBody(`This is a text/html body. <img src="cid:CID0"/>`)
	m.Embed("CID0", func(w io.Writer) (int, error) {
		return io.WriteString(w, "this is a embedded attachment.")
	})
	m.Attach("attach.txt", func(w io.Writer) (int, error) {
		return io.WriteString(w, "this is a txt attachment.")
	})

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Fatalf("write message, err: %s", err.Error())
	}
}

func TestMessage7(t *testing.T) {
	m := NewMessage()
	m.SetSender("alex@example.com")
	m.SetTo("aaaaa@example.com")

	m.SetUserAgent("ua - test mailx")
	m.SetDate(time.Now().Format(time.RFC1123Z))
	m.AddHeader("extra", "EXTRA HEADER")

	m.SetSubject("This is a subject of email.")

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Fatalf("write message, err: %s", err.Error())
	}
}

func TestErrMessage1(t *testing.T) {
	m := NewMessage()
	m.SetSender("alex@example.com")

	m.SetSubject("This is a subject of email.")

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Logf("write message, err: %s", err.Error())
	}
}

func TestErrMessage2(t *testing.T) {
	m := NewMessage()
	m.SetSender("alex@example.com")
	m.SetTo("aaaaa@example.com")

	_, err := m.WriteTo(io.Discard)
	if err != nil {
		t.Logf("write message, err: %s", err.Error())
	}
}
