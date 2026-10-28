package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	driver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
)

// Uses only a newly created disposable database; never writes to the DSN's database.
func courtTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	raw := os.Getenv("COURT_TEST_MYSQL_DSN")
	if raw == "" {
		t.Skip("set COURT_TEST_MYSQL_DSN to run isolated MySQL integration tests")
	}
	cfg, err := driver.ParseDSN(raw)
	if err != nil {
		t.Fatal("invalid integration DSN")
	}
	if cfg.Net != "tcp" || !(strings.HasPrefix(cfg.Addr, "127.0.0.1:") || strings.HasPrefix(cfg.Addr, "localhost:")) {
		t.Fatal("integration tests only allow local MySQL")
	}
	cfg.DBName = ""
	cfg.ParseTime = true
	cfg.MultiStatements = true
	admin, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		t.Fatal("open admin connection failed")
	}
	dbName := fmt.Sprintf("court_test_%d_%d", os.Getpid(), time.Now().UnixNano())
	if _, err = admin.Exec("CREATE DATABASE `" + dbName + "` CHARACTER SET utf8mb4"); err != nil {
		t.Fatal("create isolated database failed")
	}
	t.Cleanup(func() {
		_, err := admin.Exec("DROP DATABASE `" + dbName + "`")
		if err != nil {
			t.Error("cleanup isolated database failed")
		}
		admin.Close()
	})
	cfg.DBName = dbName
	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal("open isolated database failed")
	}
	pool, _ := db.DB()
	t.Cleanup(func() { pool.Close() })
	// Build previous schema, including the nullable unique index from migration 025.
	if err = db.AutoMigrate(&model.User{}, &model.TennisSession{}, &model.UserSetting{}); err != nil {
		t.Fatal("create test tables failed")
	}
	for _, column := range []string{"court_address", "court_latitude", "court_longitude", "court_coordinate_system"} {
		if err = db.Migrator().DropColumn(&model.TennisSession{}, column); err != nil {
			t.Fatal("prepare old schema failed")
		}
	}
	if err = db.Exec("ALTER TABLE user_settings MODIFY setting_value VARCHAR(255) NOT NULL DEFAULT ''").Error; err != nil {
		t.Fatal("prepare old settings failed")
	}
	if err = db.Exec("ALTER TABLE user_settings ADD UNIQUE KEY uk_user_settings_user_key_deleted(user_id,setting_key,deleted_at)").Error; err != nil {
		t.Fatal("create old unique index failed")
	}
	old := model.TennisSession{UserID: 1, Date: time.Now(), CourtName: "历史场地"}
	if err = db.Omit("CourtAddress", "CourtLatitude", "CourtLongitude", "CourtCoordinateSystem").Create(&old).Error; err != nil {
		t.Fatal("seed historical row failed")
	}
	script, err := os.ReadFile("../../migrations/029_add_court_location.sql")
	if err != nil {
		t.Fatal(err)
	}

	// Ambiguous old settings must stop the migration before any location columns are added.
	if err = db.Create(&[]model.UserSetting{{UserID: 99, SettingKey: 2, SettingValue: "A"}, {UserID: 99, SettingKey: 2, SettingValue: "B"}}).Error; err != nil {
		t.Fatal("seed duplicate preflight failed")
	}
	if err = db.Exec(string(script)).Error; err == nil {
		t.Fatal("migration accepted duplicate active settings")
	}
	if db.Migrator().HasColumn(&model.TennisSession{}, "court_latitude") {
		t.Fatal("preflight failure partially changed schema")
	}
	if err = db.Where("user_id = ?", 99).Unscoped().Delete(&model.UserSetting{}).Error; err != nil {
		t.Fatal("reset duplicate fixture failed")
	}

	// Create the routine on one physical connection and call it on another.
	poolForMigration, _ := db.DB()
	firstConnection, err := poolForMigration.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	secondConnection, err := poolForMigration.Conn(context.Background())
	if err != nil {
		firstConnection.Close()
		t.Fatal(err)
	}
	parts := strings.SplitN(string(script), "\nCALL migrate_court_location_029();", 2)
	if len(parts) != 2 {
		t.Fatal("unexpected migration layout")
	}
	_, createErr := firstConnection.ExecContext(context.Background(), parts[0])
	var callErr error
	if createErr == nil {
		_, callErr = secondConnection.ExecContext(context.Background(), "CALL migrate_court_location_029();"+parts[1])
	}
	firstConnection.Close()
	secondConnection.Close()
	if createErr != nil || callErr != nil {
		t.Fatal("cross-connection migration failed")
	}

	if err = db.Exec(string(script)).Error; err != nil {
		t.Fatal("location migration is not repeatable")
	}
	var row model.TennisSession
	if err = db.First(&row, old.ID).Error; err != nil || row.CourtName != "历史场地" || row.CourtColumns.Location() != nil {
		t.Fatal("migration changed historical row")
	}
	if err = db.Delete(&row).Error; err != nil {
		t.Fatal("delete historical fixture failed")
	}
	for i := int64(1); i <= 2; i++ {
		if err = db.Create(&model.User{ID: i, OpenID: fmt.Sprintf("court-test-%d", i)}).Error; err != nil {
			t.Fatal("create test user failed")
		}
	}
	return db
}
func locationField(t *testing.T, lat float64) model.LocationField {
	t.Helper()
	return model.LocationField{Present: true, Value: &model.CourtLocation{Address: "测试地址", Latitude: lat, Longitude: 0, CoordinateSystem: "gcj02"}}
}
func TestCourtMySQLPersistenceAndRecommendations(t *testing.T) {
	db := courtTestDB(t)
	repo := repository.NewSessionRepository(db)
	service := NewSessionService(repo, repository.NewUserRepository(db), repository.NewRacketRepository(db), repository.NewShoeRepository(db), newTestSessionCategoryResolver(t), nil)
	name := "测试场地"
	req := model.CreateSessionRequest{Date: "2026-10-27 10:00", Category: model.SessionCategoryDaily, SubCategory: model.SessionSubCategoryDoubles, CourtName: &name, CourtLocation: locationField(t, 0)}
	created, err := service.Create(1, req)
	if err != nil || created.CourtLocation == nil || created.CourtLocation.Latitude != 0 {
		t.Fatal("create location failed")
	}
	loaded, err := service.FindByID(1, created.ID)
	if err != nil || loaded.CourtLocation == nil {
		t.Fatal("location did not round trip")
	}
	if _, err = service.FindByID(2, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatal("cross-user access allowed")
	}
	updated, err := service.Update(1, created.ID, model.UpdateSessionRequest{Date: req.Date, Category: req.Category, SubCategory: req.SubCategory, CourtName: &name, CourtLocation: req.CourtLocation, Cost: 25})
	if err != nil || updated.CourtLocation == nil || updated.Cost != 25 {
		t.Fatal("non-location update lost location")
	}
	cleared := ""
	_, err = service.Update(1, created.ID, model.UpdateSessionRequest{Date: req.Date, Category: req.Category, SubCategory: req.SubCategory, CourtName: &cleared, CourtLocation: model.LocationField{Present: true}})
	if err != nil {
		t.Fatal("clear failed")
	}
	loaded, err = service.FindByID(1, created.ID)
	if err != nil || loaded.CourtLocation != nil || loaded.CourtName != "" {
		t.Fatal("clear not persisted")
	}
	req.CourtLocation = model.LocationField{}
	if _, err = service.Create(1, req); !errors.Is(err, ErrInvalidRequest) {
		t.Fatal("old client allowed to save")
	}
	now := time.Date(2026, 10, 27, 12, 0, 0, 0, time.UTC)
	// More than one page of repeats must not hide other historical courts.
	for i := 0; i < 105; i++ {
		row := model.TennisSession{UserID: 1, Date: now.Add(-time.Duration(i) * time.Minute), CourtName: "重复场地", CourtColumns: model.ColumnsFromLocation(locationField(t, 1).Value)}
		if err = db.Create(&row).Error; err != nil {
			t.Fatal("create recommendation fixture failed")
		}
	}
	for i := 0; i < 4; i++ {
		row := model.TennisSession{UserID: 1, Date: now.Add(-time.Duration(110+i) * time.Minute), CourtName: "同名球场", CourtColumns: model.ColumnsFromLocation(locationField(t, float64(10+i)).Value)}
		if err = db.Create(&row).Error; err != nil {
			t.Fatal("seed distinct venue failed")
		}
	}
	// Empty latest session should not fall back for automatic fill.
	for _, row := range []model.TennisSession{
		{UserID: 1, Date: now.Add(time.Minute)},
		{UserID: 2, Date: now.Add(2 * time.Minute), CourtName: "他人场地"},
		{UserID: 1, Date: now.Add(3 * time.Minute), CourtName: "已删除场地", DeletedAt: gorm.DeletedAt{Time: now, Valid: true}},
	} {
		if err = db.Create(&row).Error; err != nil {
			t.Fatal("seed filter fixture failed")
		}
	}
	recommendations, err := service.CourtRecommendations(context.Background(), 1)
	if err != nil || recommendations.LatestCourt != nil || len(recommendations.RecentCourts) != 4 {
		t.Fatal("recommendation pagination/filtering failed")
	}
	if recommendations.RecentCourts[0].CourtName != "重复场地" || recommendations.RecentCourts[1].CourtLocation.Latitude != 10 {
		t.Fatal("recommendation ordering failed")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = service.CourtRecommendations(ctx, 1); !errors.Is(err, context.Canceled) {
		t.Fatal("request cancellation ignored")
	}
	// Same date: higher ID wins, independent of creation timestamp.
	for _, n := range []string{"same-time-first", "same-time-last"} {
		if err = db.Create(&model.TennisSession{UserID: 1, Date: now.Add(time.Hour), CourtName: n}).Error; err != nil {
			t.Fatal(err)
		}
	}
	recommendations, err = service.CourtRecommendations(context.Background(), 1)
	if err != nil || recommendations.LatestCourt.CourtName != "same-time-last" {
		t.Fatal("tie-breaker failed")
	}
}
func TestCourtMySQLSettingsAtomicity(t *testing.T) {
	db := courtTestDB(t)
	repo := repository.NewUserSettingRepository(db)
	settings := NewUserSettingsService(repo, repository.NewUserRepository(db), nil)
	ctx := context.Background()
	var wg sync.WaitGroup
	results := make(chan error, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			n := fmt.Sprintf("场地%d", i)
			r, err := settings.Update(ctx, 1, model.UpdateUserSettingsRequest{DefaultCourtName: &n, DefaultCourtLocation: locationField(t, float64(i))})
			if err == nil && (r.DefaultCourtName != n || r.DefaultCourtLocation.Latitude != float64(i)) {
				err = errors.New("mixed response snapshot")
			}
			results <- err
		}(i)
	}
	wg.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatal("concurrent settings save failed")
		}
	}
	rows, err := repo.ListByUser(ctx, 1)
	if err != nil || len(rows) != 2 {
		t.Fatal("duplicate active setting rows created")
	}
	prior, err := settings.Get(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	sentinel := errors.New("simulated second write failure")
	// Force failure when storing the second member of a default-court pair.
	callbackName := "court_test_fail_location"
	db.Callback().Update().Before("gorm:update").Register(callbackName, func(tx *gorm.DB) {
		if row, ok := tx.Statement.Dest.(*model.UserSetting); ok && row.SettingKey == model.UserSettingDefaultCourtLocation {
			tx.AddError(sentinel)
		}
	})
	n := "must rollback"
	_, err = settings.Update(ctx, 1, model.UpdateUserSettingsRequest{DefaultCourtName: &n, DefaultCourtLocation: locationField(t, 70)})
	db.Callback().Update().Remove(callbackName)
	if !errors.Is(err, sentinel) {
		t.Fatal("did not simulate location write failure")
	}
	after, err := settings.Get(ctx, 1)
	if err != nil || after.DefaultCourtName != prior.DefaultCourtName || *after.DefaultCourtLocation != *prior.DefaultCourtLocation {
		t.Fatal("partial default court committed")
	}
	// TEXT holds long addresses, and default changes do not alter recorded snapshots.
	l := locationField(t, 1)
	l.Value.Address = strings.Repeat("地", 512)
	r, err := settings.Update(ctx, 1, model.UpdateUserSettingsRequest{DefaultCourtName: &n, DefaultCourtLocation: l})
	if err != nil || r.DefaultCourtLocation.Address != l.Value.Address {
		t.Fatal("long location JSON truncated")
	}
	oldBody := `{"defaultCourtName":"old client"}`
	var req model.UpdateUserSettingsRequest
	json.Unmarshal([]byte(oldBody), &req)
	if _, err = settings.Update(ctx, 1, req); !errors.Is(err, ErrInvalidRequest) {
		t.Fatal("old settings write accepted")
	}
	zero := 0
	if _, err = settings.Update(ctx, 1, model.UpdateUserSettingsRequest{DefaultDurationMinutes: &zero}); err != nil {
		t.Fatal("unrelated setting requires court fields")
	}
	blank := ""
	if _, err = settings.Update(ctx, 1, model.UpdateUserSettingsRequest{DefaultCourtName: &blank, DefaultCourtLocation: model.LocationField{Present: true}}); err != nil {
		t.Fatal("default clear failed")
	}
	r, err = settings.Get(ctx, 1)
	if err != nil || r.DefaultCourtLocation != nil || r.DefaultCourtName != "" {
		t.Fatal("default clear not persisted")
	}
}

