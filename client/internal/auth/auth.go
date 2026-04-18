package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/z3vxo/kronos/internal/config"
)

type AUTH struct {
	Token   string
	Refresh string
	baseURL string
	Client  *http.Client
}

type AuthResp struct {
	Tok     string `json:"token"`
	Refresh string `json:"refresh"`
}

func New(baseURL, User, Passwd string, c *http.Client) (*AUTH, error) {
	a := &AUTH{
		baseURL: config.Cfg.Http.Host,
		Client:  c,
	}
	if err := a.Login(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *AUTH) Login() error {
	endpoint := fmt.Sprintf("%s/ts/rest/login", a.baseURL)
	fmt.Printf("User: %s\n", config.Cfg.Http.User)
	fmt.Printf("Passwd: %s\n", config.Cfg.Http.Passwd)

	body, _ := json.Marshal(map[string]string{"username": config.Cfg.Http.User, "password": config.Cfg.Http.Passwd})
	resp, err := a.Client.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed logging in")
	}
	defer resp.Body.Close()
	var auth AuthResp
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return err
	}

	a.Token = auth.Tok
	a.Refresh = auth.Refresh
	fmt.Printf("Token: %s\nRefresh: %s\n", auth.Tok, auth.Refresh)
	return nil
}

func (a *AUTH) Apply(req *http.Request) {
	TokenString := fmt.Sprintf("Bearer %s", a.Token)
	req.Header.Set("Authorization", TokenString)
	return
}
