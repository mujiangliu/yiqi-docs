// backend/internal/api/responses.go
package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"jiaocheng-web/backend/internal/store"
)

// OK 返回 200 + data
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Created 返回 201 + data
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

// Fail 返回错误。store.ErrNotFound → 404；其它 → 400/500
func Fail(c *gin.Context, err error) {
	if errors.Is(err, store.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

// FailMsg 返回指定消息
func FailMsg(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}
