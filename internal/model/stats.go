package model

const (
	StatsPeriodMonth = "month"
	StatsPeriodYear  = "year"
)

type SessionStats struct {
	MonthCount   int     `json:"monthCount"`
	MonthMinutes int     `json:"monthMinutes"`
	MonthCost    float64 `json:"monthCost"`
	TotalCount   int     `json:"totalCount"`
}

type StatsChartsQuery struct {
	Period string `form:"period" binding:"required,oneof=month year"`
	Year   int    `form:"year" binding:"required,min=2000,max=2100"`
	Month  int    `form:"month" binding:"omitempty,min=1,max=12"`
}

type StatsChartsSummaryResponse struct {
	SessionCount   int64   `json:"sessionCount"`
	ActiveDayCount int64   `json:"activeDayCount"`
	TotalMinutes   int64   `json:"totalMinutes"`
	AverageMinutes float64 `json:"averageMinutes"`
	AverageRating  float64 `json:"averageRating"`
	SessionCost    float64 `json:"sessionCost"`
	RacketCost     float64 `json:"racketCost"`
	StringingCost  float64 `json:"stringingCost"`
	TotalCost      float64 `json:"totalCost"`
	TrainingCount  int64   `json:"trainingCount"`
	SinglesCount   int64   `json:"singlesCount"`
	DoublesCount   int64   `json:"doublesCount"`
	MatchCount     int64   `json:"matchCount"`
}

type StatsFrequencyChartItemResponse struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

type StatsRatingTrendItemResponse struct {
	Label  string  `json:"label"`
	Date   string  `json:"date"`
	Rating float64 `json:"rating"`
}

type StatsBreakdownItemResponse struct {
	Key     string  `json:"key"`
	Label   string  `json:"label"`
	Value   float64 `json:"value"`
	Percent float64 `json:"percent"`
}

type StatsSessionBreakdownRow struct {
	Category    SessionCategory    `gorm:"column:category"`
	SubCategory SessionSubCategory `gorm:"column:sub_category"`
	Count       int64              `gorm:"column:count"`
	Minutes     int64              `gorm:"column:minutes"`
	Cost        float64            `gorm:"column:cost"`
}

type StatsChartsResponse struct {
	Frequency                           []StatsFrequencyChartItemResponse `json:"frequency"`
	RatingTrend                         []StatsRatingTrendItemResponse    `json:"ratingTrend"`
	ExpenseBreakdown                    []StatsBreakdownItemResponse      `json:"expenseBreakdown"`
	SessionTypeBreakdown                []StatsBreakdownItemResponse      `json:"sessionTypeBreakdown"`
	SessionCategoryCountBreakdown       []StatsBreakdownItemResponse      `json:"sessionCategoryCountBreakdown"`
	SessionSubCategoryCountBreakdown    []StatsBreakdownItemResponse      `json:"sessionSubCategoryCountBreakdown"`
	SessionCategoryDurationBreakdown    []StatsBreakdownItemResponse      `json:"sessionCategoryDurationBreakdown"`
	SessionSubCategoryDurationBreakdown []StatsBreakdownItemResponse      `json:"sessionSubCategoryDurationBreakdown"`
	SessionCategoryCostBreakdown        []StatsBreakdownItemResponse      `json:"sessionCategoryCostBreakdown"`
	SessionSubCategoryCostBreakdown     []StatsBreakdownItemResponse      `json:"sessionSubCategoryCostBreakdown"`
}

type StatsChartsResultResponse struct {
	Period    string                     `json:"period"`
	Year      int                        `json:"year"`
	Month     int                        `json:"month"`
	RangeText string                     `json:"rangeText"`
	Summary   StatsChartsSummaryResponse `json:"summary"`
	Charts    StatsChartsResponse        `json:"charts"`
}

type StatsSessionAggregate struct {
	SessionCount   int64   `gorm:"column:session_count"`
	ActiveDayCount int64   `gorm:"column:active_day_count"`
	TotalMinutes   int64   `gorm:"column:total_minutes"`
	AverageMinutes float64 `gorm:"column:average_minutes"`
	AverageRating  float64 `gorm:"column:average_rating"`
	SessionCost    float64 `gorm:"column:session_cost"`
	TrainingCount  int64   `gorm:"column:training_count"`
	SinglesCount   int64   `gorm:"column:singles_count"`
	DoublesCount   int64   `gorm:"column:doubles_count"`
	MatchCount     int64   `gorm:"column:match_count"`
}

type StatsMonthlyFrequency struct {
	Month int   `gorm:"column:month"`
	Count int64 `gorm:"column:count"`
}

type StatsMonthlyRating struct {
	Month  int     `gorm:"column:month"`
	Rating float64 `gorm:"column:rating"`
}
