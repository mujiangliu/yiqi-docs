// backend/internal/api/users_admin.go
package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"jiaocheng-web/backend/internal/auth"
	"jiaocheng-web/backend/internal/model"
	"jiaocheng-web/backend/internal/store"
)

type AdminUserHandler struct {
	store *store.Store
}

func NewAdminUserHandler(s *store.Store) *AdminUserHandler {
	return &AdminUserHandler{store: s}
}

type createUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// List GET /api/admin/users
func (h *AdminUserHandler) List(c *gin.Context) {
	users, err := h.store.ListUsers()
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, users)
}

// Create POST /api/admin/users
func (h *AdminUserHandler) Create(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid request")
		return
	}
	if req.Role != model.RoleAdmin && req.Role != model.RoleSuperAdmin {
		FailMsg(c, http.StatusBadRequest, "invalid role")
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		Fail(c, err)
		return
	}
	u, err := h.store.CreateUser(model.User{
		Username:     req.Username,
		PasswordHash: hash,
		Role:         req.Role,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			FailMsg(c, http.StatusConflict, "username exists")
			return
		}
		Fail(c, err)
		return
	}
	Created(c, u)
}

type updateUserReq struct {
	Role *string `json:"role"`
}

// Update PUT /api/admin/users/:id
func (h *AdminUserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req updateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid request")
		return
	}
	fields := map[string]interface{}{}
	if req.Role != nil {
		if *req.Role != model.RoleAdmin && *req.Role != model.RoleSuperAdmin {
			FailMsg(c, http.StatusBadRequest, "invalid role")
			return
		}
		fields["role"] = *req.Role
	}
	u, err := h.store.UpdateUser(uint(id), fields)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, u)
}

type resetPasswordReq struct {
	Password string `json:"password" binding:"required"`
}

// ResetPassword POST /api/admin/users/:id/reset-password
func (h *AdminUserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req resetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid request")
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		Fail(c, err)
		return
	}
	_, err = h.store.UpdateUser(uint(id), map[string]interface{}{"password_hash": hash})
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"ok": true})
}

// Delete DELETE /api/admin/users/:id
func (h *AdminUserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid id")
		return
	}
	current := auth.CurrentUser(c)
	if current.ID == uint(id) {
		FailMsg(c, http.StatusBadRequest, "cannot delete yourself")
		return
	}
	target, err := h.store.GetUserByID(uint(id))
	if err != nil {
		Fail(c, err)
		return
	}
	if target.Role == model.RoleSuperAdmin {
		count, err := h.store.CountSuperAdmin()
		if err != nil {
			Fail(c, err)
			return
		}
		if count <= 1 {
			FailMsg(c, http.StatusBadRequest, "cannot delete the last super_admin")
			return
		}
	}
	err = h.store.DeleteUser(uint(id))
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"ok": true})
}
