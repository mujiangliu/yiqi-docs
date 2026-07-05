// backend/internal/api/auth.go
package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"jiaocheng-web/backend/internal/auth"
	"jiaocheng-web/backend/internal/store"
)

type AuthHandler struct {
	store  *store.Store
	issuer *auth.JWTIssuer
}

func NewAuthHandler(s *store.Store, issuer *auth.JWTIssuer) *AuthHandler {
	return &AuthHandler{store: s, issuer: issuer}
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid request")
		return
	}
	user, err := h.store.GetUserByUsername(req.Username)
	if err != nil {
		FailMsg(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if !auth.VerifyPassword(user.PasswordHash, req.Password) {
		FailMsg(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, err := h.issuer.Issue(user.ID, user.Role, 7*24*time.Hour)
	if err != nil {
		FailMsg(c, http.StatusInternalServerError, "token issue failed")
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, int(7*24*time.Hour/time.Second), "/", "", false, true)
	OK(c, gin.H{"id": user.ID, "username": user.Username, "role": user.Role})
}

// Logout POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	OK(c, gin.H{"ok": true})
}

// Me GET /api/me
func (h *AuthHandler) Me(c *gin.Context) {
	u := auth.CurrentUser(c)
	if u == nil {
		FailMsg(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	full, err := h.store.GetUserByID(u.ID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			FailMsg(c, http.StatusUnauthorized, "user not found")
			return
		}
		Fail(c, err)
		return
	}
	OK(c, gin.H{"id": full.ID, "username": full.Username, "role": full.Role})
}
