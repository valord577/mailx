package mailx

import (
	"io"
	"strconv"
	"testing"
)

var (
	testSmtpHost = "smtp.example.com"
	testSmtpPort = "465"
	testSmtpUser = "sender@example.com"
	testSmtpPass = "pass"
	testSmtpSSL  = "true"

	testMailRecvTo = "recv-to@example.com"
)

func TestDialer(t *testing.T) {
	smtpPort, err := strconv.ParseInt(testSmtpPort, 10, 32)
	if err != nil {
		t.Fatalf("unkonwn smtp port: %s", testSmtpPort)
	}
	sslOnConnect, err := strconv.ParseBool(testSmtpSSL)
	if err != nil {
		t.Fatalf("unkonwn smtp port: %s", testSmtpPort)
	}

	m := NewMessage()
	m.SetSender(testSmtpUser)
	m.SetTo(testMailRecvTo)
	m.SetSubject("This is a subject of email.")

	m.SetPlainBody("This is a text/plain body.")
	m.Attach("attach.txt", func(w io.Writer) (int, error) {
		return io.WriteString(w, "this is a txt attachment.")
	})

	d := &Dialer{
		Host: testSmtpHost,
		Port: int(smtpPort),

		Username: testSmtpUser,
		Password: testSmtpPass,

		SSLOnConnect: sslOnConnect,
	}

	err = d.DialAndSend(m)
	if err != nil {
		t.Fatalf("dialer err: %s", err.Error())
	}
}
