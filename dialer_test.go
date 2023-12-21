package mailx

import (
	"crypto/tls"
	"io"
	"net"
	"net/smtp"
	"testing"
)

func TestSmtpTlsLoginAuth(t *testing.T) {
	m := map[string]string{
		"AUTH": "LOGIN",
	}
	testSmtp(t, true, m)
}

func TestSmtpTlsPlainAuth(t *testing.T) {
	m := map[string]string{
		"AUTH": "PLAIN",
	}
	testSmtp(t, true, m)
}

func TestSmtpTlsCRAMMD5Auth(t *testing.T) {
	m := map[string]string{
		"AUTH": "CRAM-MD5",
	}
	testSmtp(t, true, m)
}

func TestSmtpStarttlsLoginAuth(t *testing.T) {
	m := map[string]string{
		"AUTH":     "LOGIN",
		"STARTTLS": "",
	}
	testSmtp(t, true, m)
}

func TestSmtpStarttlsPlainAuth(t *testing.T) {
	m := map[string]string{
		"AUTH":     "PLAIN",
		"STARTTLS": "",
	}
	testSmtp(t, false, m)
}

func TestSmtpStarttlsCRAMMD5Auth(t *testing.T) {
	m := map[string]string{
		"AUTH":     "CRAM-MD5",
		"STARTTLS": "",
	}
	testSmtp(t, false, m)
}

func testSmtp(t *testing.T, ssl bool, ext map[string]string) {

	smtpUser := "user"
	smtpPass := "pass"
	smtpHost := "smtp.example.com"
	smtpPort := 777

	d := &Dialer{
		Host: smtpHost,
		Port: smtpPort,

		Username: smtpUser,
		Password: smtpPass,

		SSLOnConnect: ssl,
	}
	if ssl {
		d.TLSConfig = &tls.Config{ServerName: d.Host}
	}

	netDial = func(*net.Dialer, string, string) (net.Conn, error) {
		return nil, nil
	}
	tlsDial = func(*tls.Dialer, string, string) (net.Conn, error) {
		return nil, nil
	}
	newSmtpClient = func(net.Conn, string) (smtpClient, error) {
		return &mockSmtpClient{ext}, nil
	}

	m0 := NewMessage()
	m0.SetTo("aaaaa@example.com")
	m0.SetBcc("ccccc@example.com")
	m0.SetSubject("This is a subject of email.")
	m0.SetPlainBody("This is a text/plain body.")

	err := d.DialAndSend(m0)
	if err != nil {
		t.Fatalf("dialer err: %s", err.Error())
	}

	m1 := NewMessage()
	m1.SetSender("")
	m1.SetTo("aaaaa@example.com")
	m1.SetCc("bbbbb@example.com")
	m1.SetSubject("This is a subject of email.")
	m1.SetPlainBody("This is a text/plain body.")

	err = d.DialAndSend(m1)
	if err != nil {
		t.Fatalf("dialer err: %s", err.Error())
	}
}

type mockSmtpClient struct {
	ext map[string]string
}

func (c *mockSmtpClient) Hello(localName string) error {
	return nil
}

func (c *mockSmtpClient) Extension(ext string) (bool, string) {
	value, ok := c.ext[ext]
	return ok, value
}

func (c *mockSmtpClient) StartTLS(config *tls.Config) error {
	return nil
}

func (c *mockSmtpClient) Auth(a smtp.Auth) error {
	return nil
}

func (c *mockSmtpClient) Mail(from string) error {
	return nil
}

func (c *mockSmtpClient) Rcpt(to string) error {
	return nil
}

func (c *mockSmtpClient) Data() (io.WriteCloser, error) {
	return &mockWriter{}, nil
}

func (c *mockSmtpClient) Quit() error {
	return nil
}

func (c *mockSmtpClient) Close() error {
	return nil
}

type mockWriter struct{}

func (*mockWriter) Write(p []byte) (int, error) {
	return io.Discard.Write(p)
}

func (*mockWriter) Close() error {
	return nil
}
