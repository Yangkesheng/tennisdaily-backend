package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
	"unicode"

	"github.com/tencentyun/cos-go-sdk-v5"
)

const (
	authEndpoint   = "http://api.weixin.qq.com/_/cos/getauth"
	metaIDEndpoint = "http://api.weixin.qq.com/_/cos/metaid/encode"
)

// Config 微信云托管对象存储配置。
type Config struct {
	EnvID  string
	Bucket string
	Region string
	Folder string
}

// CloudStorage 封装微信云托管对象存储的服务端上传能力：
// 通过开放接口服务获取临时密钥，写入文件元数据后用 COS SDK 上传。
type CloudStorage struct {
	cfg        Config
	httpClient *http.Client
}

func New(cfg Config) *CloudStorage {
	return &CloudStorage{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FileID 拼接对象存储文件 cloudID：
// cloud://{环境ID}.{存储桶ID}/{路径}
func (s *CloudStorage) FileID(objectKey string) string {
	return fmt.Sprintf("cloud://%s.%s/%s", s.cfg.EnvID, s.cfg.Bucket, objectKey)
}

// ObjectKey 生成球拍库图片的存储路径：
// {folder}/{brand}/{series}/{id}{ext}
// brand/series 会清洗为 cloudPath 允许的字符（中文保留，其余特殊字符转 -）。
func (s *CloudStorage) ObjectKey(id int64, brand, series, ext string) string {
	parts := []string{s.cfg.Folder}
	if brand != "" {
		parts = append(parts, SanitizeSegment(brand))
	}
	if series != "" {
		parts = append(parts, SanitizeSegment(series))
	}
	return fmt.Sprintf("%s/%d%s", strings.Join(parts, "/"), id, ext)
}

// SanitizeSegment 将目录片段清洗为 cloudPath 允许的字符。
// 允许字母数字、中文、! - _ . *，其余字符转为 -，连续 - 合并。
func SanitizeSegment(value string) string {
	value = strings.TrimSpace(value)
	mapped := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			return r
		case r == '!' || r == '-' || r == '_' || r == '.' || r == '*':
			return r
		case unicode.Is(unicode.Han, r):
			return r
		default:
			return '-'
		}
	}, value)

	for strings.Contains(mapped, "--") {
		mapped = strings.ReplaceAll(mapped, "--", "-")
	}
	mapped = strings.Trim(mapped, "-.")
	if mapped == "" {
		return "unknown"
	}
	return mapped
}

type tempAuth struct {
	TmpSecretID  string `json:"TmpSecretId"`
	TmpSecretKey string `json:"TmpSecretKey"`
	Token        string `json:"Token"`
	ExpiredTime  int64  `json:"ExpiredTime"`
}

func (s *CloudStorage) getAuth(ctx context.Context) (*tempAuth, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, authEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cos getauth: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("cos getauth read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cos getauth status=%d body=%s", resp.StatusCode, truncate(body))
	}

	var auth tempAuth
	if err := json.Unmarshal(body, &auth); err != nil {
		return nil, fmt.Errorf("cos getauth parse: %w", err)
	}
	if auth.TmpSecretID == "" || auth.TmpSecretKey == "" {
		return nil, fmt.Errorf("cos getauth empty credentials: %s", truncate(body))
	}
	return &auth, nil
}

type metaIDRequest struct {
	OpenID string   `json:"openid"`
	Bucket string   `json:"bucket"`
	Paths  []string `json:"paths"`
}

type metaIDResponse struct {
	Errcode  int    `json:"errcode"`
	Errmsg   string `json:"errmsg"`
	RespData struct {
		XcosMetaFieldStrs []string `json:"x_cos_meta_field_strs"`
	} `json:"respdata"`
}

// encodeMetaID 生成文件元数据（管理端 openid 为空），上传时必须写入 x-cos-meta-fileid，
// 否则小程序端无法访问该文件。
func (s *CloudStorage) encodeMetaID(ctx context.Context, objectKey string) (string, error) {
	payload, err := json.Marshal(metaIDRequest{
		OpenID: "",
		Bucket: s.cfg.Bucket,
		Paths:  []string{objectKey},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, metaIDEndpoint, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cos metaid encode: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("cos metaid encode read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cos metaid encode status=%d body=%s", resp.StatusCode, truncate(body))
	}

	var result metaIDResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("cos metaid encode parse: %w", err)
	}
	if result.Errcode != 0 || len(result.RespData.XcosMetaFieldStrs) == 0 {
		return "", fmt.Errorf("cos metaid encode failed errcode=%d errmsg=%s", result.Errcode, result.Errmsg)
	}
	return result.RespData.XcosMetaFieldStrs[0], nil
}

// Upload 将图片上传到对象存储，返回 cloudID。
func (s *CloudStorage) Upload(ctx context.Context, objectKey string, data []byte, contentType string) (string, error) {
	auth, err := s.getAuth(ctx)
	if err != nil {
		return "", err
	}

	metaID, err := s.encodeMetaID(ctx, objectKey)
	if err != nil {
		return "", err
	}

	bucketURL, err := cos.NewBucketURL(s.cfg.Bucket, s.cfg.Region, true)
	if err != nil {
		return "", fmt.Errorf("cos bucket url: %w", err)
	}

	client := cos.NewClient(
		&cos.BaseURL{BucketURL: bucketURL},
		&http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:     auth.TmpSecretID,
				SecretKey:    auth.TmpSecretKey,
				SessionToken: auth.Token,
			},
			Timeout: 60 * time.Second,
		},
	)

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
			XCosMetaXXX: &http.Header{
				"x-cos-meta-fileid": []string{metaID},
			},
		},
	}

	if _, err := client.Object.Put(ctx, objectKey, bytes.NewReader(data), opt); err != nil {
		return "", fmt.Errorf("cos put object %s: %w", objectKey, err)
	}
	return s.FileID(objectKey), nil
}

// ExtensionFromURL 从 URL 路径或 Content-Type 推断扩展名，兜底 .jpg。
func ExtensionFromURL(rawURL, contentType string) string {
	parsed, err := url.Parse(rawURL)
	if err == nil {
		ext := strings.ToLower(path.Ext(parsed.Path))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp", ".gif", ".avif":
			if ext == ".jpeg" {
				return ".jpg"
			}
			return ext
		}
	}

	if ct, _, err := mime.ParseMediaType(contentType); err == nil {
		switch ct {
		case "image/jpeg":
			return ".jpg"
		case "image/png":
			return ".png"
		case "image/webp":
			return ".webp"
		case "image/gif":
			return ".gif"
		case "image/avif":
			return ".avif"
		}
	}
	return ".jpg"
}

func truncate(b []byte) string {
	const max = 200
	if len(b) <= max {
		return string(b)
	}
	return string(b[:max])
}
