package repository

import (
	"errors"
	"strings"
	"time"

	"tennisdaily-backend/internal/model"

	"gorm.io/gorm"
)

type RacketRepository struct {
	db *gorm.DB
}

func NewRacketRepository(db *gorm.DB) *RacketRepository {
	return &RacketRepository{db: db}
}

func (r *RacketRepository) BrandList() ([]model.RacketBrand, error) {
	var items []model.RacketBrand
	err := r.db.Order("name ASC").Find(&items).Error
	return items, err
}

func (r *RacketRepository) SeriesList(query model.RacketSeriesQuery) ([]model.RacketSeries, error) {
	var items []model.RacketSeries
	db := r.db.Model(&model.RacketSeries{})
	if query.BrandID > 0 {
		db = db.Where("brand_id = ?", query.BrandID)
	}
	err := db.Order("brand_id ASC, name ASC").Find(&items).Error
	return items, err
}

func (r *RacketRepository) LibraryList(query model.RacketLibraryQuery) ([]model.RacketLibrary, error) {
	var items []model.RacketLibrary
	db := r.db.Model(&model.RacketLibrary{})
	if query.BrandID > 0 {
		db = db.Where("brand_id = ?", query.BrandID)
	}
	if query.SeriesID > 0 {
		db = db.Where("series_id = ?", query.SeriesID)
	}
	err := db.Order("brand ASC, model ASC, release_year DESC").Find(&items).Error
	return items, err
}

func (r *RacketRepository) LibraryBrandStats() ([]model.RacketLibraryBrandStats, error) {
	var rows []model.RacketLibraryBrandStats
	err := r.db.Model(&model.RacketLibrary{}).
		Select("brand_id, brand, COUNT(*) AS count").
		Group("brand, brand_id").
		Order("count DESC").
		Scan(&rows).Error
	return rows, err
}

func (r *RacketRepository) LibrarySeriesStats() ([]model.RacketLibrarySeriesStats, error) {
	var rows []model.RacketLibrarySeriesStats
	err := r.db.Model(&model.RacketLibrary{}).
		Select("brand_id, brand, series_id, series, COUNT(*) AS count").
		Group("brand, brand_id, series, series_id").
		Order("count DESC").
		Scan(&rows).Error
	return rows, err
}

