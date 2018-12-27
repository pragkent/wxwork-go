package wxwork

import "time"

const expiryDelta = 10 * time.Second

type Token struct {
	AccessToken string    `json:"access_token"`
	Expiry      time.Time `json:"expiry"`
}

// expired reports whether the token is expired.
// t must be non-nil.
func (t *Token) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Round(0).Add(-expiryDelta).Before(time.Now())
}

// Valid reports whether t is non-nil, has an AccessToken, and is not expired.
func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

// TokenSource returns token
type TokenSource interface {
	Token() (*Token, error)
}

type staticTokenSource struct {
	t *Token
}

func (s *staticTokenSource) Token() (*Token, error) {
	return s.t, nil
}

func StaticTokenSource(t *Token) TokenSource {
	return &staticTokenSource{t}
}
