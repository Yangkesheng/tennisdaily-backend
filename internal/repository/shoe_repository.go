package repository

import (
	"errors"
	"time"

	"tennisdaily-backend/internal/model"

	"gorm.io/gorm"
)

type ShoeRepository struct {
	db *gorm.DB
}

func NewShoeRepository(db *gorm.DB) *ShoeRepository {
	return &ShoeRepository{db: db}
}

func (r *ShoeRepository) BrandList() ([]model.ShoeBrand, error) {
	var items []model.ShoeBrand
	err := r.db.Order("name ASC").Find(&items).Error
	return items, err
}

func (r *ShoeRepository) SeriesList(query model.ShoeSeriesQuery) ([]model.ShoeSeries, error) {
	var items []model.ShoeSeries
	db := r.db.Model(&model.ShoeSeries{})
	if query.BrandID > 0 {
		db = db.Where("brand_id = ?", query.BrandID)
	}
	if query.Gender != nil {
		db = db.Where("gender = ?", *query.Gender)
	}
	err := db.Order("brand_id ASC, name ASC").Find(&items).Error
	return items, err
}

func (r *ShoeRepository) LibraryList(query model.ShoeLibraryQuery) ([]model.ShoeLibrary, error) {
	var items []model.ShoeLibrary
	db := r.db.Model(&model.ShoeLibrary{})
	if query.BrandID > 0 {
		db = db.Where("brand_id = ?", query.BrandID)
	}
	if query.SeriesID > 0 {
		db = db.Where("series_id = ?", query.SeriesID)
	}
	if query.Gender != nil {
		db = db.Where("gender = ?", *query.Gender)
	}
	err := db.Order("brand ASC, release_year DESC, id ASC").Find(&items).Error
	return items, err
}

func (r *ShoeRepository) LibraryBrandStats() ([]model.ShoeLibraryBrandStats, error) {
	var rows []model.ShoeLibraryBrandStats
	err := r.db.Model(&model.ShoeLibrary{}).
		Select("brand_id, brand, COUNT(*) AS count").
		Group("brand, brand_id").
		Order("count DESC").
		Scan(&rows).Error
	return rows, err
}

func (r *ShoeRepository) LibrarySeriesStats() ([]model.ShoeLibrarySeriesStats, error) {
	var rows []model.ShoeLibrarySeriesStats
	err := r.db.Model(&model.ShoeLibrary{}).
		Select("brand_id, brand, series_id, series, gender, COUNT(*) AS count").
		Group("brand, brand_id, series, series_id, gender").
		Order("count DESC").
		Scan(&rows).Error
	return rows, err
}

func (r *ShoeRepository) LibraryFindByID(id int64) (*model.ShoeLibrary, error) {
	var item model.ShoeLibrary
	err := r.db.Where("id = ?", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ShoeRepository) LibraryFindByIDs(ids []int64) (map[int64]model.ShoeLibrary, error) {
	items := make(map[int64]model.ShoeLibrary)
	if len(ids) == 0 {
		return items, nil
	}

	var libraries []model.ShoeLibrary
	if err := r.db.Where("id IN ?", ids).Find(&libraries).Error; err != nil {
		return nil, err
	}
	for _, item := range libraries {
		items[item.ID] = item
	}
	return items, nil
}

func (r *ShoeRepository) List(userID int64, includeRetired bool) ([]model.Shoe, error) {
	var shoes []model.Shoe
	query := r.db.Where("user_id = ? AND deleted_at IS NULL", userID)
	if !includeRetired {
		query = query.Where("status IN ?", []model.ShoeStatus{model.ShoeStatusPrimary, model.ShoeStatusActive})
	}
	err := query.Order("status ASC, created_at DESC").Find(&shoes).Error
	return shoes, err
}

func (r *ShoeRepository) Selectable(userID int64) ([]model.Shoe, error) {
	var shoes []model.Shoe
	err := r.db.Where("user_id = ? AND deleted_at IS NULL AND status IN ?", userID, []model.ShoeStatus{model.ShoeStatusPrimary, model.ShoeStatusActive}).
		Order("status ASC, created_at DESC").
		Find(&shoes).Error
	return shoes, err
}

func (r *ShoeRepository) Primary(userID int64) (*model.Shoe, error) {
	var shoe model.Shoe
	err := r.db.Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, model.ShoeStatusPrimary).
		Order("updated_at DESC, id DESC").
		First(&shoe).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &shoe, nil
}