func TestCourtMySQLFreshInitialization(t *testing.T) {
	db := courtTestDB(t)
	if err := db.Migrator().DropTable(&model.UserSetting{}, &model.TennisSession{}, &model.User{}); err != nil {
		t.Fatal("prepare fresh database failed")
	}
	var databaseName string
	if err := db.Raw("SELECT DATABASE()").Scan(&databaseName).Error; err != nil {
		t.Fatal("read isolated database name failed")
	}
	content, err := os.ReadFile("../../migrations/init.sql")
	if err != nil {
		t.Fatal(err)
	}
	script := strings.ReplaceAll(string(content), "tennis_diary", databaseName)
	// Execute mysql-client DELIMITER syntax explicitly without ever selecting a business DB.
	for run := 0; run < 2; run++ {
		delimiter := ";"
		var statement strings.Builder
		for _, line := range strings.Split(script, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "DELIMITER ") {
				delimiter = strings.TrimSpace(strings.TrimPrefix(trimmed, "DELIMITER "))
				continue
			}
			statement.WriteString(line)
			statement.WriteByte('\n')
			if strings.HasSuffix(trimmed, delimiter) {
				sql := strings.TrimSpace(statement.String())
				sql = strings.TrimSuffix(sql, delimiter)
				if err = db.Exec(sql).Error; err != nil {
					t.Fatalf("init.sql failed on run %d: %v", run+1, err)
				}
				statement.Reset()
			}
		}
	}
	if !db.Migrator().HasColumn(&model.TennisSession{}, "court_latitude") {
		t.Fatal("fresh init lacks location")
	}
	var typ string
	if err = db.Raw("SELECT DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='user_settings' AND COLUMN_NAME='setting_value'").Scan(&typ).Error; err != nil || typ != "text" {
		t.Fatal("fresh init lacks TEXT settings")
	}
}

