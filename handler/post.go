package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/idoko/letterpress/db"
	"gitlab.com/idoko/letterpress/models"
	"net/http"
	"strconv"
)

func (h *Handler) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		h.Logger.Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %s", err.Error())})
		return
	}
	err := h.DB.SavePost(&post)
	if err != nil {
		h.Logger.Err(err).Msg("could not save post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not save post: %s", err.Error())})
	} else {
		c.JSON(http.StatusCreated, gin.H{"post": post})
	}
}

func (h *Handler) UpdatePost(c *gin.Context) {
	var post models.Post
	var err error
	if err = c.ShouldBindJSON(&post); err != nil {
		h.Logger.Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %s", err.Error())})
		return
	}
	err = h.DB.UpdatePost(post)
	if err != nil {
		h.Logger.Err(err).Msg("could not update post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not update post: %s", err.Error())})
	} else {
		c.JSON(http.StatusOK, gin.H{"post": post})
	}
}

func (h *Handler) DeletePost(c *gin.Context) {
	var id int
	var err error
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}
	err = h.DB.DeletePost(id)
	switch err {
	case db.ErrNoRecord:
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find post with id: %d", id)})
		break
	case nil:
		c.JSON(http.StatusOK, gin.H{"data": map[string]string{"message": "post deleted"}})
		break
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		break
	}
}

func (h *Handler) GetPosts(c *gin.Context) {
	posts, err := h.DB.GetPosts()
	if err != nil {
		h.Logger.Err(err).Msg("Could not fetch posts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": posts})
	}
}

func (h *Handler) GetPost(c *gin.Context) {
	var id int
	var err error
	var post models.Post
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
	} else {
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
}
