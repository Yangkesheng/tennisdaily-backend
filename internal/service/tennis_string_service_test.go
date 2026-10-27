package service

import (
	"tennisdaily-backend/internal/model"
	"testing"
)

func TestNormalizeStringQuery(t *testing.T) {
	q := normalizeStringQuery(model.StringLibraryQuery{Keyword: "  ALU   Power Rough  ", Brand: " Luxilon "})
	if q.Keyword != "ALU Power Rough" || q.Brand != "Luxilon" || q.Page != 1 || q.PageSize != 20 {
		t.Fatalf("unexpected query: %+v", q)
	}
}
