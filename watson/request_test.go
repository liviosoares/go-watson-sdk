package watson

import (
	"testing"
)

func TestMissingUrl(t *testing.T) {
	creds := Credentials{
		Url:      "",
		Username: "foo",
		Password: "foo",
	}
	_, err := NewClient(creds)
	if err == nil {
		t.Error("Expected client creation to fail")
	}
}

func TestMissingUsername(t *testing.T) {
	creds := Credentials{
		Url:      "foo",
		Username: "",
		Password: "foo",
	}
	_, err := NewClient(creds)
	if err == nil {
		t.Error("Expected client creation to fail")
	}
}

func TestMissingPassword(t *testing.T) {
	creds := Credentials{
		Url:      "foo",
		Username: "foo",
		Password: "",
	}
	_, err := NewClient(creds)
	if err == nil {
		t.Error("Expected client creation to fail")
	}
}

func TestAllOk(t *testing.T) {
	creds := Credentials{
		Url:      "foo",
		Username: "foo",
		Password: "foo",
	}
	client, err := NewClient(creds)
	if err != nil {
		t.Error("Expected client creation to succeed")
	}

	if client.Creds != creds {
		t.Error("Credentials do not match")
	}
}

func TestApiKeySufficient(t *testing.T) {
	creds := Credentials{
		ApiKey: "foo",
	}
	client, err := NewClient(creds)
	if err != nil {
		t.Error("Expected client creation to succeed")
	}

	if client.Creds != creds {
		t.Error("Credentials do not match")
	}
}
