package chatgpt

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type Config struct {
	ClientID     string
	Scope        string
	RefreshToken string
	Origin       string
	Endpoint     string
}

type ChatGPT struct {
	cfg      *Config
	logger   log.Logger
	Resp     *RefreshResp
	Endpoint url.URL
}

type RefreshResp struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

func (chat *ChatGPT) Refresh() error {
	payload := url.Values{
		"client_id":     {chat.cfg.ClientID},
		"scope":         {chat.cfg.Scope},
		"refresh_token": {chat.cfg.RefreshToken},
		"grant_type":    {"refresh_token"},
	}.Encode()
	req, _ := http.NewRequest(
		"POST",
		"https://login.microsoftonline.com/common/oauth2/v2.0/token",
		strings.NewReader(payload),
	)
	req.Header.Add("Origin", chat.cfg.Origin)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var r RefreshResp
	_ = json.Unmarshal(body, &r)

	if err != nil {
		return err
	}

	chat.Resp = &r
	chat.cfg.RefreshToken = r.RefreshToken

	return nil
}

func New(config *Config, logger log.Logger) (*ChatGPT, error) {
	endpoint, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, err
	}

	mstoken := ChatGPT{
		cfg:      config,
		logger:   logger,
		Endpoint: *endpoint,
	}

	logger.Log(log.LevelInfo, "msg", "[Microsoft Token]: Try refresh token")
	err = mstoken.Refresh()

	if err != nil {
		return nil, err
	}

	return &mstoken, nil
}
