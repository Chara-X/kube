package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"k8s.io/client-go/rest"
)

type Client struct {
	*http.Client
	cfg *rest.Config
}

func New(cfg *rest.Config) *Client {
	var cli, _ = rest.HTTPClientFor(cfg)
	return &Client{cli, cfg}
}
func (c *Client) Post(url *Url, body any) (*http.Response, error) { return c.Do("POST", url, body) }
func (c *Client) Put(url *Url, body any) (*http.Response, error)  { return c.Do("PUT", url, body) }
func (c *Client) Delete(url *Url) (*http.Response, error)         { return c.Do("DELETE", url, nil) }
func (c *Client) Get(url *Url) (*http.Response, error)            { return c.Do("GET", url, nil) }
func (c *Client) Do(method string, url *Url, body any) (*http.Response, error) {
	var buf, _ = json.Marshal(body)
	var req, _ = http.NewRequest(method, c.cfg.Host+url.String(), bytes.NewBuffer(buf))
	return c.Client.Do(req)
}
