// backend/internal/api/sites_admin.go
package api

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"jiaocheng-web/backend/internal/auth"
	"jiaocheng-web/backend/internal/model"
	"jiaocheng-web/backend/internal/store"
)

type AdminSiteHandler struct {
	store *store.Store
}

func NewAdminSiteHandler(s *store.Store) *AdminSiteHandler {
	return &AdminSiteHandler{store: s}
}

type createSiteReq struct {
	Path        string `json:"path" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// List GET /api/admin/sites
func (h *AdminSiteHandler) List(c *gin.Context) {
	u := auth.CurrentUser(c)
	s := h.store.ForUser(u)
	sites, err := s.ListSites()
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, sites)
}

// Create POST /api/admin/sites
func (h *AdminSiteHandler) Create(c *gin.Context) {
	var req createSiteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid request")
		return
	}
	u := auth.CurrentUser(c)
	s := h.store.ForUser(u)
	status := req.Status
	if status == "" {
		status = model.SiteStatusDraft
	}
	site, err := s.CreateSite(model.Site{
		OwnerID:     u.ID,
		Path:        req.Path,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
	})
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, site)
}

type updateSiteReq struct {
	Path        *string `json:"path"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

// Update PUT /api/admin/sites/:id
func (h *AdminSiteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req updateSiteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid request")
		return
	}
	fields := map[string]interface{}{}
	if req.Path != nil {
		fields["path"] = *req.Path
	}
	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.Description != nil {
		fields["description"] = *req.Description
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	u := auth.CurrentUser(c)
	s := h.store.ForUser(u)
	site, err := s.UpdateSite(uint(id), fields)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, site)
}

// Delete DELETE /api/admin/sites/:id
func (h *AdminSiteHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid id")
		return
	}
	u := auth.CurrentUser(c)
	s := h.store.ForUser(u)
	err = s.DeleteSite(uint(id))
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"ok": true})
}

// UploadMedia POST /api/admin/sites/:id/media
func (h *AdminSiteHandler) UploadMedia(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid id")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "file required")
		return
	}
	src, err := file.Open()
	if err != nil {
		Fail(c, err)
		return
	}
	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		Fail(c, err)
		return
	}
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])
	u := auth.CurrentUser(c)
	s := h.store.ForUser(u)
	media, err := s.CreateMedia(model.Media{
		SiteID:   uint(id),
		Hash:     hash,
		Filename: file.Filename,
		Mime:     file.Header.Get("Content-Type"),
		Data:     data,
	})
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"hash": media.Hash, "url": "/api/media/" + media.Hash})
}
