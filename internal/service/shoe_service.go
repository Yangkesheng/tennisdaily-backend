package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"tennisdaily-backend/internal/config"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"

	"gorm.io/gorm"
)

type ShoeService struct {
	repo            *repository.ShoeRepository
	userRepo        *repository.UserRepository
	contentSecurity *ContentSecurityService
	shoeWear        *config.ShoeWearResolver
	loc             *time.Location
}

func NewShoeService(repo *repository.ShoeRepository, userRepo *repository.UserRepository, contentSecurity *ContentSecurityService, shoeWear *config.ShoeWearResolver) *ShoeService {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	return &ShoeService{repo: repo, userRepo: userRepo, contentSecurity: contentSecurity, shoeWear: shoeWear, loc: loc}
}

func (s *ShoeService) Brands() ([]model.ShoeBrandResponse, error) {
	items, err := s.repo.BrandList()
	if err != nil {
		return nil, err
	}
	responses := make([]model.ShoeBrandResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, model.NewShoeBrandResponse(item))
	}
	return responses, nil
}

func (s *ShoeService) Series(query model.ShoeSeriesQuery) ([]model.ShoeSeriesResponse, error) {
	items, err := s.repo.SeriesList(query)
	if err != nil {
		return nil, err
	}
	responses := make([]model.ShoeSeriesResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, model.NewShoeSeriesResponse(item))
	}
	return responses, nil
}

func (s *ShoeService) Library(query model.ShoeLibraryQuery) ([]model.ShoeLibraryBrandGroupResponse, error) {
	items, err := s.repo.LibraryList(query)
	if err != nil {
		return nil, err
	}

	groups := make([]model.ShoeLibraryBrandGroupResponse, 0)
	brandIndex := make(map[string]int)
	for _, item := range items {
		idx, ok := brandIndex[item.Brand]
		if !ok {
			idx = len(groups)
			brandIndex[item.Brand] = idx
			groups = append(groups, model.ShoeLibraryBrandGroupResponse{BrandID: item.BrandID, Brand: item.Brand, Items: []model.ShoeLibraryItemResponse{}})
		}
		groups[idx].Items = append(groups[idx].Items, model.NewShoeLibraryItemResponse(item))
	}
	return groups, nil
}

func (s *ShoeService) LibraryStats() ([]model.ShoeLibraryBrandStatsResponse, error) {
	brands, err := s.repo.LibraryBrandStats()
	if err != nil {
		return nil, err
	}
	series, err := s.repo.LibrarySeriesStats()
	if err != nil {
		return nil, err
	}

	responses := make([]model.ShoeLibraryBrandStatsResponse, 0, len(brands))
	type brandStatsKey struct {
		BrandID int64
		Brand   string
	}
	brandIndex := make(map[brandStatsKey]int, len(brands))
	for _, brand := range brands {
		brandIndex[brandStatsKey{BrandID: brand.BrandID, Brand: brand.Brand}] = len(responses)
		responses = append(responses, model.ShoeLibraryBrandStatsResponse{
			BrandID: brand.BrandID,
			Brand:   brand.Brand,
			Count:   brand.Count,
			Series:  map[string][]model.ShoeLibrarySeriesStatsResponse{},
		})
	}
	for _, item := range series {
		idx, ok := brandIndex[brandStatsKey{BrandID: item.BrandID, Brand: item.Brand}]
		if !ok {
			continue
		}
		genderKey := strconv.Itoa(item.Gender)
		responses[idx].Series[genderKey] = append(responses[idx].Series[genderKey], model.ShoeLibrarySeriesStatsResponse{
			SeriesID: item.SeriesID,
			Series:   item.Series,
			Count:    item.Count,
		})
	}
	return responses, nil
}

func (s *ShoeService) Stats(userID int64) (model.ShoeStatsResponse, error) {
	count, cost, err := s.repo.Stats(userID)
	if err != nil {
		return model.ShoeStatsResponse{}, err
	}
	return model.ShoeStatsResponse{
		ShoeCount:     count,
		ShoeCost:      cost,
		TotalCost:     cost,
		ShoeCostText:  fmt.Sprintf("%.2f", cost),
		TotalCostText: fmt.Sprintf("%.2f", cost),
	}, nil
}

