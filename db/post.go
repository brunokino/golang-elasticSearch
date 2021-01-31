package db

import (
	"database/sql"
	"fmt"
	"gitlab.com/idoko/letterpress/models"
)

var ErrNoRecord = fmt.Errorf("no matching record found")

func (db Database) SavePost(post *models.Post) error {
	var id int
	query := `INSERT INTO posts(title, body) VALUES ($1, $2) RETURNING id`
	err := db.Conn.QueryRow(query, post.Title, post.Body).Scan(&id)
	if err != nil {
		return err
	}
	post.ID = id
	return nil
}

func (db Database) GetPostById(postId int) (models.Post, error) {
	post := models.Post{}
	query := "SELECT id, title, body FROM posts WHERE id = $1"
	row := db.Conn.QueryRow(query, postId)
	switch err := row.Scan(&post.ID, &post.Title, &post.Body); err {
	case sql.ErrNoRows:
		return post, ErrNoRecord
	default:
		return post, err
	}
}