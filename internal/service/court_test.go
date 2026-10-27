package service

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"tennisdaily-backend/internal/model"
	"testing"
)

func TestCourtProtocol(t *testing.T) {
	cases := []struct {
		body      string
		wantError bool
		present   bool
	}{
		{`{"courtName":"A"}`, false, false},
		{`{"courtName":"A","courtLocation":null}`, false, true},
		{`{"courtName":"A","courtLocation":{"address":"","latitude":0,"longitude":0,"coordinateSystem":"gcj02"}}`, false, true},
		{`{"courtName":"A","courtLocation":{"address":"","longitude":0,"coordinateSystem":"gcj02"}}`, true, true},
		{`{"courtName":"A","courtLocation":{"address":"","latitude":null,"longitude":0,"coordinateSystem":"gcj02"}}`, true, true},
		{`{"courtName":"A","courtLocation":{"address":"","latitude":"0","longitude":0,"coordinateSystem":"gcj02"}}`, true, true},
		{`{"courtName":"A","courtLocation":[]}`, true, true},
	}
	for _, tc := range cases {
		var req model.CreateSessionRequest
		err := json.Unmarshal([]byte(tc.body), &req)
		if (err != nil) != tc.wantError {
			t.Fatalf("unexpected decode status for protocol case")
		}
		if err == nil && req.CourtLocation.Present != tc.present {
			t.Fatal("field presence lost")
		}
	}
	var req model.UpdateUserSettingsRequest
	if json.Unmarshal([]byte(`{"defaultCourtName":null}`), &req) == nil {
		t.Fatal("null name accepted")
	}
}
func TestNormalizeCourt(t *testing.T) {
	name := "\ufeff\u0085 球场　"
	l := model.CourtLocation{Address: " 地址 ", Latitude: 0, Longitude: -0.00000005, CoordinateSystem: "gcj02"}
	c, err := normalizeCourt(&name, model.LocationField{Present: true, Value: &l})
	if err != nil || c.CourtName != "球场" || c.CourtLocation.Address != "地址" || c.CourtLocation.Longitude != -0.0000001 {
		t.Fatal("normalization failed")
	}
	if l.Address != " 地址 " {
		t.Fatal("mutated input location")
	}
	if _, err := normalizeCourt(&name, model.LocationField{}); !errors.Is(err, ErrInvalidRequest) {
		t.Fatal("omitted location accepted")
	}
	long := strings.Repeat("🎾", 129)
	if _, err := normalizeCourt(&long, model.LocationField{Present: true}); err == nil {
		t.Fatal("long name accepted")
	}
	long = strings.Repeat("🎾", 128)
	if _, err := normalizeCourt(&long, model.LocationField{Present: true}); err != nil {
		t.Fatal(err)
	}
	for _, n := range []float64{91, math.Inf(1), math.NaN()} {
		invalid := l
		invalid.Latitude = n
		if _, err := normalizeCourt(&name, model.LocationField{Present: true, Value: &invalid}); err == nil {
			t.Fatal("invalid coordinate accepted")
		}
	}
	empty := ""
	if _, err := normalizeCourt(&empty, model.LocationField{Present: true, Value: &l}); err == nil {
		t.Fatal("location without name")
	}
	if c, err := normalizeCourt(&empty, model.LocationField{Present: true}); err != nil || c.CourtLocation != nil {
		t.Fatal("clear failed")
	}
}
func TestCourtIdentityAndStoredData(t *testing.T) {
	a := model.CourtSnapshot{CourtName: " A ", CourtLocation: &model.CourtLocation{Address: "one", CoordinateSystem: "gcj02"}}
	b := a
	b.CourtLocation = &model.CourtLocation{Address: "two", CoordinateSystem: "gcj02"}
	if identityOf(a) == identityOf(b) || identityOf(a) == identityOf(model.CourtSnapshot{CourtName: "A"}) {
		t.Fatal("different venues merged")
	}
	columns := model.ColumnsFromLocation(a.CourtLocation)
	if columns.Validate() != nil || columns.Location().Latitude != 0 {
		t.Fatal("zero coordinate lost")
	}
	columns.CourtLongitude = nil
	if columns.Validate() == nil {
		t.Fatal("partial stored location accepted")
	}
}
func TestSettingsResponseRejectsCorruption(t *testing.T) {
	_, err := settingsResponse([]model.UserSetting{{SettingKey: 2, SettingValue: "A"}, {SettingKey: 2, SettingValue: "B"}})
	if err == nil {
		t.Fatal("duplicate settings accepted")
	}
	r, err := settingsResponse([]model.UserSetting{{SettingKey: 2, SettingValue: "A"}})
	if err != nil || r.DefaultCourtLocation != nil {
		t.Fatal("historical settings invalid")
	}
	_, err = settingsResponse([]model.UserSetting{{SettingKey: 4, SettingValue: `{"address":"x","latitude":0,"longitude":0,"coordinateSystem":"gcj02"}`}})
	if err == nil {
		t.Fatal("orphan default location accepted")
	}
}
