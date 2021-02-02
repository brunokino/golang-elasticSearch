package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gitlab.com/idoko/letterpress/db"
)

type Handler struct {
	DB db.Database
	Logger zerolog.Logger
}

func New(database db.Database, logger zerolog.Logger) *Handler {
	return &Handler {
		DB: database,
		Logger: logger,
	}
}

func (h *Handler) Register(group *gin.RouterGroup) {
	group.GET("/posts", h.GetPosts)
	group.POST("/posts", h.CreatePost)
	group.GET("/posts/:id", h.GetPost)
}