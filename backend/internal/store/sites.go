// backend/internal/store/sites.go
package store

import (
	"errors"

	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/model"
)

// CreateSite 创建站点。owner_id 由调用方从 user 注入
func (s *Store) CreateSite(site model.Site) (*model.Site, error) {
	if s.user != nil && s.user.Role == model.RoleAdmin {
		site.OwnerID = s.user.ID // 强制归属当前 admin
	}
	err := s.db.Create(&site).Error
	if err != nil {
		return nil, err
	}
	return &site, nil
}

// ListSites 列出当前用户可见的站点（admin 只看自己的，super 看全部）
func (s *Store) ListSites() ([]model.Site, error) {
	var sites []model.Site
	err := s.scopedDB().Order("id asc").Find(&sites).Error
	return sites, err
}

// GetSiteByID 按 ID 查询（admin 仅自己可见）
func (s *Store) GetSiteByID(id uint) (*model.Site, error) {
	var site model.Site
	err := s.scopedDB().First(&site, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &site, nil
}

// GetSiteByPath 按 path 查询（admin 仅自己可见）
func (s *Store) GetSiteByPath(path string) (*model.Site, error) {
	var site model.Site
	err := s.scopedDB().Where("path = ?", path).First(&site).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &site, nil
}

// UpdateSite 更新站点字段（admin 仅自己的）
func (s *Store) UpdateSite(id uint, fields map[string]interface{}) (*model.Site, error) {
	// 先校验可见性
	_, err := s.GetSiteByID(id)
	if err != nil {
		return nil, err
	}
	// 不允许 admin 修改 owner_id
	delete(fields, "owner_id")
	err = s.scopedDB().Model(&model.Site{}).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return nil, err
	}
	return s.GetSiteByID(id)
}

// DeleteSite 删除站点（级联删 pages/media）。admin 仅自己的
func (s *Store) DeleteSite(id uint) error {
	_, err := s.GetSiteByID(id)
	if err != nil {
		return err
	}
	// 级联删除
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if e := tx.Where("site_id = ?", id).Delete(&model.Page{}).Error; e != nil {
			return e
		}
		if e := tx.Where("site_id = ?", id).Delete(&model.Media{}).Error; e != nil {
			return e
		}
		return tx.Delete(&model.Site{}, id).Error
	})
	return err
}

// GetPublishedSiteByPath 公开视角：仅返回 published 站点
func (s *Store) GetPublishedSiteByPath(path string) (*model.Site, error) {
	var site model.Site
	err := s.db.Where("path = ? AND status = ?", path, model.SiteStatusPublished).First(&site).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &site, nil
}
