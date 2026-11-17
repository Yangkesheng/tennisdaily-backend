package model

type AdminAnalyticsDay struct {
	Date                string `json:"date"`
	NewUsers            int64  `json:"newUsers"`
	NewSessions         int64  `json:"newSessions"`
	NewRackets          int64  `json:"newRackets"`
	NewStringingRecords int64  `json:"newStringingRecords"`
	NewShoes            int64  `json:"newShoes"`
	ActiveUsers         int64  `json:"activeUsers"`
	UnregisteredOpeners int64  `json:"unregisteredOpeners"`
}

type AdminAnalyticsResponse struct {
	StartDate           string              `json:"startDate"`
	EndDate             string              `json:"endDate"`
	TimeZone            string              `json:"timeZone"`
	GeneratedAt         string              `json:"generatedAt"`
	TotalUsers          int64               `json:"totalUsers"`
	UnregisteredOpeners int64               `json:"unregisteredOpeners"`
	Days                []AdminAnalyticsDay `json:"days"`
}
