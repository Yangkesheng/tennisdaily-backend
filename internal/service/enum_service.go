package service

import "tennisdaily-backend/internal/model"

type EnumService struct{}

func NewEnumService() *EnumService {
	return &EnumService{}
}

func (s *EnumService) All() model.EnumsResponse {
	return model.EnumsResponse{
		Session: model.SessionEnumsResponse{
			Categories: []model.SessionCategoryEnumItemResponse{
				sessionCategoryEnumItem(model.SessionCategoryDaily, []model.SessionSubCategory{model.SessionSubCategorySingles, model.SessionSubCategoryDoubles}),
				sessionCategoryEnumItem(model.SessionCategoryTraining, []model.SessionSubCategory{model.SessionSubCategoryServe, model.SessionSubCategoryOther}),
				sessionCategoryEnumItem(model.SessionCategoryMatch, []model.SessionSubCategory{model.SessionSubCategorySingles, model.SessionSubCategoryDoubles}),
			},
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
				{Value: int16(model.MatchRankSemiFinal), Label: model.MatchRankSemiFinal.Label()},
				{Value: int16(model.MatchRankQuarterFinal), Label: model.MatchRankQuarterFinal.Label()},
				{Value: int16(model.MatchRankGroupStage), Label: model.MatchRankGroupStage.Label()},
			},
			DefaultValues: model.SessionDefaultEnumsResponse{
				Category:        int16(model.SessionCategoryDaily),
				SubCategory:     int16(model.SessionSubCategoryDoubles),
				DurationMinutes: defaultDurationMinutes,
				Rating:          defaultRating,
				MatchRank:       int16(model.MatchRankNone),
			},
		},
		Racket: model.RacketEnumsResponse{
			Statuses: []model.EnumItemResponse{
				{Value: int16(model.RacketStatusPrimary), Label: model.RacketStatusPrimary.Label()},
				{Value: int16(model.RacketStatusActive), Label: model.RacketStatusActive.Label()},
				{Value: int16(model.RacketStatusRetired), Label: model.RacketStatusRetired.Label()},
			},
		},
	}
}

func sessionCategoryEnumItem(category model.SessionCategory, subCategories []model.SessionSubCategory) model.SessionCategoryEnumItemResponse {
	items := make([]model.SessionSubCategoryEnumItemResponse, 0, len(subCategories))
	for _, subCategory := range subCategories {
		legacyType, _ := category.ToSessionType(subCategory)
		items = append(items, model.SessionSubCategoryEnumItemResponse{
			Value:      int16(subCategory),
			Label:      subCategory.Label(category),
			Category:   int16(category),
			TypeText:   category.TypeText(subCategory),
			LegacyType: int16(legacyType),
		})
	}
	return model.SessionCategoryEnumItemResponse{
		Value:         int16(category),
		Label:         category.Label(),
		SubCategories: items,
	}
}
