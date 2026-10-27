package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"

	"gorm.io/gorm"
)

var ErrInvalidCourtLocation = errors.New("invalid court location")

// LocationField distinguishes omitted fields from explicit JSON null.
type LocationField struct {
	Present bool
	Value   *CourtLocation
}

func (f *LocationField) UnmarshalJSON(data []byte) error {
	f.Present = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		f.Value = nil
		return nil
	}
	var value CourtLocation
	if err := json.Unmarshal(data, &value); err != nil {
		return ErrInvalidCourtLocation
	}
	f.Value = &value
	return nil
}

type CourtLocation struct {
	Address          string  `json:"address"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	CoordinateSystem string  `json:"coordinateSystem"`
}

func (l *CourtLocation) UnmarshalJSON(data []byte) error {
	var raw struct {
		Address          *string  `json:"address"`
		Latitude         *float64 `json:"latitude"`
		Longitude        *float64 `json:"longitude"`
		CoordinateSystem *string  `json:"coordinateSystem"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return ErrInvalidCourtLocation
	}
	if raw.Address == nil || raw.Latitude == nil || raw.Longitude == nil || raw.CoordinateSystem == nil {
		return ErrInvalidCourtLocation
	}
	*l = CourtLocation{*raw.Address, *raw.Latitude, *raw.Longitude, *raw.CoordinateSystem}
	return nil
}
func TrimCourtText(s string) string {
	return strings.TrimFunc(s, func(r rune) bool { return unicode.IsSpace(r) || r == '\ufeff' })
}
func RoundCoordinate(n float64) float64 {
	if n == 0 {
		return 0
	}
	rounded := math.Round(n*1e7) / 1e7
	if rounded == 0 {
		return 0
	}
	return rounded
}
func (l CourtLocation) Valid() bool {
	return !math.IsNaN(l.Latitude) && !math.IsInf(l.Latitude, 0) && l.Latitude >= -90 && l.Latitude <= 90 &&
		!math.IsNaN(l.Longitude) && !math.IsInf(l.Longitude, 0) && l.Longitude >= -180 && l.Longitude <= 180 &&
		l.CoordinateSystem == "gcj02" && utf8.RuneCountInString(l.Address) <= 512
}

type CourtSnapshot struct {
	CourtName     string         `json:"courtName"`
	CourtLocation *CourtLocation `json:"courtLocation"`
}
type CourtRecommendations struct {
	LatestCourt  *CourtSnapshot  `json:"latestCourt"`
	RecentCourts []CourtSnapshot `json:"recentCourts"`
}

// CourtColumns is embedded in the database model; API responses use CourtLocation.
type CourtColumns struct {
	CourtAddress          string   `json:"-" gorm:"size:512;not null;default:''"`
	CourtLatitude         *float64 `json:"-" gorm:"type:decimal(10,7)"`
	CourtLongitude        *float64 `json:"-" gorm:"type:decimal(10,7)"`
	CourtCoordinateSystem string   `json:"-" gorm:"size:16;not null;default:''"`
}

func (c CourtColumns) Location() *CourtLocation {
	if c.CourtLatitude == nil || c.CourtLongitude == nil {
		return nil
	}
	return &CourtLocation{c.CourtAddress, *c.CourtLatitude, *c.CourtLongitude, c.CourtCoordinateSystem}
}
func (c CourtColumns) Validate() error {
	if c.CourtLatitude == nil && c.CourtLongitude == nil && c.CourtAddress == "" && c.CourtCoordinateSystem == "" {
		return nil
	}
	l := c.Location()
	if l == nil || !l.Valid() {
		return errors.New("invalid stored court location")
	}
	return nil
}
func ColumnsFromLocation(l *CourtLocation) CourtColumns {
	if l == nil {
		return CourtColumns{}
	}
	lat, lng := l.Latitude, l.Longitude
	return CourtColumns{l.Address, &lat, &lng, l.CoordinateSystem}
}
func (s *TennisSession) AfterFind(_ *gorm.DB) error {
	if err := s.CourtColumns.Validate(); err != nil {
		return err
	}
	if s.CourtColumns.Location() != nil && TrimCourtText(s.CourtName) == "" {
		return errors.New("stored court name missing")
	}
	return nil
}
