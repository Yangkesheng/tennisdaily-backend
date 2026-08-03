package service

import (
	"errors"
	"fmt"
	"time"

	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"

	"gorm.io/gorm"
)

type RacketService struct {
	repo            *repository.RacketRepository
	userRepo        *repository.UserRepository
	contentSecurity *ContentSecurityService
	loc             *time.Location
}

func NewRacketService(repo *repository.RacketRepository, userRepo *repository.UserRepository, contentSecurity *ContentSecurityService) *RacketService {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	return &RacketService{repo: repo, userRepo: userRepo, contentSecurity: contentSecurity, loc: loc}
}

func (s *RacketService) Brands() ([]model.RacketBrandResponse, error) {
	items, err := s.repo.BrandList()
	if err != nil {
		return nil, err
	}
	responses := make([]model.RacketBrandResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, model.NewRacketBrandResponse(item))
	}
	return responses, nil
}

func (s *RacketService) Series(query model.RacketSeriesQuery) ([]model.RacketSeriesResponse, error) {
	items, err := s.repo.SeriesList(query)
	if err != nil {
		return nil, err
	}
	responses := make([]model.RacketSeriesResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, model.NewRacketSeriesResponse(item))
	}
	return responses, nil
}

func (s *RacketService) Library(query model.RacketLibraryQuery) ([]model.RacketLibraryBrandGroupResponse, error) {
	items, err := s.repo.LibraryList(query)
	if err != nil {
		return nil, err
	}

	groups := make([]model.RacketLibraryBrandGroupResponse, 0)
	brandIndex := make(map[string]int)
	for _, item := range items {
		idx, ok := brandIndex[item.Brand]
		if !ok {
			idx = len(groups)
			brandIndex[item.Brand] = idx
			groups = append(groups, model.RacketLibraryBrandGroupResponse{BrandID: item.BrandID, Brand: item.Brand, Items: []model.RacketLibraryItemResponse{}})
		}
		groups[idx].Items = append(groups[idx].Items, model.NewRacketLibraryItemResponse(item))
	}
	return groups, nil
}

func (s *RacketService) LibraryStats() ([]model.RacketLibraryBrandStatsResponse, error) {
	brands, err := s.repo.LibraryBrandStats()
	if err != nil {
		return nil, err
	}
	series, err := s.repo.LibrarySeriesStats()
	if err != nil {
		return nil, err
	}

	responses := make([]model.RacketLibraryBrandStatsResponse, 0, len(brands))
	type brandStatsKey struct {
		BrandID int64
		Brand   string
	}
	brandIndex := make(map[brandStatsKey]int, len(brands))
	for _, brand := range brands {
		brandIndex[brandStatsKey{BrandID: brand.BrandID, Brand: brand.Brand}] = len(responses)
		responses = append(responses, model.RacketLibraryBrandStatsResponse{
			BrandID: brand.BrandID,
			Brand:   brand.Brand,
			Count:   brand.Count,
			Series:  []model.RacketLibrarySeriesStatsResponse{},
		})
	}
	for _, item := range series {
		idx, ok := brandIndex[brandStatsKey{BrandID: item.BrandID, Brand: item.Brand}]
		if !ok {
			continue
		}
		responses[idx].Series = append(responses[idx].Series, model.RacketLibrarySeriesStatsResponse{
			SeriesID: item.SeriesID,
			Series:   item.Series,
			Count:    item.Count,
		})
	}
	return responses, nil
}

func (s *RacketService) MyRacketsForSession(userID int64) ([]model.RacketResponse, error) {
	return s.Selectable(userID)
}

func (s *RacketService) Stats(userID int64) (model.RacketStatsResponse, error) {
	count, racketCost, stringingCost, err := s.repo.Stats(userID)
	if err != nil {
		return model.RacketStatsResponse{}, err
	}
	totalCost := racketCost + stringingCost
	return model.RacketStatsResponse{
		RacketCount:       count,
		RacketCost:        racketCost,
		StringingCost:     stringingCost,
		TotalCost:         totalCost,
		RacketCostText:    fmt.Sprintf("%.2f", racketCost),
		StringingCostText: fmt.Sprintf("%.2f", stringingCost),
		TotalCostText:     fmt.Sprintf("%.2f", totalCost),
	}, nil
}

func (s *RacketService) List(userID int64, includeRetired bool) ([]model.RacketResponse, error) {
	rackets, err := s.repo.List(userID, includeRetired)
	if err != nil {
		return nil, err
	}
	return s.enrichRackets(userID, rackets)
}

