package users

import (
	"testing"
	"time"

	"github.com/nu7hatch/gouuid"
)

func TestGenerateAPIKey(t *testing.T) {
	u := User{}

	if err := u.GenAPIKey(); err != nil {
		t.Error("Unexpected error:", err)
	}

	if u.APIKey == new(uuid.UUID).String() {
		t.Error("Unexpected zero UUID")
	}

	if u.APIKeyExp.IsZero() {
		t.Error("Expected APIKeyExp to not be zero")
	}

	apiKey := u.APIKey
	exp := u.APIKeyExp

	if err := u.GenAPIKey(); err != nil {
		t.Error("Unexpected error:", err)
	}

	if u.APIKey == apiKey {
		t.Error("Expected api key to change")
	}

	if u.APIKeyExp == exp {
		t.Error("Expected expiration to change")
	}
}

func TestValidateAPIKey(t *testing.T) {
	u := User{}

	if err := u.ValidateAPIKey(); err != ErrNoAPIKey {
		t.Error("Expected ErrNoAPIKey, got:", err)
	}

	if err := u.GenAPIKey(); err != nil {
		t.Error("Unexpected error:", err)
	}

	if err := u.ValidateAPIKey(); err != nil {
		t.Error("Expected no error, got:", err)
	}

	u.APIKeyExp = time.Time{}

	if err := u.ValidateAPIKey(); err != ErrAPIKeyExpired {
		t.Error("Expected ErrAPIKeyExpired, got:", err)
	}
}

func TestDestroyAPIKey(t *testing.T) {
	u := User{}
	if err := u.GenAPIKey(); err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	apiKey := u.APIKey

	u.DestroyAPIKey()

	if apiKey == u.APIKey {
		t.Error("Expected APIKey to have changed")
	}

	if u.APIKey != new(uuid.UUID).String() {
		t.Error("Expected APIKey to be zero")
	}

	if !u.APIKeyExp.IsZero() {
		t.Error("Expected expiration to be zero:", u.APIKeyExp)
	}
}
