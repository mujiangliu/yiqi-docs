// backend/internal/store/users.go
package store

import (
	"errors"

	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/model"
)

// CreateUser 创建用户。仅 super_admin 可调用（API 层校验，此处不重复）
func (s *Store) CreateUser(u model.User) (*model.User, error) {
	err := s.db.Create(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByUsername 按用户名查询（登录用）
func (s *Store) GetUserByUsername(username string) (*model.User, error) {
	var u model.User
	err := s.db.Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// GetUserByID 按 ID 查询
func (s *Store) GetUserByID(id uint) (*model.User, error) {
	var u model.User
	err := s.db.First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// ListUsers 列出所有用户（仅 super_admin）
func (s *Store) ListUsers() ([]model.User, error) {
	var users []model.User
	err := s.db.Order("id asc").Find(&users).Error
	return users, err
}

// UpdateUser 更新指定字段
func (s *Store) UpdateUser(id uint, fields map[string]interface{}) (*model.User, error) {
	err := s.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return nil, err
	}
	return s.GetUserByID(id)
}

// CountSuperAdmin 统计 super_admin 数量（防止删到最后一个）
func (s *Store) CountSuperAdmin() (int64, error) {
	var n int64
	err := s.db.Model(&model.User{}).Where("role = ?", model.RoleSuperAdmin).Count(&n).Error
	return n, err
}

// DeleteUser 删除用户
func (s *Store) DeleteUser(id uint) error {
	return s.db.Delete(&model.User{}, id).Error
}