func (s *RacketService) Selectable(userID int64) ([]model.RacketResponse, error) {
	rackets, err := s.repo.Selectable(userID)
	if err != nil {
		return nil, err
	}
	return s.enrichRackets(userID, rackets)
}

func (s *RacketService) Primary(userID int64) (*model.RacketResponse, error) {
	racket, err := s.repo.Primary(userID)
	if err != nil {
		return nil, err
	}
	if racket == nil {
		return nil, nil
	}
	responses, err := s.enrichRackets(userID, []model.Racket{*racket})
	if err != nil {
		return nil, err
	}
	return &responses[0], nil
}

func (s *RacketService) Create(userID int64, req model.CreateRacketRequest) (model.RacketResponse, error) {
	if err := s.applyLibraryDefaults(&req.LibraryID, &req.Name, &req.Brand, &req.Model); err != nil {
		return model.RacketResponse{}, err
	}
	if req.Name == "" {
		return model.RacketResponse{}, ErrInvalidRequest
	}
	if err := s.checkUserInputTexts(userID, req.Name, req.Brand, req.Model); err != nil {
		return model.RacketResponse{}, err
	}

	purchaseDate, err := s.parseOptionalDate(req.PurchaseDate)
	if err != nil {
		return model.RacketResponse{}, ErrInvalidRequest
	}

	status := model.RacketStatusActive
	if req.Status != 0 {
		if req.Status != model.RacketStatusPrimary && req.Status != model.RacketStatusActive {
			return model.RacketResponse{}, ErrInvalidRequest
		}
		status = req.Status
	}

	racket := model.Racket{
		UserID:        userID,
		LibraryID:     req.LibraryID,
		Name:          req.Name,
		Brand:         req.Brand,
		Model:         req.Model,
		Status:        status,
		PurchaseDate:  purchaseDate,
		PurchasePrice: req.PurchasePrice,
	}
	if err := s.repo.Create(&racket); err != nil {
		return model.RacketResponse{}, err
	}
	if status == model.RacketStatusPrimary {
		if err := s.repo.SetPrimary(userID, racket.ID); err != nil {
			return model.RacketResponse{}, err
		}
		racket.Status = model.RacketStatusPrimary
	}

	responses, err := s.enrichRackets(userID, []model.Racket{racket})
	if err != nil {
		return model.RacketResponse{}, err
	}
	return responses[0], nil
}

func (s *RacketService) Detail(userID, id int64) (model.RacketDetailResponse, error) {
	racket, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.RacketDetailResponse{}, err
	}
	if racket == nil {
		return model.RacketDetailResponse{}, ErrNotFound
	}

	responses, err := s.enrichRackets(userID, []model.Racket{*racket})
	if err != nil {
		return model.RacketDetailResponse{}, err
	}

	records, err := s.repo.ListStringingRecords(userID, id)
	if err != nil {
		return model.RacketDetailResponse{}, err
	}
	recordResponses := make([]model.StringingRecordResponse, 0, len(records))
	for _, record := range records {
		recordResponses = append(recordResponses, model.NewStringingRecordResponse(record))
	}

	return model.RacketDetailResponse{Racket: responses[0], StringingRecords: recordResponses}, nil
}

func (s *RacketService) Update(userID, id int64, req model.UpdateRacketRequest) (model.RacketResponse, error) {
	if err := s.applyLibraryDefaults(&req.LibraryID, &req.Name, &req.Brand, &req.Model); err != nil {
		return model.RacketResponse{}, err
	}
	if req.Name == "" {
		return model.RacketResponse{}, ErrInvalidRequest
	}
	if err := s.checkUserInputTexts(userID, req.Name, req.Brand, req.Model); err != nil {
		return model.RacketResponse{}, err
	}

	racket, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.RacketResponse{}, err
	}
	if racket == nil {
		return model.RacketResponse{}, ErrNotFound
	}

	purchaseDate, err := s.parseOptionalDate(req.PurchaseDate)
	if err != nil {
		return model.RacketResponse{}, ErrInvalidRequest
	}

	racket.LibraryID = req.LibraryID
	racket.Name = req.Name
	racket.Brand = req.Brand
	racket.Model = req.Model
	racket.PurchaseDate = purchaseDate
	racket.PurchasePrice = req.PurchasePrice
	if req.Status >= model.RacketStatusPrimary && req.Status <= model.RacketStatusRetired {
		racket.Status = req.Status
	}

	if racket.Status == model.RacketStatusPrimary {
		if err := s.repo.Update(racket); err != nil {
			return model.RacketResponse{}, err
		}
		if err := s.repo.SetPrimary(userID, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return model.RacketResponse{}, ErrNotFound
			}
			return model.RacketResponse{}, err
		}
		racket.Status = model.RacketStatusPrimary
	} else if err := s.repo.Update(racket); err != nil {
		return model.RacketResponse{}, err
	}

	updated, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.RacketResponse{}, err
	}
	responses, err := s.enrichRackets(userID, []model.Racket{*updated})
	if err != nil {
		return model.RacketResponse{}, err
	}
	return responses[0], nil
}

