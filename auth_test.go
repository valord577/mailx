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
		TLS:  true,
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
	if string(toServer) != smtpPass {
		t.Fatalf("invalid password, got '%s', want '%s'", toServer, smtpPass)
	}
}

func TestLoginAuthStartErr(t *testing.T) {
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
	if err == nil || proto != "" || toServer != nil {
		t.Fatalf("invalid response")
	}

	server = &smtp.ServerInfo{
		Name: "localhost",
		TLS:  true,
		Auth: []string{"LOGIN", "PLAIN"},
	}
	proto, toServer, err = auth.Start(server)
	if err == nil || proto != "" || toServer != nil {
		t.Fatalf("invalid response")
	}

	server = &smtp.ServerInfo{
		Name: "abcd.example.com",
		TLS:  true,
		Auth: []string{"LOGIN", "PLAIN"},
	}
	proto, toServer, err = auth.Start(server)
	if err == nil || proto != "" || toServer != nil {
		t.Fatalf("invalid response")
	}
}

func TestLoginAuthNextErr(t *testing.T) {
	smtpUser := "user"
	smtpPass := "pass"
	smtpHost := "smtp.example.com"

	auth := &loginAuth{
		username: smtpUser,
		password: smtpPass,
		host:     smtpHost,
	}

	toServer, err := auth.Next([]byte("everything"), false)
	if err != nil || toServer != nil {
		t.Fatalf("got err: %s", err.Error())
	}

	toServer, err = auth.Next([]byte("everything"), true)
	if err == nil || toServer != nil {
		t.Fatalf("invalid response")
	}
}
