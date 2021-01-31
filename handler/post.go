package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetPosts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}