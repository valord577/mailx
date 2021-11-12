package mailx

import (
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// @author valor.

// Dialer is a dialer to an SMTP server.
type Dialer struct {
	// Host represents the host of the SMTP server.
	Host string
	// Port represents the port of the SMTP server.
	Port int
	// Username is the username to use to authenticate to the SMTP server.
	Username string
	// Password is the password to use to authenticate to the SMTP server.
	Password string
	// SSLOnConnect defines whether an SSL connection is used.
	// It should be false while SMTP server use the STARTTLS extension.
	SSLOnConnect bool
	// TSLConfig represents the TLS configuration.
	// It is used for the TLS (when the
	// STARTTLS extension is used) or SSL connection.
	TLSConfig *tls.Config
	// Timeout is passed to net.Dialer's Timeout.
	Timeout time.Duration
}

func (d *Dialer) addr() string {
	return d.Host + ":" + strconv.FormatInt(int64(d.Port), 10)
}

func (d *Dialer) tlsConfig() *tls.Config {
	if d.TLSConfig == nil {
		return &tls.Config{ServerName: d.Host}
	}
	return d.TLSConfig
}

func (d *Dialer) smtpClient(conn net.Conn, host string) (*smtp.Client, error) {
	return smtp.NewClient(conn, host)
}

func (d *Dialer) smtpAuth(c *smtp.Client) (smtp.Auth, error) {
	if d.Username == "" {
		return nil, nil
	}

	ok, auths := c.Extension("AUTH")
	if !ok {
		return nil, errors.New("smtp server doesn't support AUTH")
	}

	if strings.Contains(auths, "CRAM-MD5") {
		return smtp.CRAMMD5Auth(d.Username, d.Password), nil
	}
	if strings.Contains(auths, "PLAIN") {
		return smtp.PlainAuth("", d.Username, d.Password, d.Host), nil
	}
	if strings.Contains(auths, "LOGIN") {
		return &loginAuth{
			username: d.Username,
			password: d.Password,
			host:     d.Host,
		}, nil
	}
	return nil, errors.New("no authentication mechanism is implemented: " + auths)
}

// Dial dials and authenticates to an SMTP server.
// The returned *Sender should be closed when done using it.
func (d *Dialer) Dial() (*Sender, error) {
	if d.SSLOnConnect {
		return d.dialSSL()
	}
	// debug: openssl s_client -starttls smtp -ign_eof -crlf -connect <host>:<port>
	return d.dialStartTLS()
}

func (d *Dialer) dialSSL() (*Sender, error) {
	netDialer := &net.Dialer{Timeout: d.Timeout}
	conn, err := tls.DialWithDialer(netDialer, "tcp", d.addr(), d.tlsConfig())
	if err != nil {
		return nil, err
	}

	c, err := d.smtpClient(conn, d.Host)
	if err != nil {
		return nil, err
	}

	auth, err := d.smtpAuth(c)
	if err != nil {
		c.Close()
		return nil, err
	}

	if auth != nil {
		if err = c.Auth(auth); err != nil {
			c.Close()
			return nil, err
		}
	}
	return &Sender{c}, nil
}

func (d *Dialer) dialStartTLS() (*Sender, error) {
	netDialer := &net.Dialer{Timeout: d.Timeout}
	conn, err := netDialer.Dial("tcp", d.addr())
	if err != nil {
		return nil, err
	}

	c, err := d.smtpClient(conn, d.Host)
	if err != nil {
		return nil, err
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(d.tlsConfig()); err != nil {
			c.Close()
			return nil, err
		}
	}

	auth, err := d.smtpAuth(c)
	if err != nil {
		c.Close()
		return nil, err
	}

	if auth != nil {
		if err = c.Auth(auth); err != nil {
			c.Close()
			return nil, err
		}
	}
	return &Sender{c}, nil
}

// DialAndSend opens a connection to the SMTP server,
// sends the given emails and closes the connection.
func (d *Dialer) DialAndSend(m *message) error {
	s, err := d.Dial()
	if err != nil {
		return err
	}
	defer s.Close()

	return s.Send(m)
}
