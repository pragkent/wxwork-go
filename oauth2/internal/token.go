package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Token represents the credentials used to authorize
// the requests to access protected resources on the OAuth 2.0
// provider's backend.
//
// This type is a mirror of oauth2.Token and exists to break
// an otherwise-circular dependency. Other internal packages
// should convert this Token into an oauth2.Token before use.
type Token struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time
}

// tokenJSON is the struct representing the HTTP response from OAuth2
// providers returning a token in JSON form.
type tokenJSON struct {
	ErrCode     int            `json:"errcode"`
	ErrMsg      string         `json:"errmsg"`
	AccessToken string         `json:"access_token"`
	ExpiresIn   expirationTime `json:"expires_in"`
}

func (e *tokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}

	return
}

type expirationTime int32

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	*e = expirationTime(i)
	return nil
}

func RetrieveToken(ctx context.Context, tokenURL string, q url.Values, v interface{}) (*Token, error) {
	req, err := newRequest(ctx, tokenURL, q, v)
	if err != nil {
		return nil, err
	}

	req.WithContext(ctx)

	r, err := ContextClient(ctx).Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}

	if code := r.StatusCode; code < 200 || code > 299 {
		return nil, &RetrieveError{
			Response: r,
			Body:     body,
		}
	}

	var tj tokenJSON
	if err = json.Unmarshal(body, &tj); err != nil {
		return nil, err
	}

	if tj.ErrCode != 0 {
		return nil, &RetrieveError{
			Response: r,
			Body:     body,
		}
	}

	token := &Token{
		AccessToken: tj.AccessToken,
		Expiry:      tj.expiry(),
	}

	if token.AccessToken == "" {
		return token, errors.New("oauth2: server response missing access_token")
	}

	return token, nil
}

func newRequest(ctx context.Context, tokenURL string, q url.Values, v interface{}) (*http.Request, error) {
	if len(q) != 0 {
		tokenURL += "?"
		tokenURL += q.Encode()
	}

	var body io.Reader
	if v != nil {
		out, err := json.Marshal(v)
		if err != nil {
			return nil, errors.New("oauth2: marshal request body error")
		}
		body = bytes.NewReader(out)
	}

	return http.NewRequest("POST", tokenURL, body)
}

type RetrieveError struct {
	Response *http.Response
	Body     []byte
}

func (r *RetrieveError) Error() string {
	return fmt.Sprintf("oauth2: cannot fetch token: %v: %s", r.Response.Status, r.Body)
}