func (r *RacketRepository) LibraryFindByID(id int64) (*model.RacketLibrary, error) {
	var item model.RacketLibrary
	err := r.db.Where("id = ?", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RacketRepository) LibraryFindByIDs(ids []int64) (map[int64]model.RacketLibrary, error) {
	items := make(map[int64]model.RacketLibrary)
	if len(ids) == 0 {
		return items, nil
	}

	var libraries []model.RacketLibrary
	if err := r.db.Where("id IN ?", ids).Find(&libraries).Error; err != nil {
		return nil, err
	}
	for _, item := range libraries {
		items[item.ID] = item
	}
	return items, nil
}

func (r *RacketRepository) List(userID int64, includeRetired bool) ([]model.Racket, error) {
	var rackets []model.Racket
	query := r.db.Where("user_id = ? AND deleted_at IS NULL", userID)
	if !includeRetired {
		query = query.Where("status IN ?", []model.RacketStatus{model.RacketStatusPrimary, model.RacketStatusActive})
	}
	err := query.Order("status ASC, created_at DESC").Find(&rackets).Error
	return rackets, err
}

func (r *RacketRepository) Selectable(userID int64) ([]model.Racket, error) {
	var rackets []model.Racket
	err := r.db.Where("user_id = ? AND deleted_at IS NULL AND status IN ?", userID, []model.RacketStatus{model.RacketStatusPrimary, model.RacketStatusActive}).
		Order("status ASC, created_at DESC").
		Find(&rackets).Error
	return rackets, err
}

func (r *RacketRepository) Primary(userID int64) (*model.Racket, error) {
	var racket model.Racket
	err := r.db.Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, model.RacketStatusPrimary).
		Order("updated_at DESC, id DESC").
		First(&racket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &racket, nil
}

func (r *RacketRepository) Stats(userID int64) (int, float64, float64, error) {
	type racketRow struct {
		Count int     `gorm:"column:count"`
		Cost  float64 `gorm:"column:cost"`
	}
	var racketStats racketRow
	if err := r.db.Model(&model.Racket{}).
		Select("COUNT(*) AS count, COALESCE(SUM(purchase_price), 0) AS cost").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Scan(&racketStats).Error; err != nil {
		return 0, 0, 0, err
	}

	type stringingRow struct {
		Cost float64 `gorm:"column:cost"`
	}
	var stringingStats stringingRow
	if err := r.db.Table("racket_stringing_record AS rsr").
		Select("COALESCE(SUM(rsr.cost), 0) AS cost").
		Joins("JOIN racket AS r ON r.id = rsr.racket_id AND r.user_id = rsr.user_id").
		Where("rsr.user_id = ? AND rsr.deleted_at IS NULL AND r.deleted_at IS NULL", userID).
		Scan(&stringingStats).Error; err != nil {
		return 0, 0, 0, err
	}

	return racketStats.Count, racketStats.Cost, stringingStats.Cost, nil
}

func (r *RacketRepository) SumPurchaseCostByMonth(userID int64, start, end time.Time) (float64, error) {
	return r.SumPurchaseCostByRange(userID, start, end)
}

func (r *RacketRepository) SumStringingCostByMonth(userID int64, start, end time.Time) (float64, error) {
	return r.SumStringingCostByRange(userID, start, end)
}

func (r *RacketRepository) SumPurchaseCostByRange(userID int64, start, end time.Time) (float64, error) {
	type row struct {
		Cost float64 `gorm:"column:cost"`
	}
	var result row
	err := r.db.Model(&model.Racket{}).
		Select("COALESCE(SUM(purchase_price), 0) AS cost").
		Where("user_id = ? AND deleted_at IS NULL AND purchase_date >= ? AND purchase_date < ?", userID, start, end).
		Scan(&result).Error
	return result.Cost, err
}

func (r *RacketRepository) SumStringingCostByRange(userID int64, start, end time.Time) (float64, error) {
	type row struct {
		Cost float64 `gorm:"column:cost"`
	}
	var result row
	err := r.db.Model(&model.RacketStringingRecord{}).
		Select("COALESCE(SUM(cost), 0) AS cost").
		Where("user_id = ? AND deleted_at IS NULL AND string_date >= ? AND string_date < ?", userID, start, end).
		Scan(&result).Error
	return result.Cost, err
}

func (r *RacketRepository) Create(racket *model.Racket) error {
	return r.db.Create(racket).Error
}

func (r *RacketRepository) FindByID(userID, id int64) (*model.Racket, error) {
	var racket model.Racket
	err := r.db.Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).First(&racket).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &racket, nil
}

func (r *RacketRepository) Update(racket *model.Racket) error {
	return r.db.Save(racket).Error
}

