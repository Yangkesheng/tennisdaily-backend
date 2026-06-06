package service

import (
	"bytes"
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

type wechatAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type wechatPhoneResponse struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
	} `json:"phone_info"`
}

type wechatIdentity struct {
	OpenID  string
	UnionID string
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

func (s *AuthService) PhoneLogin(req model.PhoneLoginRequest) (model.PhoneLoginResponse, error) {
	if req.LoginCode == "" || req.PhoneCode == "" {
		return model.PhoneLoginResponse{}, ErrInvalidRequest
	}

	identity, err := s.code2WechatIdentity(req.LoginCode)
	if err != nil {
		return model.PhoneLoginResponse{}, err
	}
	phone, err := s.phoneCode2Phone(req.PhoneCode)
	if err != nil {
		return model.PhoneLoginResponse{}, err
	}

	appID := s.cfg.WechatAppID
	if appID == "" {
		appID = "dev_appid"
	}
	user, isNewUser, err := s.userRepo.FindOrCreateByPhoneAndWechat(phone, appID, identity.OpenID, identity.UnionID)
	if err != nil {
		return model.PhoneLoginResponse{}, err
	}
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return model.PhoneLoginResponse{}, err
	}

	return model.PhoneLoginResponse{Token: token, User: model.NewUserResponse(*user), IsNewUser: isNewUser}, nil
}

func (s *AuthService) CurrentUser(userID int64) (model.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return model.UserResponse{}, err
	}
	if user == nil {
		return model.UserResponse{}, ErrUnauthorized
	}
	return model.NewUserResponse(*user), nil
}

func (s *AuthService) UpdateProfile(userID int64, req model.UpdateProfileRequest) (model.UserResponse, error) {
	user, err := s.userRepo.UpdateProfile(userID, req.Nickname, req.AvatarURL)
	if err != nil {
		return model.UserResponse{}, err
	}
	if user == nil {
		return model.UserResponse{}, ErrNotFound
	}
	return model.NewUserResponse(*user), nil
}

func (s *AuthService) Logout() map[string]bool {
	return map[string]bool{"success": true}
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

func (s *AuthService) code2WechatIdentity(code string) (wechatIdentity, error) {
	if s.cfg.WechatAppID == "" || s.cfg.WechatSecret == "" {
		return wechatIdentity{OpenID: "dev_openid_" + code}, nil
	}

	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.cfg.WechatAppID,
		s.cfg.WechatSecret,
		code,
	)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return wechatIdentity{}, ErrWechatService
	}
	defer resp.Body.Close()

	var result wechatCode2SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return wechatIdentity{}, ErrWechatService
	}
	if result.ErrCode != 0 || result.OpenID == "" {
		return wechatIdentity{}, ErrWechatLoginCode
	}
	return wechatIdentity{OpenID: result.OpenID, UnionID: result.UnionID}, nil
}

func (s *AuthService) code2OpenID(code string) (string, error) {
	identity, err := s.code2WechatIdentity(code)
	if err != nil {
		return "", err
	}
	return identity.OpenID, nil
}

func (s *AuthService) phoneCode2Phone(code string) (string, error) {
	if s.cfg.WechatAppID == "" || s.cfg.WechatSecret == "" {
		if len(code) >= 11 && code[:3] == "dev" {
			return "13800138000", nil
		}
		return "13800138000", nil
	}

	accessToken, err := s.getWechatAccessToken()
	if err != nil {
		return "", err
	}
	payload, _ := json.Marshal(map[string]string{"code": code})
	url := "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=" + accessToken
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", ErrWechatService
	}
	defer resp.Body.Close()

	var result wechatPhoneResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", ErrWechatService
	}
	if result.ErrCode != 0 || result.PhoneInfo.PurePhoneNumber == "" {
		return "", ErrWechatPhoneCode
	}
	return result.PhoneInfo.PurePhoneNumber, nil
}

func (s *AuthService) getWechatAccessToken() (string, error) {
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

	var result wechatAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", ErrWechatService
	}
	if result.ErrCode != 0 || result.AccessToken == "" {
		return "", ErrWechatService
	}
	return result.AccessToken, nil
}
