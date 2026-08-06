package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"

	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
	"tennisdaily-backend/internal/storage"
)

// ShoeImageMigration 将 shoe_library.image_url 的旧外链图片迁移到微信云托管对象存储：
// 下载 → 上传（带文件元数据）→ 更新 file_id（不动 image_url）→ 写入上传记录表。
type ShoeImageMigration struct {
	repo       *repository.ShoeRepository
	cloudStore *storage.CloudStorage
	httpClient *http.Client
	running    atomic.Bool
}

func NewShoeImageMigration(repo *repository.ShoeRepository, cloudStore *storage.CloudStorage) *ShoeImageMigration {
	return &ShoeImageMigration{
		repo:       repo,
		cloudStore: cloudStore,
		httpClient: &http.Client{Timeout: imageDownloadTimeout},
	}
}

// Pending 返回当前待迁移的球鞋库图片（用于日志/预览）。
func (m *ShoeImageMigration) Pending() ([]model.ShoeLibrary, error) {
	return m.repo.LibraryImagesPending()
}

// Run 执行一轮迁移。同一时间只允许一轮运行，避免定时任务与手动触发重叠。
func (m *ShoeImageMigration) Run(ctx context.Context) (*ImageMigrationResult, error) {
	if !m.running.CompareAndSwap(false, true) {
		return nil, errors.New("shoe image migration already running")
	}
	defer m.running.Store(false)

	items, err := m.repo.LibraryImagesPending()
	if err != nil {
		return nil, fmt.Errorf("query pending shoe images: %w", err)
	}

	result := &ImageMigrationResult{Total: len(items)}
	logger.Debug("shoe image migration start total=%d", len(items))

	for _, item := range items {
		fileID, err := m.migrateOne(ctx, item)
		if err != nil {
			result.Failed++
			result.Failures = append(result.Failures, ImageMigrationFailure{
				ID:    item.ID,
				URL:   item.ImageURL,
				Error: err.Error(),
			})
			logger.Warn("shoe image migration failed id=%d url=%s err=%v", item.ID, item.ImageURL, err)
			continue
		}

		result.Uploaded++
		logger.Debug("shoe image migration ok id=%d fileID=%s", item.ID, fileID)
	}

	logger.Debug("shoe image migration finish total=%d uploaded=%d failed=%d", result.Total, result.Uploaded, result.Failed)
	return result, nil
}

func (m *ShoeImageMigration) migrateOne(ctx context.Context, item model.ShoeLibrary) (string, error) {
	if item.ImageURL == "" {
		return "", errors.New("empty image url")
	}

	data, contentType, err := downloadImage(ctx, m.httpClient, item.ImageURL)
	if err != nil {
		return "", fmt.Errorf("download: %w", err)
	}

	ext := storage.ExtensionFromURL(item.ImageURL, contentType)
	objectKey := m.cloudStore.ObjectKey(item.ID, item.Brand, item.Series, ext)

	fileID, err := m.cloudStore.Upload(ctx, objectKey, data, contentType)
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}

	// 上传成功，写入 file_id 与上传记录（不动 image_url，保留旧链接兜底）。
	if err := m.repo.ApplyLibraryImageUpload(item.ID, fileID, objectKey, item.ImageURL); err != nil {
		return "", fmt.Errorf("update db: %w", err)
	}
	return fileID, nil
}
