package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"tennisdaily-backend/internal/config"
)

const contentSecuritySceneProfile = 1

type ContentSecurityService struct {
	cfg        config.Config
	httpClient *http.Client
}

type contentSecurityAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type msgSecCheckResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  struct {
		Suggest string `json:"suggest"`
		Label   int    `json:"label"`
	} `json:"result"`
}

func NewContentSecurityService(cfg config.Config) *ContentSecurityService {
	return &ContentSecurityService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *ContentSecurityService) CheckTexts(openid string, texts ...string) error {
	if s == nil || s.cfg.WechatAppID == "" || s.cfg.WechatSecret == "" {
		return nil
	}
	if openid == "" {
		return ErrInvalidRequest
	}

	for _, text := range texts {
		content := strings.TrimSpace(text)
		if content == "" {
			continue
		}
		if err := s.checkText(openid, content); err != nil {
			return err
		}
	}
	return nil
}

func (s *ContentSecurityService) checkText(openid, content string) error {
	accessToken, err := s.getAccessToken()
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"openid":  openid,
		"scene":   contentSecuritySceneProfile,
		"version": 2,
		"content": content,
	})
	url := "https://api.weixin.qq.com/wxa/msg_sec_check?access_token=" + accessToken
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return ErrWechatService
	}
	defer resp.Body.Close()

	var result msgSecCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ErrWechatService
	}
	if result.ErrCode != 0 {
		if result.ErrCode == 87014 {
			return ErrSensitiveContent
		}
		return ErrWechatService
	}
	if result.Result.Suggest == "risky" {
		return ErrSensitiveContent
	}
	return nil
}

func (s *ContentSecurityService) getAccessToken() (string, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		s.cfg.WechatAppID,
		s.cfg.WechatSecret,
	)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return "", ErrWechatService
	}
	defer resp.Body.Close()

	var result contentSecurityAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", ErrWechatService
	}
	if result.ErrCode != 0 || result.AccessToken == "" {
		return "", ErrWechatService
	}
	return result.AccessToken, nil
}
