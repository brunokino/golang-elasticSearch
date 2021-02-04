package db

import (
	"database/sql"
	"fmt"
	"gitlab.com/idoko/letterpress/models"
)

var (
	ErrNoRecord = fmt.Errorf("no matching record found")
)

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

func (db Database) UpdatePost(post models.Post) error {
	query := "UPDATE posts SET title=$1, body=$2 WHERE id=$3"
	_, err := db.Conn.Exec(query, post.Title, post.Body, post.ID)
	return err
}

func (db Database) DeletePost(postId int) error {
	query := "DELETE FROM Posts WHERE id=$1"
	_, err := db.Conn.Exec(query, postId)
	return err
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

func (db Database) GetPosts() ([]models.Post, error) {
	var list []models.Post
	query := "SELECT id, title, body FROM posts ORDER BY id DESC"
	rows, err := db.Conn.Query(query)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Body)
		if err != nil {
			return list, err
		}
		list = append(list, post)
	}
	return list, nil
}
