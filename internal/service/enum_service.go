package service

import (
	"tennisdaily-backend/internal/config"
	"tennisdaily-backend/internal/model"
)

type EnumService struct {
	sessionCategoryResolver *config.SessionCategoryResolver
}

func NewEnumService(sessionCategoryResolver *config.SessionCategoryResolver) *EnumService {
	return &EnumService{sessionCategoryResolver: sessionCategoryResolver}
}

func (s *EnumService) SessionConfig() model.SessionEnumsResponse {
	return model.SessionEnumsResponse{
		Categories: s.sessionCategories(),
		LegacyTypes: []model.EnumItemResponse{
			{Value: int16(model.SessionTypeDoubles), Label: model.SessionTypeDoubles.Label()},
			{Value: int16(model.SessionTypeSingles), Label: model.SessionTypeSingles.Label()},
			{Value: int16(model.SessionTypeTraining), Label: model.SessionTypeTraining.Label()},
			{Value: int16(model.SessionTypeSinglesMatch), Label: model.SessionTypeSinglesMatch.Label()},
			{Value: int16(model.SessionTypeDoublesMatch), Label: model.SessionTypeDoublesMatch.Label()},
		},
		MatchRanks: []model.EnumItemResponse{
			{Value: int16(model.MatchRankNone), Label: model.MatchRankNone.Label()},
			{Value: int16(model.MatchRankChampion), Label: model.MatchRankChampion.Label()},
			{Value: int16(model.MatchRankRunnerUp), Label: model.MatchRankRunnerUp.Label()},
			{Value: int16(model.MatchRankThirdPlace), Label: model.MatchRankThirdPlace.Label()},
			{Value: int16(model.MatchRankSemiFinal), Label: model.MatchRankSemiFinal.Label()},
			{Value: int16(model.MatchRankQuarterFinal), Label: model.MatchRankQuarterFinal.Label()},
			{Value: int16(model.MatchRankRoundOf16), Label: model.MatchRankRoundOf16.Label()},
			{Value: int16(model.MatchRankGroupStage), Label: model.MatchRankGroupStage.Label()},
		},
		DefaultValues: model.SessionDefaultEnumsResponse{
			Category:        int16(model.SessionCategoryDaily),
			SubCategory:     int16(model.SessionSubCategoryDoubles),
			DurationMinutes: defaultDurationMinutes,
			Rating:          defaultRating,
			MatchRank:       int16(model.MatchRankNone),
		},
	}
}

func (s *EnumService) sessionCategories() []model.SessionCategoryEnumItemResponse {
	categories := s.sessionCategoryResolver.Categories()
	items := make([]model.SessionCategoryEnumItemResponse, 0, len(categories))
	for _, category := range categories {
		items = append(items, s.sessionCategoryEnumItem(category))
	}
	return items
}

func (s *EnumService) sessionCategoryEnumItem(category config.SessionCategoryConfig) model.SessionCategoryEnumItemResponse {
	items := make([]model.SessionSubCategoryEnumItemResponse, 0, len(category.SubCategories))
	for _, subCategory := range category.SubCategories {
		categoryValue := model.SessionCategory(category.Value)
		subCategoryValue := model.SessionSubCategory(subCategory.Value)
		option, _ := s.sessionCategoryResolver.Option(categoryValue, subCategoryValue)
		typeText, _ := s.sessionCategoryResolver.TypeText(categoryValue, subCategoryValue)
		items = append(items, model.SessionSubCategoryEnumItemResponse{
			Value:      subCategory.Value,
			Label:      subCategory.Label,
			Category:   category.Value,
			TypeText:   typeText,
			LegacyType: int16(option.LegacyType),
		})
	}
	return model.SessionCategoryEnumItemResponse{
		Value:         category.Value,
		Label:         category.Label,
		SubCategories: items,
	}
}
