package httpclient

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/z3vxo/kronos/internal/auth"
	"github.com/z3vxo/kronos/internal/config"
)

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

	client := &Client{
		HttpClient: c,
		Auth:       a,
		Hostname:   config.Cfg.Http.Host,
		Stream:     s,
	}

	return client, nil
}

func (c *Client) ConnectToSSE() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ts/events", c.Hostname), nil)
	if err != nil {
		return err
	}
	c.Auth.Apply(req)

	resp, err := c.Stream.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("SSE connect failed: %s", resp.Status)
	}

	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		raw := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		var event Event
		if err := json.Unmarshal([]byte(raw), &event); err != nil {
			continue
		}
		switch event.CmdType {
		case TYPE_NEW_AGENT:
			fmt.Println("[+] New agent:", event.User.CodeName)
		case TYPE_CMD_OUTPUT:
			fmt.Println(event.Data.Output)
		}
	}
	return nil
}

func (c *Client) DoGet(endpoint string, out any) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.Hostname, endpoint), nil)
	if err != nil {
		return err
	}
	c.Auth.Apply(req)
	return c.Do(req, out)
}

func (c *Client) DoPost(endpoint string, data []byte, out any) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.Hostname, endpoint), bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	c.Auth.Apply(req)
	return c.Do(req, out)

}

func (c *Client) Do(req *http.Request, out any) error {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var e ErrorRes
		_ = json.NewDecoder(resp.Body).Decode(&e)
		return fmt.Errorf("[!] Error: %s", e.ErrorStr)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
