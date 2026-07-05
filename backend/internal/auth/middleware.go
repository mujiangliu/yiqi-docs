// backend/internal/auth/middleware.go
package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"jiaocheng-web/backend/internal/model"
)

const ContextUserKey = "current_user"

// Middleware 返回 Gin 中间件：解析 Cookie 中的 JWT，注入 user 到 context
func Middleware(issuer *JWTIssuer) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""
		// 优先从 Cookie 读
		if ck, err := c.Cookie("token"); err == nil && ck != "" {
			tokenStr = ck
		} else if auth := c.GetHeader("Authorization"); strings.HasPrefix(auth, "Bearer ") {
			// 兼容 Bearer token（便于 API 测试）
			tokenStr = strings.TrimPrefix(auth, "Bearer ")
		}
		if tokenStr == "" {
			c.Next()
			return
		}
		claims, err := issuer.Parse(tokenStr)
		if err != nil {
			c.Next()
			return
		}
		c.Set(ContextUserKey, &model.User{ID: claims.UID, Role: claims.Role})
		c.Next()
	}
}

// RequireAuth 要求已登录，否则 401
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := c.Get(ContextUserKey)
		if !ok || u == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

// RequireRole 要求特定角色
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		u, ok := c.Get(ContextUserKey)
		if !ok || u == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		user := u.(*model.User)
		if !allowed[user.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

// CurrentUser 从 context 取出当前用户（可能为 nil）
func CurrentUser(c *gin.Context) *model.User {
	v, ok := c.Get(ContextUserKey)
	if !ok || v == nil {
		return nil
	}
	return v.(*model.User)
}
