package mailx

import (
	"errors"
	"net/smtp"
)

// @author valor.

// loginAuth implements the LOGIN authentication mechanism of the SMTP.
type loginAuth struct {
	username string
	password string
	host     string
}

func isLocalhost(name string) bool {
	return name == "localhost" || name == "127.0.0.1" || name == "::1"
}

// Start implements the stmp.Auth's Start.
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// Must have TLS, or else localhost server.
	// Note: If TLS is not true, then we can't trust ANYTHING in ServerInfo.
	// In particular, it doesn't matter if the server advertises PLAIN auth.
	// That might just be the attacker saying
	// "it's ok, you can trust me with your password."
	if !server.TLS && !isLocalhost(server.Name) {
		return "", nil, errors.New("unencrypted connection")
	}
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	return "LOGIN", nil, nil
}

// Next implements the stmp.Auth's Next.
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	if string(fromServer) == "Username:" {
		return []byte(a.username), nil
	}
	if string(fromServer) == "Password:" {
		return []byte(a.password), nil
	}
	return nil, errors.New("unexpected server challenge: " + string(fromServer))
}
