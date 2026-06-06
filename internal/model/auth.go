package model

type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"userId"`
	OpenID string `json:"openid"`
}

type PhoneLoginRequest struct {
	LoginCode string `json:"loginCode" binding:"required"`
	PhoneCode string `json:"phoneCode" binding:"required"`
}

type PhoneLoginResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	IsNewUser bool         `json:"isNewUser"`
}

type UpdateProfileRequest struct {
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}
