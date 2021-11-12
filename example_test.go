package mailx

import (
	"fmt"
	"io"
	"net/mail"
	"testing"
	"time"
)

// @author valor.

func TestSample(t *testing.T) {
	const (
		smtpHost = "smtp.example.com"
		smtpPort = 465
		username = "user"
		password = "123456"

		sslOnConnect = true
	)

	m := NewMessage()
	m.SetSender("alex@example.com")
	m.SetTo("bob@example.com", "cora@example.com")
	m.SetRcptCc(&mail.Address{Name: "Dan", Address: "dan@example.com"})
	m.SetSubject("This is a subject of email.")

	m.SetPlainBody("This is a text/plain body.")
	m.Attach("attach.txt", "download.txt", func(w io.Writer) (int, error) {
		return io.WriteString(w, "this is a txt attachment.")
	})

	d := &Dialer{
		Host: smtpHost,
		Port: smtpPort,

		Username: username,
		Password: password,

		SSLOnConnect: sslOnConnect,
	}
	err := d.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}

func TestSampleDaemon(t *testing.T) {
	const (
		smtpHost = "smtp.example.com"
		smtpPort = 465
		username = "user"
		password = "123456"

		sslOnConnect = true
	)

	// Use the channel in your program to send emails.
	ch := make(chan *message)

	go func() {
		d := &Dialer{
			Host: smtpHost,
			Port: smtpPort,

			Username: username,
			Password: password,

			SSLOnConnect: sslOnConnect,
		}

		var ser *Sender
		var err error
		open := false
		for {
			select {
			case m := <-ch:
				if !open {
					if ser, err = d.Dial(); err != nil {
						fmt.Printf("%s\n", err.Error())
					}
					open = true
				}
				if err := ser.Send(m); err != nil {
					fmt.Printf("%s\n", err.Error())
				}
			// Close the connection to the SMTP server
			// if no email was sent in the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					if err := ser.Close(); err != nil {
						fmt.Printf("%s\n", err.Error())
					}
					open = false
				}
			}
		}
	}()

	// Close the channel to stop the mail daemon.
	close(ch)
}