func TestCourtMySQLConsoleMigration(t *testing.T) {
	db := courtTestDB(t)
	script, err := os.ReadFile("../../migrations/029_court_location_console.sql")
	if err != nil {
		t.Fatal(err)
	}
	type statement struct {
		MigrationStatus string
		SQLToExecute    *string `gorm:"column:sql_to_execute"`
	}
	read := func() []statement {
		var rows []statement
		if err := db.Raw(string(script)).Scan(&rows).Error; err != nil {
			t.Fatal(err)
		}
		return rows
	}
	if rows := read(); len(rows) != 1 || rows[0].MigrationStatus != "DONE" {
		t.Fatal("current schema not recognized")
	}
	if err := db.Migrator().DropColumn(&model.TennisSession{}, "court_address"); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("ALTER TABLE user_settings MODIFY setting_value VARCHAR(255) NOT NULL DEFAULT ''").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&[]model.UserSetting{{UserID: 1, SettingKey: 2, SettingValue: "A"}, {UserID: 1, SettingKey: 2, SettingValue: "B"}}).Error; err != nil {
		t.Fatal(err)
	}
	for _, row := range read() {
		if row.SQLToExecute != nil || !strings.HasPrefix(row.MigrationStatus, "BLOCKED") {
			t.Fatal("duplicate settings did not block console migration")
		}
	}
	if err := db.Where("user_id = ?", 1).Unscoped().Delete(&model.UserSetting{}).Error; err != nil {
		t.Fatal(err)
	}
	rows := read()
	if len(rows) != 2 {
		t.Fatal("partial schema not recognized")
	}
	for _, row := range rows {
		if row.MigrationStatus != "READY" || row.SQLToExecute == nil {
			t.Fatal("missing executable statement")
		}
		if err := db.Exec(*row.SQLToExecute).Error; err != nil {
			t.Fatal(err)
		}
	}
	if rows := read(); len(rows) != 1 || rows[0].MigrationStatus != "DONE" {
		t.Fatal("migration not completed")
	}
}
