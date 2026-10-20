package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"tennisdaily-backend/internal/model"
	"testing"
)

func TestLibraryReviewScopeFiltersBeforeAggregation(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: "test:test@tcp(127.0.0.1:3306)/test", SkipInitializeWithVersion: true}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, status := range []int{0, 1} {
		var rows []model.RacketLibraryBrandStats
		stmt := db.Model(&model.RacketLibrary{}).Scopes(libraryReviewScope(&status)).Select("brand_id, COUNT(*) AS count").Group("brand_id").Find(&rows).Statement
		if !strings.Contains(stmt.SQL.String(), "WHERE is_reviewed = ? GROUP BY") || len(stmt.Vars) != 1 || stmt.Vars[0] != status {
			t.Fatalf("unexpected query: %s %v", stmt.SQL.String(), stmt.Vars)
		}
	}
	var rows []model.RacketLibrary
	stmt := db.Scopes(libraryReviewScope(nil)).Find(&rows).Statement
	if strings.Contains(stmt.SQL.String(), "is_reviewed =") {
		t.Fatal("admin all filter must include pending")
	}
}