func (s *ShoeService) List(userID int64, includeRetired bool) ([]model.ShoeResponse, error) {
	shoes, err := s.repo.List(userID, includeRetired)
	if err != nil {
		return nil, err
	}
	return s.enrichShoes(userID, shoes)
}

func (s *ShoeService) Selectable(userID int64) ([]model.ShoeResponse, error) {
	shoes, err := s.repo.Selectable(userID)
	if err != nil {
		return nil, err
	}
	return s.enrichShoes(userID, shoes)
}

func (s *ShoeService) MyShoesForSession(userID int64) ([]model.ShoeResponse, error) {
	return s.Selectable(userID)
}

func (s *ShoeService) Primary(userID int64) (*model.ShoeResponse, error) {
	shoe, err := s.repo.Primary(userID)
	if err != nil {
		return nil, err
	}
	if shoe == nil {
		return nil, nil
	}
	responses, err := s.enrichShoes(userID, []model.Shoe{*shoe})
	if err != nil {
		return nil, err
	}
	return &responses[0], nil
}

func (s *ShoeService) Create(userID int64, req model.CreateShoeRequest) (model.ShoeResponse, error) {
	library, err := s.applyLibraryDefaults(&req.LibraryID, &req.Name, &req.Brand, &req.Model, &req.Colorway)
	if err != nil {
		return model.ShoeResponse{}, err
	}
	if req.Name == "" {
		return model.ShoeResponse{}, ErrInvalidRequest
	}
	if err := s.checkUserInputTexts(userID, req.Name, req.Brand, req.Model, req.Size, req.Colorway); err != nil {
		return model.ShoeResponse{}, err
	}

	purchaseDate, err := s.parseOptionalDate(req.PurchaseDate)
	if err != nil {
		return model.ShoeResponse{}, ErrInvalidRequest
	}

	status := model.ShoeStatusActive
	if req.Status != 0 {
		if req.Status != model.ShoeStatusPrimary && req.Status != model.ShoeStatusActive {
			return model.ShoeResponse{}, ErrInvalidRequest
		}
		status = req.Status
	}

	shoe := model.Shoe{
		UserID:        userID,
		LibraryID:     req.LibraryID,
		Name:          req.Name,
		Brand:         req.Brand,
		Model:         req.Model,
		Status:        status,
		Size:          req.Size,
		Colorway:      req.Colorway,
		PurchaseDate:  purchaseDate,
		PurchasePrice: req.PurchasePrice,
	}
	if library != nil {
		shoe.Gender = library.Gender
		shoe.ReleaseYear = library.ReleaseYear
		shoe.FileID = library.FileID
	}
	if err := s.repo.Create(&shoe); err != nil {
		return model.ShoeResponse{}, err
	}
	if status == model.ShoeStatusPrimary {
		if err := s.repo.SetPrimary(userID, shoe.ID); err != nil {
			return model.ShoeResponse{}, err
		}
		shoe.Status = model.ShoeStatusPrimary
	}

	responses, err := s.enrichShoes(userID, []model.Shoe{shoe})
	if err != nil {
		return model.ShoeResponse{}, err
	}
	return responses[0], nil
}

func (s *ShoeService) Detail(userID, id int64) (model.ShoeResponse, error) {
	shoe, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.ShoeResponse{}, err
	}
	if shoe == nil {
		return model.ShoeResponse{}, ErrNotFound
	}

	responses, err := s.enrichShoes(userID, []model.Shoe{*shoe})
	if err != nil {
		return model.ShoeResponse{}, err
	}
	return responses[0], nil
}

