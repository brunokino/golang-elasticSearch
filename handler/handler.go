package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gitlab.com/idoko/letterpress/db"
)

type Handler struct {
	DB     db.Database
	Logger zerolog.Logger
}

func New(database db.Database, logger zerolog.Logger) *Handler {
	return &Handler{
		DB:     database,
		Logger: logger,
	}
}

func (h *Handler) Register(group *gin.RouterGroup) {
	group.GET("/posts/:id", h.GetPost)
	group.DELETE("/posts/:id", h.DeletePost)

	group.GET("/posts", h.GetPosts)
	group.POST("/posts", h.CreatePost)
	group.PUT("/posts", h.UpdatePost)
}
