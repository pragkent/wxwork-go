package oauth2

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pragkent/wxwork-go/oauth2/internal"
)

type Config struct {
	ClientID string

	Endpoint Endpoint

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string

	// Scope specifies optional requested permissions.
	Scopes []string
}

type Endpoint struct {
	AuthURL string
}

// An AuthCodeOption is passed to Config.AuthCodeURL.
type AuthCodeOption interface {
	setValue(url.Values)
}

type setParam struct{ k, v string }

func (p setParam) setValue(m url.Values) { m.Set(p.k, p.v) }

// SetAuthURLParam builds an AuthCodeOption which passes key/value parameters
// to a provider's authorization endpoint.
func SetAuthURLParam(key, value string) AuthCodeOption {
	return setParam{key, value}
}

func (c *Config) AuthCodeURL(state string, opts ...AuthCodeOption) string {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.AuthURL)

	v := url.Values{
		"response_type": {"code"},
		"appid":         {c.ClientID},
	}

	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}

	if len(c.Scopes) > 0 {
		v.Set("scope", strings.Join(c.Scopes, " "))
	}

	if state != "" {
		v.Set("state", state)
	}

	for _, opt := range opts {
		opt.setValue(v)
	}

	if strings.Contains(c.Endpoint.AuthURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}

	buf.WriteString(v.Encode())
	return buf.String()
}

var HTTPClient internal.ContextKey

func NewClient(ctx context.Context, src TokenSource) *http.Client {
	if src == nil {
		return internal.ContextClient(ctx)
	}

	return &http.Client{
		Transport: &Transport{
			Base:   internal.ContextClient(ctx).Transport,
			Source: ReuseTokenSource(nil, src),
		},
	}
}

// RetrieveError is the error returned when the token endpoint returns a
// non-2XX HTTP status code.
type RetrieveError struct {
	Response *http.Response
	// Body is the body that was consumed by reading Response.Body.
	// It may be truncated.
	Body []byte
}

func (r *RetrieveError) Error() string {
	return fmt.Sprintf("oauth2: cannot fetch token: %v: %s", r.Response.Status, r.Body)
}
