package model

type SessionStats struct {
	MonthCount   int     `json:"monthCount"`
	MonthMinutes int     `json:"monthMinutes"`
	MonthCost    float64 `json:"monthCost"`
	TotalCount   int     `json:"totalCount"`
}
