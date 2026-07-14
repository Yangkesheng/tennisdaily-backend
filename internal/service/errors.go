package service

import "errors"

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrNotFound         = errors.New("not found")
	ErrWechatLoginCode  = errors.New("wechat login code invalid")
	ErrWechatPhoneCode  = errors.New("wechat phone code invalid")
	ErrWechatService    = errors.New("wechat service unavailable")
	ErrSensitiveContent = errors.New("sensitive content")
)
