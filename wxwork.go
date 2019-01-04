package wxwork

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "https://qyapi.weixin.qq.com/"
	userAgent      = "wxwork-go/0.1.0"
)

type Client struct {
	client *http.Client

	BaseURL *url.URL

	UserAgent string

	// Reuse a single struct for all the services.
	common service

	Message *MessageService
}

func NewClient(hc *http.Client) *Client {
	if hc == nil {
		hc = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    hc,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}

	c.common.client = c
	c.Message = (*MessageService)(&c.common)

	return c
}

type service struct {
	client *Client
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(url).String()
				return nil, e
			}
		}

		return nil, err
	}
	defer resp.Body.Close()

	body, err := CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, body)
		} else {
			decErr := json.NewDecoder(body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

var sensitiveParams = []string{
	"corpsecret",
	"provider_secret",
	"suite_secret",
	"suite_access_token",
	"access_token",
}

func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}

	params := uri.Query()

	changed := false
	for _, p := range sensitiveParams {
		if len(params.Get(p)) > 0 {
			params.Set(p, "REDACTED")
			changed = true
		}
	}

	if changed {
		uri.RawQuery = params.Encode()
	}

	return uri
}

type ErrorResponse struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`

	HTTPCode int
	Body     string
	Header   http.Header
}

func (e *ErrorResponse) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("wxwork: HTTP response code %d with body: %v", e.HTTPCode, e.Body)
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "wxwork: Error: %d %d: ", e.HTTPCode, e.Code)
	if e.Message != "" {
		fmt.Fprintf(&buf, "%s", e.Message)
	}

	return strings.TrimSpace(buf.String())
}

// CheckResponse returns an error (of type *Error) if the
// status code is not 2xx or errcode is not 0.
func CheckResponse(res *http.Response) (io.Reader, error) {
	var buf bytes.Buffer
	reader := io.TeeReader(res.Body, &buf)

	slurp, err := ioutil.ReadAll(reader)
	if err == nil {
		jerr := new(ErrorResponse)
		err = json.Unmarshal(slurp, jerr)
		if err == nil && jerr.Code == 0 {
			return &buf, nil
		}

		jerr.HTTPCode = res.StatusCode
		jerr.Body = string(slurp)
		return nil, jerr
	}

	return nil, &ErrorResponse{
		HTTPCode: res.StatusCode,
		Body:     string(slurp),
		Header:   res.Header,
	}
}
