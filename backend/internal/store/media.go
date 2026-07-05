// backend/internal/store/media.go
package store

import (
	"errors"

	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/model"
)

// CreateMedia 上传媒体。校验站点可见
func (s *Store) CreateMedia(m model.Media) (*model.Media, error) {
	_, err := s.GetSiteByID(m.SiteID)
	if err != nil {
		return nil, err
	}
	err = s.db.Create(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMediaByHash 公开访问：按 hash 查询（无需鉴权）
func (s *Store) GetMediaByHash(hash string) (*model.Media, error) {
	var m model.Media
	err := s.db.Where("hash = ?", hash).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}
