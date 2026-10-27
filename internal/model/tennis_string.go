package model

// StringLibrary stores independently reviewed product data; it does not alter personal records.
type StringLibrary struct {
	ID        int64           `json:"id" gorm:"primaryKey"`
	Brand     string          `json:"brand"`
	Model     string          `json:"model"`
	NameZh    string          `json:"nameZh"`
	Aliases   string          `json:"-"`
	Material  string          `json:"material"`
	SourceURL string          `json:"sourceUrl"`
	Published bool            `json:"-"`
	Variants  []StringVariant `json:"variants" gorm:"foreignKey:StringID"`
}

func (StringLibrary) TableName() string { return "string_library" }

type StringVariant struct {
	ID         int64   `json:"id" gorm:"primaryKey"`
	StringID   int64   `json:"-"`
	DiameterMm float64 `json:"diameterMm"`
}

func (StringVariant) TableName() string { return "string_variants" }

type StringLibraryQuery struct {
	StringID int64  `form:"stringId" binding:"omitempty,min=1"`
	Keyword  string `form:"keyword" binding:"max=100"`
	Brand    string `form:"brand" binding:"max=100"`
	Page     int    `form:"page" binding:"omitempty,min=1,max=10000"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=50"`
}
type StringLibraryPage struct {
	Items    []StringLibrary `json:"items"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Total    int64           `json:"total"`
}
