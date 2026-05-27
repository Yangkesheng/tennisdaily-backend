package model

type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"userId"`
	OpenID string `json:"openid"`
}
