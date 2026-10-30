package service

import (
	"strings"
	"tennisdaily-backend/internal/model"
	"testing"
)

func TestNormalizeStringingStore(t *testing.T) {
	valid := &model.CourtLocation{Address: " 地址 ", Latitude: 34.123456789, Longitude: 108.1, CoordinateSystem: "gcj02"}
	for _, tc := range []struct {
		name  string
		field model.LocationField
	}{
		{"", model.LocationField{}},
		{"手填", model.LocationField{Present: true}},
		{"", model.LocationField{Present: true, Value: valid}},
		{strings.Repeat("店", 101), model.LocationField{Present: true, Value: valid}},
		{"店", model.LocationField{Present: true, Value: &model.CourtLocation{Latitude: 91, CoordinateSystem: "gcj02"}}},
	} {
		if _, err := normalizeStringingStore(tc.name, tc.field); err == nil {
			t.Errorf("accepted invalid store %q", tc.name)
		}
	}
	if loc, err := normalizeStringingStore("", model.LocationField{Present: true}); err != nil || loc != nil {
		t.Fatal("empty location must be accepted")
	}
	loc, err := normalizeStringingStore("店", model.LocationField{Present: true, Value: valid})
	if err != nil || loc.Address != "地址" || loc.Latitude != 34.1234568 {
		t.Fatalf("normalization failed: %+v %v", loc, err)
	}
	if valid.Address != " 地址 " || valid.Latitude != 34.123456789 {
		t.Fatal("must not mutate input")
	}
}
