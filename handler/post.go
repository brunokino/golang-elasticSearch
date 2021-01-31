package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/idoko/letterpress/db"
	"gitlab.com/idoko/letterpress/models"
	"net/http"
	"strconv"
)

func (h *Handler) GetPosts(c *gin.Context) {
	posts, err := h.DB.GetPosts()
	if err != nil {
		h.Logger.Err(err).Msg("Could not fetch posts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func (h *Handler) GetPost(c *gin.Context) {
	var id int
	var err error
	var post models.Post
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
	}
	post, err = h.DB.GetPostById(id)
	switch err {
	case db.ErrNoRecord:
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find post with id: %d", id)})
		break
	case nil:
		c.JSON(http.StatusOK, gin.H{"data": post})
		break
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}