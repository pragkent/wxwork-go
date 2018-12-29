package oauth2

import (
	"testing"
	"time"
)

func TestTokenValidNoAccessToken(t *testing.T) {
	token := &Token{}
	if token.Valid() {
		t.Errorf("got valid with no access token; want invalid")
	}
}

func TestExpiredWithExpiry(t *testing.T) {
	token := &Token{
		AccessToken: "abc",
		Expiry:      time.Now().Add(-5 * time.Minute),
	}

	if token.Valid() {
		t.Errorf("got valid with expired token; want invalid")
	}
}
