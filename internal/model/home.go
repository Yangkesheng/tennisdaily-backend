package model

type HomeSummaryQuery struct {
	Year  int `form:"year" binding:"omitempty,min=2000,max=2100"`
	Month int `form:"month" binding:"omitempty,min=1,max=12"`
}

type HomeSessionSummaryResponse struct {
	MonthCount   int     `json:"monthCount"`
	MonthMinutes int     `json:"monthMinutes"`
	MonthCost    float64 `json:"monthCost"`
	TotalCount   int     `json:"totalCount"`
}

type HomeExpenseSummaryResponse struct {
	SessionCost   float64 `json:"sessionCost"`
	RacketCost    float64 `json:"racketCost"`
	StringingCost float64 `json:"stringingCost"`
	TotalCost     float64 `json:"totalCost"`
}

type HomeRatingTrendItemResponse struct {
	ID     int64  `json:"id"`
	Date   string `json:"date"`
	Rating int16  `json:"rating"`
}

type HomeSummaryResponse struct {
	Year          int                           `json:"year"`
	Month         int                           `json:"month"`
	Session       HomeSessionSummaryResponse    `json:"session"`
	Expense       HomeExpenseSummaryResponse    `json:"expense"`
	LatestSession *SessionResponse              `json:"latestSession"`
	RatingTrend   []HomeRatingTrendItemResponse `json:"ratingTrend"`
}
