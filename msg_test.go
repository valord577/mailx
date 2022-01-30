package mailx

import (
	"io"
	"net/mail"
	"os"
	"testing"
)

func TestMessage(t *testing.T) {
	m := NewMessage()
	m.SetSender("alex@example.com")
	m.SetTo("bob@example.com", "cora@example.com")
	m.SetRcptCc(&mail.Address{Name: "Dan", Address: "dan@example.com"})
	m.SetSubject("This is a subject of email.")

	m.SetPlainBody("This is a text/plain body.")
	m.Attach("attach.txt", func(w io.Writer) (int, error) {
		return io.WriteString(w, "this is a txt attachment.")
	})

	_, err := m.WriteTo(os.Stdout)
	if err != nil {
		t.Fatalf("write message, err: %s", err.Error())
	}
}
