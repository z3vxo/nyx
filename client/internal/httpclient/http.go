package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/z3vxo/kronos/internal/auth"
	"github.com/z3vxo/kronos/internal/config"
)

type Client struct {
	Hostname   string
	Auth       *auth.AUTH
	HttpClient *http.Client
	Stream     *http.Client
}

func NewClient() (*Client, error) {

	c := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	s := &http.Client{
		Timeout: 0,
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}

	a, err := auth.New(config.Cfg.Http.Host, config.Cfg.Http.User, config.Cfg.Http.Passwd, c)
	if err != nil {
		return nil, err
	}

	return &Client{
		HttpClient: c,
		Auth:       a,
		Hostname:   config.Cfg.Http.Host,
		Stream:     s,
	}, nil
}

func (c *Client) DoGet(method string, endpoint string, out any) error {
	req, _ := http.NewRequest(method, fmt.Sprintf("%s/%s", c.Hostname, endpoint), nil)
	c.Auth.Apply(req)
	return c.Do(req, out)
}

func (c *Client) Do(req *http.Request, out any) error {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
