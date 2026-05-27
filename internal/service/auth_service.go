package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"tennisdaily-backend/internal/config"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	cfg        config.Config
	userRepo   *repository.UserRepository
	httpClient *http.Client
}

type Claims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

type wechatCode2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func NewAuthService(cfg config.Config, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userRepo: userRepo,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *AuthService) WechatLogin(code string) (model.LoginResponse, error) {
	if code == "" {
		return model.LoginResponse{}, ErrInvalidRequest
	}

	openid, err := s.code2OpenID(code)
	if err != nil {
		return model.LoginResponse{}, err
	}

	user, err := s.userRepo.FindOrCreateByOpenID(openid)
	if err != nil {
		return model.LoginResponse{}, err
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{
		Token:  token,
		UserID: user.ID,
		OpenID: user.OpenID,
	}, nil
}

func (s *AuthService) GenerateToken(userID int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.JWTExpire)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnauthorized
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return 0, ErrUnauthorized
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.UserID <= 0 {
		return 0, ErrUnauthorized
	}
	return claims.UserID, nil
}

func (s *AuthService) code2OpenID(code string) (string, error) {
	if s.cfg.WechatAppID == "" || s.cfg.WechatSecret == "" {
		// 方便本地 MVP 联调：未配置微信密钥时使用稳定的开发 openid。
		return "dev_openid_" + code, nil
	}

	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.cfg.WechatAppID,
		s.cfg.WechatSecret,
		code,
	)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result wechatCode2SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 || result.OpenID == "" {
		return "", fmt.Errorf("wechat code2session failed: %d %s", result.ErrCode, result.ErrMsg)
	}
	return result.OpenID, nil
}
