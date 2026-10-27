package service

import (
	"context"
	"time"
	"unicode/utf8"

	"tennisdaily-backend/internal/model"
)

func normalizeCourt(name *string, field model.LocationField) (model.CourtSnapshot, error) {
	if name == nil || !field.Present {
		return model.CourtSnapshot{}, NewInvalidRequestError("请更新小程序后重试")
	}
	court := model.CourtSnapshot{CourtName: model.TrimCourtText(*name)}
	if utf8.RuneCountInString(court.CourtName) > 128 {
		return court, NewInvalidRequestError("球场名称最多 128 字")
	}
	if field.Value == nil {
		return court, nil
	}
	l := *field.Value
	l.Address = model.TrimCourtText(l.Address)
	if utf8.RuneCountInString(l.Address) > 512 {
		return court, NewInvalidRequestError("球场地址最多 512 字")
	}
	if court.CourtName == "" || !l.Valid() {
		return court, NewInvalidRequestError("球场位置信息有误，请重新选择")
	}
	l.Latitude = model.RoundCoordinate(l.Latitude)
	l.Longitude = model.RoundCoordinate(l.Longitude)
	court.CourtLocation = &l
	return court, nil
}
func courtAddress(c model.CourtSnapshot) string {
	if c.CourtLocation == nil {
		return ""
	}
	return c.CourtLocation.Address
}

type courtIdentity struct {
	Name, Address, System string
	Latitude, Longitude   float64
	HasLocation           bool
}

func identityOf(c model.CourtSnapshot) courtIdentity {
	key := courtIdentity{Name: model.TrimCourtText(c.CourtName)}
	if l := c.CourtLocation; l != nil {
		key.Address = model.TrimCourtText(l.Address)
		key.System = l.CoordinateSystem
		key.Latitude = model.RoundCoordinate(l.Latitude)
		key.Longitude = model.RoundCoordinate(l.Longitude)
		key.HasLocation = true
	}
	return key
}

// Scan only court projections, stopping at four distinct courts rather than four sessions.
func (s *SessionService) CourtRecommendations(ctx context.Context, userID int64) (model.CourtRecommendations, error) {
	var cursorDate *time.Time
	var cursorID int64
	result := model.CourtRecommendations{RecentCourts: make([]model.CourtSnapshot, 0, 4)}
	seen := make(map[courtIdentity]bool)
	first := true
	for {
		if err := ctx.Err(); err != nil {
			return model.CourtRecommendations{}, err
		}
		rows, err := s.repo.CourtHistoryPage(ctx, userID, cursorDate, cursorID, 100)
		if err != nil {
			return model.CourtRecommendations{}, err
		}
		for _, row := range rows {
			c := model.CourtSnapshot{CourtName: model.TrimCourtText(row.CourtName), CourtLocation: row.CourtColumns.Location()}
			if first {
				if c.CourtName != "" {
					latest := c
					result.LatestCourt = &latest
				}
				first = false
			}
			if c.CourtName == "" {
				continue
			}
			key := identityOf(c)
			if seen[key] {
				continue
			}
			seen[key] = true
			result.RecentCourts = append(result.RecentCourts, c)
			if len(result.RecentCourts) == 4 {
				return result, nil
			}
		}
		if len(rows) < 100 {
			return result, nil
		}
		last := rows[len(rows)-1]
		cursorDate = &last.Date
		cursorID = last.ID
	}
}
