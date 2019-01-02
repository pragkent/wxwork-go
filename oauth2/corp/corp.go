package corp

import (
	"context"
	"net/http"
	"net/url"

	"github.com/pragkent/wxwork-go/oauth2"
	"github.com/pragkent/wxwork-go/oauth2/internal"
)

const defaultTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"

// Config describes a 2-legged OAuth2 flow, with both the
// client application information and the server's endpoint URLs.
type Config struct {
	CorpID string

	CorpSecret string

	// TokenURL is the resource server's token endpoint
	// URL. This is a constant specific to each server.
	TokenURL string
}

func NewConfig(corpID, corpSecret string) *Config {
	return &Config{
		CorpID:     corpID,
		CorpSecret: corpSecret,
		TokenURL:   defaultTokenURL,
	}
}

// Token uses client credentials to retrieve a token.
//
// The provided context optionally controls which HTTP client is used. See the oauth2.HTTPClient variable.
func (c *Config) Token(ctx context.Context) (*oauth2.Token, error) {
	return c.TokenSource(ctx).Token()
}

// Client returns an HTTP client using the provided token.
// The token will auto-refresh as necessary.
//
// The provided context optionally controls which HTTP client
// is returned. See the oauth2.HTTPClient variable.
//
// The returned Client and its Transport should not be modified.
func (c *Config) Client(ctx context.Context) *http.Client {
	return oauth2.NewClient(ctx, c.TokenSource(ctx))
}

// TokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary using the provided context and the
// client ID and client secret.
//
// Most users will use Config.Client instead.
func (c *Config) TokenSource(ctx context.Context) oauth2.TokenSource {
	source := &tokenSource{
		ctx:  ctx,
		conf: c,
	}

	return oauth2.ReuseTokenSource(nil, source)
}

type tokenSource struct {
	ctx  context.Context
	conf *Config
}

// Token refreshes the token by using a new client credentials request.
func (c *tokenSource) Token() (*oauth2.Token, error) {
	q := url.Values{}
	q.Set("corpid", c.conf.CorpID)
	q.Set("corpsecret", c.conf.CorpSecret)

	tk, err := internal.RetrieveToken(c.ctx, c.conf.TokenURL, q, nil)
	if err != nil {
		if rErr, ok := err.(*internal.RetrieveError); ok {
			return nil, (*oauth2.RetrieveError)(rErr)
		}
		return nil, err
	}

	t := &oauth2.Token{
		AccessToken: tk.AccessToken,
		Expiry:      tk.Expiry,
	}

	return t, nil
}