// NamesByIDs 返回当前用户未删除球鞋的 id -> name 映射，用于打球记录展示时关联我的球鞋名称。
func (r *ShoeRepository) NamesByIDs(userID int64, ids []int64) (map[int64]string, error) {
	names := make(map[int64]string)
	if len(ids) == 0 {
		return names, nil
	}

	var shoes []model.Shoe
	err := r.db.Select("id", "name").
		Where("user_id = ? AND id IN ?", userID, ids).
		Find(&shoes).Error
	if err != nil {
		return nil, err
	}
	for _, shoe := range shoes {
		names[shoe.ID] = shoe.Name
	}
	return names, nil
}

// UsageStats 按 shoe_id 聚合当前用户未删除打球记录的上场次数与时长，用于球鞋使用统计。
func (r *ShoeRepository) UsageStats(userID int64, shoeIDs []int64) (map[int64]model.ShoeUsageStats, error) {
	usage := make(map[int64]model.ShoeUsageStats)
	if len(shoeIDs) == 0 {
		return usage, nil
	}

	type row struct {
		ShoeID  int64 `gorm:"column:shoe_id"`
		Count   int   `gorm:"column:count"`
		Minutes int   `gorm:"column:minutes"`
	}
	var rows []row
	err := r.db.Model(&model.TennisSession{}).
		Select("shoe_id, COUNT(id) AS count, COALESCE(SUM(duration_minutes), 0) AS minutes").
		Where("user_id = ? AND shoe_id IN ? AND deleted_at IS NULL", userID, shoeIDs).
		Group("shoe_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		usage[row.ShoeID] = model.ShoeUsageStats{
			Count:   row.Count,
			Minutes: row.Minutes,
			Hours:   row.Minutes / 60,
		}
	}
	return usage, nil
}

func (r *ShoeRepository) FindByID(userID, id int64) (*model.Shoe, error) {
	var shoe model.Shoe
	err := r.db.Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).First(&shoe).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &shoe, nil
}

func (r *ShoeRepository) Create(shoe *model.Shoe) error {
	return r.db.Create(shoe).Error
}

func (r *ShoeRepository) Update(shoe *model.Shoe) error {
	return r.db.Save(shoe).Error
}

func (r *ShoeRepository) SoftDelete(userID, id int64) (bool, error) {
	result := r.db.Model(&model.Shoe{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).
		Updates(map[string]interface{}{"deleted_at": time.Now(), "updated_at": time.Now()})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *ShoeRepository) SetPrimary(userID, id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Shoe{}).
			Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, model.ShoeStatusPrimary).
			Updates(map[string]interface{}{"status": model.ShoeStatusActive, "updated_at": time.Now()}).Error; err != nil {
			return err
		}

		result := tx.Model(&model.Shoe{}).
			Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).
			Updates(map[string]interface{}{"status": model.ShoeStatusPrimary, "updated_at": time.Now()})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (r *ShoeRepository) Retire(userID, id int64) (bool, error) {
	result := r.db.Model(&model.Shoe{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, id).
		Updates(map[string]interface{}{"status": model.ShoeStatusRetired, "updated_at": time.Now()})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *ShoeRepository) Stats(userID int64) (int, float64, error) {
	type shoeRow struct {
		Count int     `gorm:"column:count"`
		Cost  float64 `gorm:"column:cost"`
	}
	var shoeStats shoeRow
	if err := r.db.Model(&model.Shoe{}).
		Select("COUNT(*) AS count, COALESCE(SUM(purchase_price), 0) AS cost").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Scan(&shoeStats).Error; err != nil {
		return 0, 0, err
	}
	return shoeStats.Count, shoeStats.Cost, nil
}

func (r *ShoeRepository) SumPurchaseCostByMonth(userID int64, start, end time.Time) (float64, error) {
	return r.SumPurchaseCostByRange(userID, start, end)
}

func (r *ShoeRepository) SumPurchaseCostByRange(userID int64, start, end time.Time) (float64, error) {
	type row struct {
		Cost float64 `gorm:"column:cost"`
	}
	var result row
	err := r.db.Model(&model.Shoe{}).
		Select("COALESCE(SUM(purchase_price), 0) AS cost").
		Where("user_id = ? AND deleted_at IS NULL AND purchase_date >= ? AND purchase_date < ?", userID, start, end).
		Scan(&result).Error
	return result.Cost, err
}
