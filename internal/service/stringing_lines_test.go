package service

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"net/http/httptest"
	"strings"
	"tennisdaily-backend/internal/model"
	"testing"
)

type fakeStringCatalog struct {
	item *model.StringLibrary
	err  error
}

func (f fakeStringCatalog) CatalogString(int64) (*model.StringLibrary, error) { return f.item, f.err }
func ptrID(v int64) *int64                                                    { return &v }
func TestStructuredChoiceValidationAndSnapshot(t *testing.T) {
	catalog := fakeStringCatalog{item: &model.StringLibrary{ID: 1, Brand: "Brand", Model: "Current", Published: true, Variants: []model.StringVariant{{ID: 2, DiameterMm: 1.25}}}}
	req := model.StringingChoice{StringID: ptrID(1), VariantID: ptrID(2)}
	line, err := resolveStringingChoice(catalog, req, model.StringingMain, nil)
	if err != nil || line.Name != "Brand Current" || line.DiameterMm == nil || *line.DiameterMm != 1.25 {
		t.Fatalf("bad snapshot: %+v %v", line, err)
	}
	old := line
	old.Name = "Original name"
	kept, err := resolveStringingChoice(catalog, req, model.StringingMain, []model.StringingLine{old})
	if err != nil || kept.Name != old.Name {
		t.Fatal("renaming catalog must not change historical snapshot")
	}
	catalog.item.Published = false
	if _, err := resolveStringingChoice(catalog, req, model.StringingMain, nil); !errors.Is(err, ErrInvalidRequest) {
		t.Fatal("new reference to withdrawn item accepted")
	}
	if _, err := resolveStringingChoice(catalog, req, model.StringingMain, []model.StringingLine{old}); err != nil {
		t.Fatal("unchanged existing withdrawn reference rejected")
	}
	if _, err := resolveStringingChoice(catalog, req, model.StringingCross, []model.StringingLine{old}); err == nil {
		t.Fatal("withdrawn reference reused on a different side")
	}
	catalog.item.Published = true
	req.VariantID = ptrID(99)
	if _, err := resolveStringingChoice(catalog, req, model.StringingMain, nil); err == nil {
		t.Fatal("foreign variant accepted")
	}
	req.VariantID = nil
	line, err = resolveStringingChoice(catalog, req, model.StringingMain, nil)
	if err != nil || line.DiameterMm != nil {
		t.Fatal("unknown diameter must remain nil")
	}
	for _, bad := range []model.StringingChoice{{}, {VariantID: ptrID(2), CustomName: "Manual"}, {StringID: ptrID(0)}, {StringID: ptrID(1), CustomName: "Fake"}, {CustomName: strings.Repeat("线", 101)}} {
		if _, err := resolveStringingChoice(catalog, bad, model.StringingMain, nil); err == nil {
			t.Fatalf("invalid choice accepted: %+v", bad)
		}
	}
	line, err = resolveStringingChoice(catalog, model.StringingChoice{CustomName: "  Historical text  "}, model.StringingMain, nil)
	if err != nil || line.StringID != nil || line.Name != "Historical text" {
		t.Fatal("manual name must not be matched to catalog")
	}
}
func TestOldStringingWriteRejected(t *testing.T) {
	for _, body := range []string{`{"stringName":"Old client","stringDate":"2026-10-21 10:00"}`, `{"lines":null,"stringDate":"2026-10-21 10:00"}`} {
		req := httptest.NewRequest("POST", "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		var dto model.CreateStringingRecordRequest
		if err := binding.JSON.Bind(req, &dto); err == nil {
			t.Fatal("old write must fail validation")
		}
	}
}
func TestHistoricalRecordDisplay(t *testing.T) {
	r := model.RacketStringingRecord{StringName: "Original legacy name"}
	response := model.NewStringingRecordResponse(r)
	if response.Lines != nil || response.StringName != r.StringName {
		t.Fatal("legacy text changed")
	}
	r.Lines = []model.StringingLine{{Axis: model.StringingMain, Name: "A"}, {Axis: model.StringingCross, Name: "B"}}
	response = model.NewStringingRecordResponse(r)
	if response.StringName != "竖 A / 横 B" || response.Lines == nil {
		t.Fatal("hybrid snapshot not rendered")
	}
}

func TestManualDiameterSnapshot(t *testing.T) {
	for _, v := range []float64{0.5, 1.1, 1.25, 2} {
		line, err := resolveStringingChoice(fakeStringCatalog{}, model.StringingChoice{CustomName: "Manual", CustomDiameterMm: &v}, model.StringingMain, nil)
		if err != nil || line.DiameterMm == nil || *line.DiameterMm != v {
			t.Fatalf("diameter lost: %v %v", line, err)
		}
	}
	for _, v := range []float64{0, 0.49, 2.01, 1.255} {
		if _, err := resolveStringingChoice(fakeStringCatalog{}, model.StringingChoice{CustomName: "Manual", CustomDiameterMm: &v}, model.StringingMain, nil); err == nil {
			t.Fatalf("accepted invalid diameter %v", v)
		}
	}
}