func (s *ShoeService) Update(userID, id int64, req model.UpdateShoeRequest) (model.ShoeResponse, error) {
	if _, err := s.applyLibraryDefaults(&req.LibraryID, &req.Name, &req.Brand, &req.Model, &req.Colorway); err != nil {
		return model.ShoeResponse{}, err
	}
	if req.Name == "" {
		return model.ShoeResponse{}, ErrInvalidRequest
	}
	if err := s.checkUserInputTexts(userID, req.Name, req.Brand, req.Model, req.Size, req.Colorway); err != nil {
		return model.ShoeResponse{}, err
	}

	shoe, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.ShoeResponse{}, err
	}
	if shoe == nil {
		return model.ShoeResponse{}, ErrNotFound
	}

	purchaseDate, err := s.parseOptionalDate(req.PurchaseDate)
	if err != nil {
		return model.ShoeResponse{}, ErrInvalidRequest
	}

	shoe.LibraryID = req.LibraryID
	shoe.Name = req.Name
	shoe.Brand = req.Brand
	shoe.Model = req.Model
	shoe.Size = req.Size
	shoe.Colorway = req.Colorway
	shoe.PurchaseDate = purchaseDate
	shoe.PurchasePrice = req.PurchasePrice
	if req.Status >= model.ShoeStatusPrimary && req.Status <= model.ShoeStatusRetired {
		shoe.Status = req.Status
	}

	if shoe.Status == model.ShoeStatusPrimary {
		if err := s.repo.Update(shoe); err != nil {
			return model.ShoeResponse{}, err
		}
		if err := s.repo.SetPrimary(userID, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return model.ShoeResponse{}, ErrNotFound
			}
			return model.ShoeResponse{}, err
		}
		shoe.Status = model.ShoeStatusPrimary
	} else if err := s.repo.Update(shoe); err != nil {
		return model.ShoeResponse{}, err
	}

	updated, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.ShoeResponse{}, err
	}
	responses, err := s.enrichShoes(userID, []model.Shoe{*updated})
	if err != nil {
		return model.ShoeResponse{}, err
	}
	return responses[0], nil
}

func (s *ShoeService) SetPrimary(userID, id int64) (model.ShoeResponse, error) {
	if err := s.repo.SetPrimary(userID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ShoeResponse{}, ErrNotFound
		}
		return model.ShoeResponse{}, err
	}
	shoe, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.ShoeResponse{}, err
	}
	responses, err := s.enrichShoes(userID, []model.Shoe{*shoe})
	if err != nil {
		return model.ShoeResponse{}, err
	}
	return responses[0], nil
}

func (s *ShoeService) Retire(userID, id int64) (model.ShoeResponse, error) {
	ok, err := s.repo.Retire(userID, id)
	if err != nil {
		return model.ShoeResponse{}, err
	}
	if !ok {
		return model.ShoeResponse{}, ErrNotFound
	}
	shoe, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.ShoeResponse{}, err
	}
	responses, err := s.enrichShoes(userID, []model.Shoe{*shoe})
	if err != nil {
		return model.ShoeResponse{}, err
	}
	return responses[0], nil
}

func (s *ShoeService) Delete(userID, id int64) error {
	deleted, err := s.repo.SoftDelete(userID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func (s *ShoeService) checkUserInputTexts(userID int64, texts ...string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUnauthorized
	}
	return s.contentSecurity.CheckTexts(user.OpenID, texts...)
}

func (s *ShoeService) enrichShoes(userID int64, shoes []model.Shoe) ([]model.ShoeResponse, error) {
	ids := make([]int64, 0, len(shoes))
	libraryIDs := make([]int64, 0, len(shoes))
	libraryIDSet := make(map[int64]struct{})
	for _, shoe := range shoes {
		ids = append(ids, shoe.ID)
		if shoe.LibraryID > 0 {
			if _, ok := libraryIDSet[shoe.LibraryID]; !ok {
				libraryIDSet[shoe.LibraryID] = struct{}{}
				libraryIDs = append(libraryIDs, shoe.LibraryID)
			}
		}
	}
	libraries, err := s.repo.LibraryFindByIDs(libraryIDs)
	if err != nil {
		return nil, err
	}
	usage, err := s.repo.UsageStats(userID, ids)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ShoeResponse, 0, len(shoes))
	for _, shoe := range shoes {
		if library, ok := libraries[shoe.LibraryID]; ok {
			shoe.Gender = library.Gender
			shoe.ReleaseYear = library.ReleaseYear
			shoe.FileID = library.FileID
		}
		stats := usage[shoe.ID]
		shoe.UsageCount = stats.Count
		shoe.UsageMinutes = stats.Minutes
		shoe.UsageHours = stats.Hours
		shoe.TotalMinutes = stats.Minutes
		shoe.TotalHours = stats.Hours
		shoe.Wear = s.calculateShoeWear(shoe.PurchaseDate, stats.Minutes)
		responses = append(responses, model.NewShoeResponse(shoe))
	}
	return responses, nil
}

