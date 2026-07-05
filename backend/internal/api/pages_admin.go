// backend/internal/api/pages_admin.go
package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"jiaocheng-web/backend/internal/auth"
	"jiaocheng-web/backend/internal/model"
	"jiaocheng-web/backend/internal/store"
)

type AdminPageHandler struct {
	store *store.Store
}

func NewAdminPageHandler(s *store.Store) *AdminPageHandler {
	return &AdminPageHandler{store: s}
}

type createPageReq struct {
	ParentID  *uint  `json:"parent_id"`
	Slug      string `json:"slug" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Sort      int    `json:"sort"`
	ContentMD string `json:"content_md"`
}

// ListBySite GET /api/admin/sites/:id/pages
func (h *AdminPageHandler) ListBySite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid id")
		return
	}
	u := auth.CurrentUser(c)
	s := h.store.ForUser(u)
	pages, err := s.ListPagesBySite(uint(id))
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, pages)
}

// Create POST /api/admin/sites/:id/pages
func (h *AdminPageHandler) Create(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req createPageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid request")
		return
	}
	u := auth.CurrentUser(c)
	s := h.store.ForUser(u)
	page, err := s.CreatePage(model.Page{
		SiteID:    uint(id),
		ParentID:  req.ParentID,
		Slug:      req.Slug,
		Title:     req.Title,
		Sort:      req.Sort,
		ContentMD: req.ContentMD,
	})
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, page)
}

type updatePageReq struct {
	ParentID  *uint   `json:"parent_id"`
	Slug      *string `json:"slug"`
	Title     *string `json:"title"`
	Sort      *int    `json:"sort"`
	ContentMD *string `json:"content_md"`
}

// Update PUT /api/admin/pages/:id
func (h *AdminPageHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req updatePageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid request")
		return
	}
	fields := map[string]interface{}{}
	if req.ParentID != nil {
		fields["parent_id"] = *req.ParentID
	}
	if req.Slug != nil {
		fields["slug"] = *req.Slug
	}
	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.Sort != nil {
		fields["sort"] = *req.Sort
	}
	if req.ContentMD != nil {
		fields["content_md"] = *req.ContentMD
	}
	u := auth.CurrentUser(c)
	s := h.store.ForUser(u)
	page, err := s.UpdatePage(uint(id), fields)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, page)
}

// Delete DELETE /api/admin/pages/:id
func (h *AdminPageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid id")
		return
	}
	u := auth.CurrentUser(c)
	s := h.store.ForUser(u)
	err = s.DeletePage(uint(id))
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"ok": true})
}

type reorderReq struct {
	IDs []uint `json:"ids" binding:"required"`
}

// Reorder POST /api/admin/pages/:id/reorder
// 注意：:id 在此为占位（路由需要），实际用 body 中的 ids
func (h *AdminPageHandler) Reorder(c *gin.Context) {
	var req reorderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		FailMsg(c, http.StatusBadRequest, "invalid request")
		return
	}
	u := auth.CurrentUser(c)
	s := h.store.ForUser(u)
	// 校验所有页面都属于当前用户
	for _, id := range req.IDs {
		_, err := s.GetPageByID(id)
		if err != nil {
			Fail(c, err)
			return
		}
	}
	err := s.ReorderPages(req.IDs)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, gin.H{"ok": true})
}
