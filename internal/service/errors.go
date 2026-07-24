package service

import (
	"errors"
	"fmt"
	"strings"
)

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

func NewInvalidRequestError(message string) error {
	if message == "" {
		return ErrInvalidRequest
	}
	return fmt.Errorf("%w: %s", ErrInvalidRequest, message)
}

func InvalidRequestMessage(err error) string {
	if !errors.Is(err, ErrInvalidRequest) {
		return "invalid request"
	}

	prefix := ErrInvalidRequest.Error() + ": "
	message := strings.TrimPrefix(err.Error(), prefix)
	if message == err.Error() || message == "" {
		return ErrInvalidRequest.Error()
	}
	return message
}
