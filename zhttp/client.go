package zhttp

import (
	"net/http"
	"time"
)

var DefaultClient = &Client{
	Client: &http.Client{
		Transport: Transport(false, nil),
	},
}

type (
	Client struct {
		*http.Client
	}

	ClientOption func(c *Client)
)

func WithTransport(t *http.Transport) ClientOption {
	return func(c *Client) {
		c.Transport = t
	}
}

func WithCheckRedirect(redirect func(*http.Request, []*http.Request) error) ClientOption {
	return func(c *Client) {
		c.CheckRedirect = redirect
	}
}

func WithJar(jar http.CookieJar) ClientOption {
	return func(c *Client) {
		c.Jar = jar
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.Timeout = timeout
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{}
	c.Transport = Transport(false, nil)
	for _, opt := range opts {
		opt(c)
	}
	return c
}
