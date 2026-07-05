// backend/internal/api/media_public.go
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMedia GET /api/media/:hash
func (h *PublicHandler) GetMedia(c *gin.Context) {
	hash := c.Param("hash")
	m, err := h.store.Public().GetMediaByHash(hash)
	if err != nil {
		Fail(c, err)
		return
	}
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Data(http.StatusOK, m.Mime, m.Data)
}
