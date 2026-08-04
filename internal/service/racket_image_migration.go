package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/model"
	"tennisdaily-backend/internal/repository"
	"tennisdaily-backend/internal/storage"
)

const (
	maxImageBytes          = 10 << 20
	imageDownloadTimeout   = 30 * time.Second
	imageDownloadUserAgent = "Mozilla/5.0 (compatible; TennisDaily/1.0)"
)

// ImageMigrationFailure 单条迁移失败信息。
type ImageMigrationFailure struct {
	ID    int64  `json:"id"`
	URL   string `json:"url"`
	Error string `json:"error"`
}

// ImageMigrationResult 一次迁移的结果汇总。
type ImageMigrationResult struct {
	Total    int                     `json:"total"`
	Uploaded int                     `json:"uploaded"`
	Failed   int                     `json:"failed"`
	Failures []ImageMigrationFailure `json:"failures"`
}

// RacketImageMigration 将 racket_library.image_url 的旧外链图片迁移到微信云托管对象存储：
// 下载 → 上传（带文件元数据）→ 更新 file_id（不动 image_url）→ 写入上传记录表。
type RacketImageMigration struct {
	repo       *repository.RacketRepository
	cloudStore *storage.CloudStorage
	httpClient *http.Client
	running    atomic.Bool
}

func NewRacketImageMigration(repo *repository.RacketRepository, cloudStore *storage.CloudStorage) *RacketImageMigration {
	return &RacketImageMigration{
		repo:       repo,
		cloudStore: cloudStore,
		httpClient: &http.Client{Timeout: imageDownloadTimeout},
	}
}

// Pending 返回当前待迁移的球拍库图片（用于日志/预览）。
func (m *RacketImageMigration) Pending() ([]model.RacketLibrary, error) {
	return m.repo.LibraryImagesPending()
}

// Run 执行一轮迁移。同一时间只允许一轮运行，避免定时任务与手动触发重叠。
func (m *RacketImageMigration) Run(ctx context.Context) (*ImageMigrationResult, error) {
	if !m.running.CompareAndSwap(false, true) {
		return nil, errors.New("racket image migration already running")
	}
	defer m.running.Store(false)

	items, err := m.repo.LibraryImagesPending()
	if err != nil {
		return nil, fmt.Errorf("query pending racket images: %w", err)
	}

	result := &ImageMigrationResult{Total: len(items)}
	logger.Debug("racket image migration start total=%d", len(items))

	for _, item := range items {
		fileID, err := m.migrateOne(ctx, item)
		if err != nil {
			result.Failed++
			result.Failures = append(result.Failures, ImageMigrationFailure{
				ID:    item.ID,
				URL:   item.ImageURL,
				Error: err.Error(),
			})
			logger.Debug("racket image migration failed id=%d url=%s err=%v", item.ID, item.ImageURL, err)
			continue
		}

		result.Uploaded++
		logger.Debug("racket image migration ok id=%d fileID=%s", item.ID, fileID)
	}

	logger.Debug("racket image migration finish total=%d uploaded=%d failed=%d", result.Total, result.Uploaded, result.Failed)
	return result, nil
}

func (m *RacketImageMigration) migrateOne(ctx context.Context, item model.RacketLibrary) (string, error) {
	if item.ImageURL == "" {
		return "", errors.New("empty image url")
	}

	data, contentType, err := m.download(ctx, item.ImageURL)
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

// download 将图片下载到内存；不落盘，因此无需再删除临时文件。
func (m *RacketImageMigration) download(ctx context.Context, rawURL string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", imageDownloadUserAgent)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxImageBytes+1))
	if err != nil {
		return nil, "", err
	}
	if len(data) > maxImageBytes {
		return nil, "", fmt.Errorf("image exceeds %d bytes", maxImageBytes)
	}
	return data, contentType, nil
}
