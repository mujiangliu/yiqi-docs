// backend/internal/scraper/media.go
package scraper

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/model"
)

var imgURLRegex = regexp.MustCompile(`!\[[^\]]*\]\((https?://[^)]+)\)`)

// LocalizeImages 扫描 Markdown 中的图片 URL，下载到 media 表，重写为 /api/media/:hash
func LocalizeImages(db *gorm.DB, siteID uint, md string) string {
	matches := imgURLRegex.FindAllStringSubmatch(md, -1)
	for _, m := range matches {
		originalURL := m[1]
		newURL, err := downloadAndStore(db, siteID, originalURL)
		if err != nil {
			continue // 失败保留原 URL
		}
		md = strings.ReplaceAll(md, originalURL, newURL)
	}
	return md
}

func downloadAndStore(db *gorm.DB, siteID uint, url string) (string, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])
	mime := resp.Header.Get("Content-Type")
	if mime == "" {
		mime = "application/octet-stream"
	}
	filename := url[strings.LastIndex(url, "/")+1:]

	if err := InsertMedia(db, siteID, hash, filename, mime, data); err != nil {
		return "", err
	}
	return "/api/media/" + hash, nil
}

// 防 unused（model 在本文件未直接引用，预留后续扩展使用）
var _ = model.Media{}
