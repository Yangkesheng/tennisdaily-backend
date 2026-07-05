package model

type EnumItemResponse struct {
	Value int16  `json:"value"`
	Label string `json:"label"`
}

type SessionSubCategoryEnumItemResponse struct {
	Value      int16  `json:"value"`
	Label      string `json:"label"`
	Category   int16  `json:"category"`
	TypeText   string `json:"typeText"`
	LegacyType int16  `json:"legacyType"`
}

type SessionCategoryEnumItemResponse struct {
	Value         int16                                `json:"value"`
	Label         string                               `json:"label"`
	SubCategories []SessionSubCategoryEnumItemResponse `json:"subCategories"`
}

type SessionEnumsResponse struct {
	Categories []SessionCategoryEnumItemResponse `json:"categories"`
	MatchRanks []EnumItemResponse                `json:"matchRanks"`
}
