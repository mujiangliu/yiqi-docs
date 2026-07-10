// backend/internal/api/sites_public.go
package api

import (
	"github.com/gin-gonic/gin"

	"jiaocheng-web/backend/internal/model"
	"jiaocheng-web/backend/internal/store"
)

type PublicHandler struct {
	store *store.Store
}

func NewPublicHandler(s *store.Store) *PublicHandler {
	return &PublicHandler{store: s}
}

type pageNode struct {
	ID        uint   `json:"id"`
	ParentID  *uint  `json:"parent_id"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Sort      int    `json:"sort"`
	ContentMD string `json:"content_md,omitempty"`
	Path      string `json:"path"`
}

// GetSite GET /api/sites/:path
func (h *PublicHandler) GetSite(c *gin.Context) {
	path := c.Param("path")
	pub := h.store.Public()
	site, err := pub.GetPublishedSiteByPath(path)
	if err != nil {
		Fail(c, err)
		return
	}
	pages, err := pub.ListPublishedPageTreeBySite(site.ID)
	if err != nil {
		Fail(c, err)
		return
	}
	idToPage := map[uint]model.Page{}
	for _, p := range pages {
		idToPage[p.ID] = p
	}
	nodes := make([]pageNode, 0, len(pages))
	for _, p := range pages {
		nodes = append(nodes, pageNode{
			ID:       p.ID,
			ParentID: p.ParentID,
			Slug:     p.Slug,
			Title:    p.Title,
			Sort:     p.Sort,
			Path:     buildPagePath(p, idToPage),
		})
	}
	OK(c, gin.H{
		"title":       site.Title,
		"description": site.Description,
		"pages":       nodes,
	})
}

// GetPage GET /api/sites/:path/pages/*pagePath
func (h *PublicHandler) GetPage(c *gin.Context) {
	path := c.Param("path")
	pagePath := c.Param("pagePath")
	if len(pagePath) > 0 && pagePath[0] == '/' {
		pagePath = pagePath[1:]
	}

	pub := h.store.Public()
	site, err := pub.GetPublishedSiteByPath(path)
	if err != nil {
		Fail(c, err)
		return
	}
	pages, err := pub.ListPublishedPageTreeBySite(site.ID)
	if err != nil {
		Fail(c, err)
		return
	}
	idToPage := map[uint]model.Page{}
	for _, p := range pages {
		idToPage[p.ID] = p
	}
	var match *model.Page
	for _, p := range pages {
		if buildPagePath(p, idToPage) == pagePath {
			page := p
			match = &page
			break
		}
	}
	if match == nil {
		Fail(c, store.ErrNotFound)
		return
	}
	page, err := pub.GetPublishedPageByID(site.ID, match.ID)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, pageNode{
		ID:        page.ID,
		ParentID:  page.ParentID,
		Slug:      page.Slug,
		Title:     page.Title,
		Sort:      page.Sort,
		ContentMD: page.ContentMD,
		Path:      pagePath,
	})
}

// buildPagePath 从根拼到当前节点
func buildPagePath(p model.Page, idToPage map[uint]model.Page) string {
	parts := []string{p.Slug}
	cur := p
	for cur.ParentID != nil {
		parent, ok := idToPage[*cur.ParentID]
		if !ok {
			break
		}
		parts = append([]string{parent.Slug}, parts...)
		cur = parent
	}
	return joinSlug(parts)
}

func joinSlug(parts []string) string {
	out := ""
	for i, s := range parts {
		if i > 0 {
			out += "/"
		}
		out += s
	}
	return out
}
