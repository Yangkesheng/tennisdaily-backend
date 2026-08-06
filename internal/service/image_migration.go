package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
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

// downloadImage 将图片下载到内存；不落盘，因此无需再删除临时文件。
func downloadImage(ctx context.Context, client *http.Client, rawURL string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", imageDownloadUserAgent)

	resp, err := client.Do(req)
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
