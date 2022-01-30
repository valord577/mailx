package mailx

import (
	"net/smtp"
	"testing"
)

func TestLoginAuth(t *testing.T) {

	smtpUser := "user"
	smtpPass := "pass"
	smtpHost := "smtp.example.com"

	auth := &loginAuth{
		username: smtpUser,
		password: smtpPass,
		host:     smtpHost,
	}
	server := &smtp.ServerInfo{
		Name: smtpHost,
		TLS:  false,
		Auth: []string{"LOGIN", "PLAIN"},
	}

	proto, toServer, err := auth.Start(server)
	if err != nil {
		t.Fatalf("loginAuth Start(): %s", err.Error())
	}
	if proto != "LOGIN" {
		t.Fatalf("invalid protocol, got '%s', want 'LOGIN'", proto)
	}
	if toServer != nil {
		t.Fatalf("invalid response, got '%s', want 'nil'", toServer)
	}

	toServer, err = auth.Next([]byte("Username:"), true)
	if err != nil {
		t.Fatalf("loginAuth Next(): %s", err.Error())
	}
	if string(toServer) != smtpUser {
		t.Fatalf("invalid username, got '%s', want '%s'", toServer, smtpUser)
	}

	toServer, err = auth.Next([]byte("Password:"), true)
	if err != nil {
		t.Fatalf("loginAuth Next(): %s", err.Error())
	}
	if string(toServer) != smtpUser {
		t.Fatalf("invalid password, got '%s', want '%s'", toServer, smtpPass)
	}
}
