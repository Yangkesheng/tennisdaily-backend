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
	err := db.Order("brand ASC, release_year DESC, model ASC, weight ASC").Find(&items).Error
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

// ---- 管理员维护球拍库 ----

func (r *RacketRepository) BrandFindByName(name string) (*model.RacketBrand, error) {
	var item model.RacketBrand
	err := r.db.Where("name = ?", name).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RacketRepository) CreateBrand(brand *model.RacketBrand) error {
	return r.db.Create(brand).Error
}

func (r *RacketRepository) SeriesFindUnique(brandID int64, name string) (*model.RacketSeries, error) {
	var item model.RacketSeries
	err := r.db.Where("brand_id = ? AND name = ?", brandID, name).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RacketRepository) CreateSeries(series *model.RacketSeries) error {
	return r.db.Create(series).Error
}

// LibraryFindDuplicate 按 品牌+系列+型号+年份 查重，避免管理员重复录入。
func (r *RacketRepository) LibraryFindDuplicate(brandID, seriesID int64, modelName string, releaseYear int) (bool, error) {
	var count int64
	err := r.db.Model(&model.RacketLibrary{}).
		Where("brand_id = ? AND series_id = ? AND model = ? AND release_year = ?",
			brandID, seriesID, modelName, releaseYear).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RacketRepository) CreateLibraryItems(items []model.RacketLibrary) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i := range items {
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// LibraryImagesPending 返回待上传对象存储的球拍库图片：
// image_url 非空，且尚未写入上传记录表（记录表是幂等依据，不依赖 file_id）。
func (r *RacketRepository) LibraryImagesPending() ([]model.RacketLibrary, error) {
	var items []model.RacketLibrary
	err := r.db.Model(&model.RacketLibrary{}).
		Where("image_url IS NOT NULL AND image_url <> ''").
		Where("NOT EXISTS (SELECT 1 FROM racket_library_image_upload u WHERE u.racket_library_id = racket_library.id)").
		Order("id ASC").
		Find(&items).Error
	return items, err
}

// ApplyLibraryImageUpload 在事务中更新 file_id（不动 image_url）并写入上传记录。
func (r *RacketRepository) ApplyLibraryImageUpload(libraryID int64, fileID, objectKey, sourceURL string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.RacketLibrary{}).
			Where("id = ?", libraryID).
			Update("file_id", fileID).Error; err != nil {
			return err
		}

		record := model.RacketLibraryImageUpload{
			RacketLibraryID: libraryID,
			FileID:          fileID,
			ObjectKey:       objectKey,
			SourceURL:       sourceURL,
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
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

// NamesByIDs 返回当前用户未删除球拍的 id -> name 映射，用于打球记录展示时关联我的球拍名称。
func (r *RacketRepository) NamesByIDs(userID int64, ids []int64) (map[int64]string, error) {
	names := make(map[int64]string)
	if len(ids) == 0 {
		return names, nil
	}

	var rackets []model.Racket
	err := r.db.Select("id", "name").
		Where("user_id = ? AND id IN ?", userID, ids).
		Find(&rackets).Error
	if err != nil {
		return nil, err
	}
	for _, racket := range rackets {
		names[racket.ID] = racket.Name
	}
	return names, nil
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
