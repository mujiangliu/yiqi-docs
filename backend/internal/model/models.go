// backend/internal/model/models.go
package model

import (
	"time"

	"gorm.io/gorm"
)

// RoleSuperAdmin 总管角色，可看所有内容、管理用户
const RoleSuperAdmin = "super_admin"

// RoleAdmin 内容管理员角色，只能看自己 owner 的内容
const RoleAdmin = "admin"

// SiteStatusPublished 已发布，对外可见
const SiteStatusPublished = "published"

// SiteStatusDraft 草稿，对外不可见
const SiteStatusDraft = "draft"

// User 用户表
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"not null;default:admin" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Site 站点表。owner_id 是权限隔离的关键字段
type Site struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OwnerID     uint           `gorm:"not null;index" json:"owner_id"`
	Path        string         `gorm:"uniqueIndex;not null" json:"path"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Status      string         `gorm:"not null;default:draft" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Page 页面表，树形结构（parent_id 自引用）
type Page struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	SiteID    uint           `gorm:"not null;index;uniqueIndex:uniq_page" json:"site_id"`
	ParentID  *uint          `gorm:"index;uniqueIndex:uniq_page" json:"parent_id"`
	Slug      string         `gorm:"not null;uniqueIndex:uniq_page" json:"slug"`
	Title     string         `gorm:"not null" json:"title"`
	Sort      int            `gorm:"default:0" json:"sort"`
	ContentMD string         `gorm:"type:text" json:"content_md"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Media 图片二进制本地存储，hash 唯一
type Media struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SiteID    uint      `gorm:"not null;index" json:"site_id"`
	Hash      string    `gorm:"uniqueIndex;not null" json:"hash"`
	Filename  string    `json:"filename"`
	Mime      string    `json:"mime"`
	Data      []byte    `gorm:"type:blob" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
