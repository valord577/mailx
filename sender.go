package mailx

import (
	"io"
)

// @author valor.

// Sender sends emails via *smtp.Client
type Sender struct {
	smtpClient
	from string
}

// Send sends the given emails.
func (s *Sender) Send(m *Message) error {
	rcpt, err := m.rcpt()
	if err != nil {
		return err
	}

	m.setFrom(s.from)
	return s.send(s.from, rcpt, m)
}

// SendOne sends a message implements io.WriterTo
func (s *Sender) send(from string, to []string, msg io.WriterTo) error {
	if err := s.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err := s.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := s.Data()
	if err != nil {
		return err
	}

	if _, err = msg.WriteTo(w); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

// Close sends the QUIT command and closes the connection to the server.
func (s *Sender) Close() error {
	return s.Quit()
}
