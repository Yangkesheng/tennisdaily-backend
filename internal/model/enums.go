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
	Categories    []SessionCategoryEnumItemResponse `json:"categories"`
	LegacyTypes   []EnumItemResponse                `json:"legacyTypes"`
	MatchRanks    []EnumItemResponse                `json:"matchRanks"`
	DefaultValues SessionDefaultEnumsResponse       `json:"defaultValues"`
}

type SessionDefaultEnumsResponse struct {
	Category        int16 `json:"category"`
	SubCategory     int16 `json:"subCategory"`
	DurationMinutes int   `json:"durationMinutes"`
	Rating          int16 `json:"rating"`
	MatchRank       int16 `json:"matchRank"`
}
