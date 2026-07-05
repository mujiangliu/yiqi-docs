// backend/internal/store/store.go
package store

import (
	"errors"

	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/model"
)

// ErrNotFound 资源不存在或无权访问（统一返回 404 以免泄露存在性）
var ErrNotFound = errors.New("record not found")

// Store 数据访问层。ForUser 返回带 owner 过滤的子 store。
type Store struct {
	db   *gorm.DB
	user *model.User // 可为 nil，表示公开视角
}

// NewStore 创建一个 store 根
func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// ForUser 返回绑定到某用户的 store。super_admin 不加 owner 过滤；admin 加 owner 过滤。
func (s *Store) ForUser(user *model.User) *Store {
	return &Store{db: s.db, user: user}
}

// Public 返回公开视角的 store（只查 published）
func (s *Store) Public() *Store {
	return &Store{db: s.db, user: nil}
}

// scopedDB 返回带 owner 过滤的 db（admin 限定 owner_id；super_admin 不过滤）
func (s *Store) scopedDB() *gorm.DB {
	if s.user == nil {
		return s.db
	}
	if s.user.Role == model.RoleSuperAdmin {
		return s.db
	}
	// admin: 强制按 owner 过滤（在 site 层面）
	return s.db.Where("owner_id = ?", s.user.ID)
}

// isSuper 是否总管
func (s *Store) isSuper() bool {
	return s.user != nil && s.user.Role == model.RoleSuperAdmin
}
