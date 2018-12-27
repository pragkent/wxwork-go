package wxwork

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func NewClient() *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
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

type Response struct {
	*http.Response
}

type Error struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`

	HTTPCode int
	Body     string
	Header   http.Header
}

func (e *Error) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("wxwork: got HTTP response code %d with body: %v", e.HTTPCode, e.Body)
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
func CheckResponse(res *Response) error {
	slurp, err := ioutil.ReadAll(res.Body)
	if err == nil {
		jerr := new(Error)
		err = json.Unmarshal(slurp, jerr)
		if err == nil && jerr.Code == 0 {
			return nil
		}

		jerr.HTTPCode = res.StatusCode
		jerr.Body = string(slurp)
		return jerr
	}

	return &Error{
		HTTPCode: res.StatusCode,
		Body:     string(slurp),
		Header:   res.Header,
	}
}