func (s *ShoeService) calculateShoeWear(purchaseDate *time.Time, usageMinutes int) *model.ShoeWearResponse {
	if s.shoeWear == nil {
		return nil
	}
	var purchase time.Time
	if purchaseDate != nil {
		purchase = *purchaseDate
	}
	wear := calculateShoeWearAt(time.Now().In(s.loc), purchase, usageMinutes, s.loc, s.shoeWear)
	return &wear
}

// calculateShoeWearAt 参考球线健康度估算球鞋磨损度：
// effectiveWear = 上场小时数 + 购买后天数 * 每日静置老化系数，
// score = 100 * (1 - effectiveWear / 标准寿命)，剩余寿命按小时取整。
func calculateShoeWearAt(now, purchaseDate time.Time, usageMinutes int, loc *time.Location, wear *config.ShoeWearResolver) model.ShoeWearResponse {
	now = now.In(loc)
	daysSincePurchase := 0
	if !purchaseDate.IsZero() && now.After(purchaseDate) {
		startDay := time.Date(purchaseDate.Year(), purchaseDate.Month(), purchaseDate.Day(), 0, 0, 0, 0, loc)
		nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		daysSincePurchase = int(nowDay.Sub(startDay).Hours() / 24)
	}

	hoursPlayed := float64(usageMinutes) / 60
	effectiveWear := hoursPlayed + float64(daysSincePurchase)*wear.RestWearPerDay()
	standardLifeHours := wear.StandardLifeHours()
	score := math.Max(0, 100*(1-effectiveWear/standardLifeHours))
	remainingHours := int(math.Round(math.Max(0, standardLifeHours-effectiveWear)))

	state := wear.State(score)
	return model.ShoeWearResponse{
		State:          state.Key,
		Display:        wear.RenderDisplay(state.Display, remainingHours),
		Score:          score,
		RemainingHours: remainingHours,
	}
}

// applyLibraryDefaults 从球鞋库补全品牌、型号、名称与配色；
// libraryID 无效时返回 ErrInvalidRequest。
func (s *ShoeService) applyLibraryDefaults(libraryID *int64, name *string, brand *string, shoeModel *string, colorway *string) (*model.ShoeLibrary, error) {
	if libraryID == nil || *libraryID == 0 {
		return nil, nil
	}
	item, err := s.repo.LibraryFindByID(*libraryID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrInvalidRequest
	}
	if *brand == "" {
		*brand = item.Brand
	}
	if *shoeModel == "" {
		*shoeModel = item.Model
	}
	if *name == "" {
		*name = item.Brand + " " + item.Model
	}
	if *colorway == "" {
		*colorway = item.Colorway
	}
	return item, nil
}

func (s *ShoeService) parseOptionalDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", value, s.loc)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// ---- 管理员维护球鞋库 ----

// CreateBrand 管理员新增球鞋品牌；slug 未传时由品牌名自动生成。
func (s *ShoeService) CreateBrand(req model.CreateShoeBrandRequest) (model.ShoeBrandResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return model.ShoeBrandResponse{}, NewInvalidRequestError("品牌名称不能为空")
	}
	existing, err := s.repo.BrandFindByName(name)
	if err != nil {
		return model.ShoeBrandResponse{}, err
	}
	if existing != nil {
		return model.ShoeBrandResponse{}, NewInvalidRequestError("品牌已存在")
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = deriveBrandSlug(name)
	}

	brand := model.ShoeBrand{Name: name, Slug: slug, FileID: req.FileID}
	if err := s.repo.CreateBrand(&brand); err != nil {
		return model.ShoeBrandResponse{}, err
	}
	return model.NewShoeBrandResponse(brand), nil
}

// CreateSeries 管理员新增球鞋系列；同一品牌、性别下系列名唯一。
func (s *ShoeService) CreateSeries(req model.CreateShoeSeriesRequest) (model.ShoeSeriesResponse, error) {
	if req.BrandID <= 0 {
		return model.ShoeSeriesResponse{}, NewInvalidRequestError("品牌不能为空")
	}
	brand, err := s.repo.BrandFindByID(req.BrandID)
	if err != nil {
		return model.ShoeSeriesResponse{}, err
	}
	if brand == nil {
		return model.ShoeSeriesResponse{}, NewInvalidRequestError("品牌不存在")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return model.ShoeSeriesResponse{}, NewInvalidRequestError("系列名称不能为空")
	}
	if req.Gender < 0 || req.Gender > 3 {
		return model.ShoeSeriesResponse{}, NewInvalidRequestError("性别取值不合法")
	}

	existing, err := s.repo.SeriesFindUnique(req.BrandID, req.Gender, name)
	if err != nil {
		return model.ShoeSeriesResponse{}, err
	}
	if existing != nil {
		return model.ShoeSeriesResponse{}, NewInvalidRequestError("该品牌下同名系列已存在")
	}

	series := model.ShoeSeries{BrandID: req.BrandID, Gender: req.Gender, Name: name}
	if err := s.repo.CreateSeries(&series); err != nil {
		return model.ShoeSeriesResponse{}, err
	}
	return model.NewShoeSeriesResponse(series), nil
}