func (r *RacketRepository) SoftDelete(userID, id int64) (bool, error) {
	result := r.db.Model(&model.Racket{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).
		Updates(map[string]interface{}{"deleted_at": time.Now(), "updated_at": time.Now()})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *RacketRepository) SetPrimary(userID, id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Racket{}).
			Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, model.RacketStatusPrimary).
			Updates(map[string]interface{}{"status": model.RacketStatusActive, "updated_at": time.Now()}).Error; err != nil {
			return err
		}

		result := tx.Model(&model.Racket{}).
			Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).
			Updates(map[string]interface{}{"status": model.RacketStatusPrimary, "updated_at": time.Now()})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (r *RacketRepository) Retire(userID, id int64) (bool, error) {
	result := r.db.Model(&model.Racket{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).
		Updates(map[string]interface{}{"status": model.RacketStatusRetired, "updated_at": time.Now()})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *RacketRepository) CreateStringingRecord(record *model.RacketStringingRecord) error {
	return r.db.Create(record).Error
}

func (r *RacketRepository) FindStringingRecordByID(userID, racketID, recordID int64) (*model.RacketStringingRecord, error) {
	var record model.RacketStringingRecord
	err := r.db.Where("user_id = ? AND racket_id = ? AND id = ? AND deleted_at IS NULL", userID, racketID, recordID).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *RacketRepository) UpdateStringingRecord(record *model.RacketStringingRecord) error {
	return r.db.Save(record).Error
}

func (r *RacketRepository) SoftDeleteStringingRecord(userID, racketID, recordID int64) (bool, error) {
	result := r.db.Model(&model.RacketStringingRecord{}).
		Where("user_id = ? AND racket_id = ? AND id = ? AND deleted_at IS NULL", userID, racketID, recordID).
		Updates(map[string]interface{}{"deleted_at": time.Now(), "updated_at": time.Now()})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *RacketRepository) ListStringingRecords(userID, racketID int64) ([]model.RacketStringingRecord, error) {
	var records []model.RacketStringingRecord
	err := r.db.Where("user_id = ? AND racket_id = ? AND deleted_at IS NULL", userID, racketID).
		Order("string_date DESC, created_at DESC").
		Find(&records).Error
	return records, err
}

func (r *RacketRepository) LatestStringingRecords(userID int64, racketIDs []int64) (map[int64]model.RacketStringingRecord, error) {
	latest := make(map[int64]model.RacketStringingRecord)
	if len(racketIDs) == 0 {
		return latest, nil
	}

	var records []model.RacketStringingRecord
	err := r.db.Where("user_id = ? AND racket_id IN ? AND deleted_at IS NULL", userID, racketIDs).
		Order("racket_id ASC, string_date DESC, id DESC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if _, ok := latest[record.RacketID]; !ok {
			latest[record.RacketID] = record
		}
	}
	return latest, nil
}

func (r *RacketRepository) UsageStats(userID int64, racketIDs []int64) (map[int64]model.RacketUsageStats, error) {
	usage := make(map[int64]model.RacketUsageStats)
	if len(racketIDs) == 0 {
		return usage, nil
	}

	type row struct {
		RacketID int64 `gorm:"column:racket_id"`
		Count    int   `gorm:"column:count"`
		Minutes  int   `gorm:"column:minutes"`
	}
	var rows []row
	err := r.db.Model(&model.TennisSession{}).
		Select("racket_id, COUNT(id) AS count, COALESCE(SUM(duration_minutes), 0) AS minutes").
		Where("user_id = ? AND racket_id IN ? AND deleted_at IS NULL", userID, racketIDs).
		Group("racket_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		usage[row.RacketID] = model.RacketUsageStats{
			Count:   row.Count,
			Minutes: row.Minutes,
			Hours:   row.Minutes / 60,
		}
	}
	return usage, nil
}

func (r *RacketRepository) UsageStatsSince(userID int64, latestStringing map[int64]model.RacketStringingRecord) (map[int64]model.RacketUsageStats, error) {
	usage := make(map[int64]model.RacketUsageStats)
	if len(latestStringing) == 0 {
		return usage, nil
	}

	conditions := make([]string, 0, len(latestStringing))
	args := make([]interface{}, 0, len(latestStringing)*2+1)
	args = append(args, userID)
	for racketID, record := range latestStringing {
		conditions = append(conditions, "(racket_id = ? AND date >= ?)")
		args = append(args, racketID, record.StringDate)
	}

	type row struct {
		RacketID int64 `gorm:"column:racket_id"`
		Count    int   `gorm:"column:count"`
		Minutes  int   `gorm:"column:minutes"`
	}
	var rows []row
	err := r.db.Model(&model.TennisSession{}).
		Select("racket_id, COUNT(id) AS count, COALESCE(SUM(duration_minutes), 0) AS minutes").
		Where("user_id = ? AND deleted_at IS NULL AND ("+strings.Join(conditions, " OR ")+")", args...).
		Group("racket_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		usage[row.RacketID] = model.RacketUsageStats{
			Count:   row.Count,
			Minutes: row.Minutes,
			Hours:   row.Minutes / 60,
		}
	}
	return usage, nil
}
