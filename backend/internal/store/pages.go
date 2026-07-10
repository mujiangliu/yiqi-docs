// backend/internal/store/pages.go
package store

import (
	"errors"

	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/model"
)

// CreatePage 创建页面。先校验站点对当前用户可见
func (s *Store) CreatePage(p model.Page) (*model.Page, error) {
	_, err := s.GetSiteByID(p.SiteID)
	if err != nil {
		return nil, err
	}
	err = s.db.Create(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPagesBySite 列出站点下所有页面（含正文，admin 仅自己站点可见）
func (s *Store) ListPagesBySite(siteID uint) ([]model.Page, error) {
	_, err := s.GetSiteByID(siteID)
	if err != nil {
		return nil, err
	}
	var pages []model.Page
	err = s.db.Where("site_id = ?", siteID).Order("sort asc, id asc").Find(&pages).Error
	return pages, err
}

// ListPageTreeBySite 列出站点页面树所需字段，不加载正文。
func (s *Store) ListPageTreeBySite(siteID uint) ([]model.Page, error) {
	_, err := s.GetSiteByID(siteID)
	if err != nil {
		return nil, err
	}
	var pages []model.Page
	err = s.db.
		Select("id", "site_id", "parent_id", "slug", "title", "sort", "created_at", "updated_at", "deleted_at").
		Where("site_id = ?", siteID).
		Order("sort asc, id asc").
		Find(&pages).Error
	return pages, err
}

// GetPageByID 按 ID 查询（admin 仅自己站点下的页面可见）
func (s *Store) GetPageByID(id uint) (*model.Page, error) {
	var page model.Page
	err := s.db.First(&page, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	// 校验页面所属站点对当前用户可见
	_, err = s.GetSiteByID(page.SiteID)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// UpdatePage 更新页面字段
func (s *Store) UpdatePage(id uint, fields map[string]interface{}) (*model.Page, error) {
	_, err := s.GetPageByID(id)
	if err != nil {
		return nil, err
	}
	delete(fields, "site_id") // 不允许改所属站点
	err = s.db.Model(&model.Page{}).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return nil, err
	}
	return s.GetPageByID(id)
}

// DeletePage 删除页面
func (s *Store) DeletePage(id uint) error {
	_, err := s.GetPageByID(id)
	if err != nil {
		return err
	}
	return s.db.Delete(&model.Page{}, id).Error
}

// ReorderPages 批量更新同级页面 sort
func (s *Store) ReorderPages(ids []uint) error {
	for i, id := range ids {
		err := s.db.Model(&model.Page{}).Where("id = ?", id).Update("sort", i).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// ListPublishedPagesBySite 公开视角：列出 published 站点的所有页面
func (s *Store) ListPublishedPagesBySite(siteID uint) ([]model.Page, error) {
	var pages []model.Page
	err := s.db.Where("site_id = ?", siteID).Order("sort asc, id asc").Find(&pages).Error
	return pages, err
}

// ListPublishedPageTreeBySite 公开视角：只列出目录所需字段，不加载正文。
func (s *Store) ListPublishedPageTreeBySite(siteID uint) ([]model.Page, error) {
	var pages []model.Page
	err := s.db.
		Select("id", "site_id", "parent_id", "slug", "title", "sort", "created_at", "updated_at", "deleted_at").
		Where("site_id = ?", siteID).
		Order("sort asc, id asc").
		Find(&pages).Error
	return pages, err
}

// GetPublishedPageByID 公开视角：读取单页正文。
func (s *Store) GetPublishedPageByID(siteID uint, pageID uint) (*model.Page, error) {
	var page model.Page
	err := s.db.Where("site_id = ?", siteID).First(&page, pageID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &page, nil
}
