package service

import (
	"strings"

	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

const (
	maxFeedbackContentLength = 1000
	maxFeedbackContactLength = 100
)

type FeedbackService struct {
	repo            *repository.FeedbackRepository
	userRepo        *repository.UserRepository
	contentSecurity *ContentSecurityService
}

func NewFeedbackService(repo *repository.FeedbackRepository, userRepo *repository.UserRepository, contentSecurity *ContentSecurityService) *FeedbackService {
	return &FeedbackService{repo: repo, userRepo: userRepo, contentSecurity: contentSecurity}
}

func (s *FeedbackService) Create(userID int64, req model.CreateFeedbackRequest) (model.FeedbackResponse, error) {
	content := strings.TrimSpace(req.Content)
	contact := strings.TrimSpace(req.Contact)

	if content == "" {
		return model.FeedbackResponse{}, NewInvalidRequestError("反馈内容不能为空")
	}
	if len([]rune(content)) > maxFeedbackContentLength {
		return model.FeedbackResponse{}, NewInvalidRequestError("反馈内容过长，请控制在1000字以内")
	}
	if len([]rune(contact)) > maxFeedbackContactLength {
		return model.FeedbackResponse{}, NewInvalidRequestError("联系方式过长")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return model.FeedbackResponse{}, err
	}
	if user == nil {
		return model.FeedbackResponse{}, ErrUnauthorized
	}

	if err := s.contentSecurity.CheckTexts(user.OpenID, content, contact); err != nil {
		return model.FeedbackResponse{}, err
	}

	feedback := model.Feedback{
		UserID:  userID,
		Content: content,
		Contact: contact,
		Status:  model.FeedbackStatusPending,
	}
	if err := s.repo.Create(&feedback); err != nil {
		return model.FeedbackResponse{}, err
	}
	return model.NewFeedbackResponse(feedback), nil
}