func (s *RacketService) SetPrimary(userID, id int64) (model.RacketResponse, error) {
	if err := s.repo.SetPrimary(userID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.RacketResponse{}, ErrNotFound
		}
		return model.RacketResponse{}, err
	}
	racket, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.RacketResponse{}, err
	}
	responses, err := s.enrichRackets(userID, []model.Racket{*racket})
	if err != nil {
		return model.RacketResponse{}, err
	}
	return responses[0], nil
}

func (s *RacketService) Retire(userID, id int64) (model.RacketResponse, error) {
	ok, err := s.repo.Retire(userID, id)
	if err != nil {
		return model.RacketResponse{}, err
	}
	if !ok {
		return model.RacketResponse{}, ErrNotFound
	}
	racket, err := s.repo.FindByID(userID, id)
	if err != nil {
		return model.RacketResponse{}, err
	}
	responses, err := s.enrichRackets(userID, []model.Racket{*racket})
	if err != nil {
		return model.RacketResponse{}, err
	}
	return responses[0], nil
}

func (s *RacketService) Delete(userID, id int64) error {
	deleted, err := s.repo.SoftDelete(userID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func (s *RacketService) AddStringingRecord(userID, racketID int64, req model.CreateStringingRecordRequest) error {
	if req.StringName == "" || req.StringDate == "" {
		return ErrInvalidRequest
	}
	if err := s.checkUserInputTexts(userID, req.StringName, req.StoreName); err != nil {
		return err
	}
	if _, err := s.requireRacket(userID, racketID); err != nil {
		return err
	}
	stringDate, err := s.parseStringingDate(req.StringDate)
	if err != nil {
		return ErrInvalidRequest
	}
	record := model.RacketStringingRecord{
		UserID:            userID,
		RacketID:          racketID,
		StringName:        req.StringName,
		StoreName:         req.StoreName,
		VerticalTension:   req.VerticalTension,
		HorizontalTension: req.HorizontalTension,
		Cost:              req.Cost,
		StringDate:        stringDate,
	}
	return s.repo.CreateStringingRecord(&record)
}

func (s *RacketService) CreateStringingRecord(userID, racketID int64, req model.CreateStringingRecordRequest) (model.StringingRecordResponse, error) {
	if req.StringName == "" || req.StringDate == "" {
		return model.StringingRecordResponse{}, ErrInvalidRequest
	}
	if err := s.checkUserInputTexts(userID, req.StringName, req.StoreName); err != nil {
		return model.StringingRecordResponse{}, err
	}
	if _, err := s.requireRacket(userID, racketID); err != nil {
		return model.StringingRecordResponse{}, err
	}
	stringDate, err := s.parseStringingDate(req.StringDate)
	if err != nil {
		return model.StringingRecordResponse{}, ErrInvalidRequest
	}
	record := model.RacketStringingRecord{
		UserID:            userID,
		RacketID:          racketID,
		StringName:        req.StringName,
		StoreName:         req.StoreName,
		VerticalTension:   req.VerticalTension,
		HorizontalTension: req.HorizontalTension,
		Cost:              req.Cost,
		StringDate:        stringDate,
	}
	if err := s.repo.CreateStringingRecord(&record); err != nil {
		return model.StringingRecordResponse{}, err
	}
	return model.NewStringingRecordResponse(record), nil
}

func (s *RacketService) UpdateStringingRecord(userID, racketID, recordID int64, req model.UpdateStringingRecordRequest) (model.StringingRecordResponse, error) {
	if req.StringName == "" || req.StringDate == "" {
		return model.StringingRecordResponse{}, ErrInvalidRequest
	}
	if err := s.checkUserInputTexts(userID, req.StringName, req.StoreName); err != nil {
		return model.StringingRecordResponse{}, err
	}
	if _, err := s.requireRacket(userID, racketID); err != nil {
		return model.StringingRecordResponse{}, err
	}
	record, err := s.repo.FindStringingRecordByID(userID, racketID, recordID)
	if err != nil {
		return model.StringingRecordResponse{}, err
	}
	if record == nil {
		return model.StringingRecordResponse{}, ErrNotFound
	}
	stringDate, err := s.parseStringingDate(req.StringDate)
	if err != nil {
		return model.StringingRecordResponse{}, ErrInvalidRequest
	}
	record.StringName = req.StringName
	record.StoreName = req.StoreName
	record.VerticalTension = req.VerticalTension
	record.HorizontalTension = req.HorizontalTension
	record.Cost = req.Cost
	record.StringDate = stringDate
	if err := s.repo.UpdateStringingRecord(record); err != nil {
		return model.StringingRecordResponse{}, err
	}
	return model.NewStringingRecordResponse(*record), nil
}

func (s *RacketService) DeleteStringingRecord(userID, racketID, recordID int64) error {
	if _, err := s.requireRacket(userID, racketID); err != nil {
		return err
	}
	deleted, err := s.repo.SoftDeleteStringingRecord(userID, racketID, recordID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func (s *RacketService) checkUserInputTexts(userID int64, texts ...string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUnauthorized
	}
	return s.contentSecurity.CheckTexts(user.OpenID, texts...)
}

func (s *RacketService) requireRacket(userID, racketID int64) (*model.Racket, error) {
	racket, err := s.repo.FindByID(userID, racketID)
	if err != nil {
		return nil, err
	}
	if racket == nil {
		return nil, ErrNotFound
	}
	return racket, nil
}

func (s *RacketService) enrichRackets(userID int64, rackets []model.Racket) ([]model.RacketResponse, error) {
	ids := make([]int64, 0, len(rackets))
	libraryIDs := make([]int64, 0, len(rackets))
	libraryIDSet := make(map[int64]struct{})
	for _, racket := range rackets {
		ids = append(ids, racket.ID)
		if racket.LibraryID > 0 {
			if _, ok := libraryIDSet[racket.LibraryID]; !ok {
				libraryIDSet[racket.LibraryID] = struct{}{}
				libraryIDs = append(libraryIDs, racket.LibraryID)
			}
		}
	}
	libraries, err := s.repo.LibraryFindByIDs(libraryIDs)
	if err != nil {
		return nil, err
	}
	latest, err := s.repo.LatestStringingRecords(userID, ids)
	if err != nil {
		return nil, err
	}
	usage, err := s.repo.UsageStats(userID, ids)
	if err != nil {
		return nil, err
	}
	afterStringingUsage, err := s.repo.UsageStatsSince(userID, latest)
	if err != nil {
		return nil, err
	}

	responses := make([]model.RacketResponse, 0, len(rackets))
	for _, racket := range rackets {
		if library, ok := libraries[racket.LibraryID]; ok {
			racket.ReleaseYear = library.ReleaseYear
			racket.Weight = library.Weight
			racket.HeadSize = library.HeadSize
			racket.StringPattern = library.StringPattern
			racket.FileID = library.FileID
		}
		if record, ok := latest[racket.ID]; ok {
			racket.StringName = record.StringName
			racket.StoreName = record.StoreName
			racket.VerticalTension = record.VerticalTension
			racket.HorizontalTension = record.HorizontalTension
			racket.LastStringDate = record.StringDate.Format("2006-01-02")
			racket.LastStringCost = &record.Cost
			latestRecord := model.NewStringingRecordResponse(record)
			racket.LatestStringingRecord = &latestRecord
		}
		stats := usage[racket.ID]
		racket.UsageCount = stats.Count
		racket.UsageMinutes = stats.Minutes
		racket.UsageHours = stats.Hours
		racket.TotalMinutes = stats.Minutes
		racket.TotalHours = stats.Hours
		if stats, ok := afterStringingUsage[racket.ID]; ok {
			racket.AfterStringingUsageCount = stats.Count
			racket.AfterStringingUsageMinutes = stats.Minutes
			racket.AfterStringingUsageHours = stats.Hours
		}
		responses = append(responses, model.NewRacketResponse(racket))
	}
	return responses, nil
}

func (s *RacketService) applyLibraryDefaults(libraryID *int64, name *string, brand *string, racketModel *string) error {
	if libraryID == nil || *libraryID == 0 {
		return nil
	}
	item, err := s.repo.LibraryFindByID(*libraryID)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrInvalidRequest
	}
	if *brand == "" {
		*brand = item.Brand
	}
	if *racketModel == "" {
		*racketModel = item.Model
	}
	if *name == "" {
		*name = item.Brand + " " + item.Model
	}
	return nil
}

func (s *RacketService) parseStringingDate(value string) (time.Time, error) {
	layouts := []string{"2006-01-02 15:04", "2006-01-02"}
	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, s.loc)
		if err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, ErrInvalidRequest
}

func (s *RacketService) parseOptionalDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.ParseInLocation("2006-01-02", value, s.loc)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
