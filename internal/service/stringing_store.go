package service

import (
	"context"
	"tennisdaily-backend/internal/model"
	"time"
	"unicode/utf8"
)

func normalizeStringingStore(name string, field model.LocationField) (*model.CourtLocation, error) {
	name = model.TrimCourtText(name)
	if !field.Present {
		return nil, NewInvalidRequestError("请提交穿线地点信息")
	}
	location := field.Value
	if utf8.RuneCountInString(name) > 100 {
		return nil, NewInvalidRequestError("穿线场所名称最多 100 字")
	}
	if location == nil {
		if name != "" {
			return nil, NewInvalidRequestError("请从地图选择穿线地点")
		}
		return nil, nil
	}
	value := *location
	value.Address = model.TrimCourtText(value.Address)
	if name == "" || !value.Valid() {
		return nil, NewInvalidRequestError("穿线地点有误，请重新选择")
	}
	value.Latitude = model.RoundCoordinate(value.Latitude)
	value.Longitude = model.RoundCoordinate(value.Longitude)
	return &value, nil
}

func (s *RacketService) RecentStringingStores(ctx context.Context, userID int64) ([]model.StringingStore, error) {
	stores := make([]model.StringingStore, 0, 3)
	seen := make(map[courtIdentity]bool)
	var cursorDate *time.Time
	var cursorID int64
	for {
		rows, err := s.repo.StringingStoreHistory(ctx, userID, cursorDate, cursorID)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			name := model.TrimCourtText(row.StoreName)
			location := row.StoreColumns.Location()
			if name == "" || location == nil || !location.Valid() {
				continue
			}
			key := identityOf(model.CourtSnapshot{CourtName: name, CourtLocation: location})
			if seen[key] {
				continue
			}
			seen[key] = true
			stores = append(stores, model.StringingStore{StoreName: name, StoreLocation: location})
			if len(stores) == 3 {
				return stores, nil
			}
		}
		if len(rows) < 100 {
			return stores, nil
		}
		last := rows[len(rows)-1]
		cursorDate, cursorID = &last.StringDate, last.ID
	}
}
