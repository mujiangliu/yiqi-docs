// backend/internal/scraper/seed.go
package scraper

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/auth"
	"jiaocheng-web/backend/internal/model"
)

// SeedOptions 种子参数
type SeedOptions struct {
	SuperAdminUser   string
	SuperAdminPass   string
	ContentAdminUser string
	ContentAdminPass string
}

// SeedResult 种子结果
type SeedResult struct {
	SuperAdminID   uint
	ContentAdminID uint
	ApiSiteID      uint
	ListSiteID     uint
}

// EnsureAccounts 确保 super_admin 和内容管理员账号存在，返回各自 ID
func EnsureAccounts(db *gorm.DB, opts SeedOptions) (superID, contentID uint, err error) {
	superID, err = upsertUser(db, opts.SuperAdminUser, opts.SuperAdminPass, model.RoleSuperAdmin)
	if err != nil {
		return 0, 0, fmt.Errorf("ensure super_admin: %w", err)
	}
	contentID, err = upsertUser(db, opts.ContentAdminUser, opts.ContentAdminPass, model.RoleAdmin)
	if err != nil {
		return 0, 0, fmt.Errorf("ensure content admin: %w", err)
	}
	return superID, contentID, nil
}

func upsertUser(db *gorm.DB, username, password, role string) (uint, error) {
	if username == "" || password == "" {
		return 0, errors.New("username/password required")
	}
	var u model.User
	err := db.Where("username = ?", username).First(&u).Error
	if err == nil {
		// 已存在，不覆盖密码
		return u.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return 0, err
	}
	u = model.User{Username: username, PasswordHash: hash, Role: role}
	if err := db.Create(&u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

// EnsureSite 幂等创建/取回站点（按 path），归属到 ownerID
func EnsureSite(db *gorm.DB, ownerID uint, path, title, description, status string) (uint, error) {
	var site model.Site
	err := db.Where("path = ?", path).First(&site).Error
	if err == nil {
		// 已存在，更新标题/描述/owner
		site.Title = title
		site.Description = description
		site.Status = status
		site.OwnerID = ownerID
		if e := db.Save(&site).Error; e != nil {
			return 0, e
		}
		return site.ID, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	site = model.Site{
		OwnerID:     ownerID,
		Path:        path,
		Title:       title,
		Description: description,
		Status:      status,
	}
	if err := db.Create(&site).Error; err != nil {
		return 0, err
	}
	return site.ID, nil
}

// ClearPages 删除站点下所有页面（重新抓取前清理）
func ClearPages(db *gorm.DB, siteID uint) error {
	return db.Where("site_id = ?", siteID).Delete(&model.Page{}).Error
}

// InsertPage 插入一个页面，返回 id
func InsertPage(db *gorm.DB, siteID uint, parentID *uint, slug, title, contentMD string, sort int) (uint, error) {
	p := model.Page{
		SiteID:    siteID,
		ParentID:  parentID,
		Slug:      slug,
		Title:     title,
		ContentMD: contentMD,
		Sort:      sort,
	}
	if err := db.Create(&p).Error; err != nil {
		return 0, err
	}
	return p.ID, nil
}

// InsertMedia 插入媒体（如已存在 hash 则复用）
func InsertMedia(db *gorm.DB, siteID uint, hash, filename, mime string, data []byte) error {
	var m model.Media
	err := db.Where("hash = ?", hash).First(&m).Error
	if err == nil {
		return nil // 已存在
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	m = model.Media{
		SiteID:   siteID,
		Hash:     hash,
		Filename: filename,
		Mime:     mime,
		Data:     data,
	}
	return db.Create(&m).Error
}