// CreateLibraryItems 管理员新增球鞋库条目；一次只插入一条，配色为合并后的单值（如 Red/Black）。
func (s *ShoeService) CreateLibraryItems(req model.CreateShoeLibraryRequest) (model.CreateShoeLibraryResponse, error) {
	if req.BrandID <= 0 || req.SeriesID <= 0 {
		return model.CreateShoeLibraryResponse{}, NewInvalidRequestError("品牌和系列不能为空")
	}
	brand, err := s.repo.BrandFindByID(req.BrandID)
	if err != nil {
		return model.CreateShoeLibraryResponse{}, err
	}
	if brand == nil {
		return model.CreateShoeLibraryResponse{}, NewInvalidRequestError("品牌不存在")
	}
	series, err := s.repo.SeriesFindByID(req.SeriesID)
	if err != nil {
		return model.CreateShoeLibraryResponse{}, err
	}
	if series == nil {
		return model.CreateShoeLibraryResponse{}, NewInvalidRequestError("系列不存在")
	}
	if series.BrandID != req.BrandID {
		return model.CreateShoeLibraryResponse{}, NewInvalidRequestError("系列不属于该品牌")
	}

	modelName := strings.TrimSpace(req.Model)
	if modelName == "" {
		return model.CreateShoeLibraryResponse{}, NewInvalidRequestError("型号不能为空")
	}
	if req.Gender < 0 || req.Gender > 3 {
		return model.CreateShoeLibraryResponse{}, NewInvalidRequestError("性别取值不合法")
	}

	colorway := strings.TrimSpace(req.Colorway)
	exists, err := s.repo.LibraryFindDuplicate(req.BrandID, req.SeriesID, modelName, req.Gender, colorway)
	if err != nil {
		return model.CreateShoeLibraryResponse{}, err
	}
	if exists {
		label := colorway
		if label == "" {
			label = "无配色"
		}
		return model.CreateShoeLibraryResponse{}, NewInvalidRequestError(fmt.Sprintf("配色 %s 已存在", label))
	}

	items := []model.ShoeLibrary{{
		BrandID:       req.BrandID,
		Brand:         brand.Name,
		SeriesID:      req.SeriesID,
		Series:        series.Name,
		Model:         modelName,
		Gender:        req.Gender,
		Colorway:      colorway,
		ReleaseYear:   req.ReleaseYear,
		Weight:        req.Weight,
		Width:         req.Width,
		Surface:       req.Surface,
		Price:         req.Price,
		ColorwayCount: req.ColorwayCount,
		FileID:        req.FileID,
		ImageURL:      req.ImageURL,
	}}
	if err := s.repo.CreateLibraryItems(items); err != nil {
		return model.CreateShoeLibraryResponse{}, err
	}

	responses := make([]model.ShoeLibraryItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, model.NewShoeLibraryItemResponse(item))
	}
	return model.CreateShoeLibraryResponse{Created: len(responses), Items: responses}, nil
}

// deriveBrandSlug 由品牌名生成内部 slug：转小写、保留字母数字与连字符；
// 纯中文等无法转换的名字退化为 brand-<md5 前 4 字节>。
func deriveBrandSlug(name string) string {
	source := strings.ToLower(strings.TrimSpace(name))
	var builder strings.Builder
	lastDash := false
	for _, r := range source {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			builder.WriteRune(r)
			lastDash = false
		case r == ' ' || r == '-' || r == '_':
			if !lastDash && builder.Len() > 0 {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}
	slug := strings.Trim(builder.String(), "-")
	if slug != "" {
		return slug
	}
	sum := md5.Sum([]byte(name))
	return fmt.Sprintf("brand-%x", sum[:4])
}
