package db

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"gitlab.com/idoko/letterpress/models"
	"strconv"
)

var (
	ErrNoRecord = fmt.Errorf("no matching record found")
	elasticPostIndex = "posts"
)

func (db Database) SavePost(post *models.Post) error {
	var id int
	query := `INSERT INTO posts(title, body) VALUES ($1, $2) RETURNING id`
	err := db.Conn.QueryRow(query, post.Title, post.Body).Scan(&id)
	if err != nil {
		return err
	}
	post.ID = id

	// dereference post since we have no plan to mutate it in the indexPost function
	if res, err := indexPost(db.esClient, *post); err != nil {
		db.Logger.Err(err).Msg(fmt.Sprintf("could not index document ID=%d", post.ID))
		return err
	} else {
		db.Logger.Info().Msg(fmt.Sprintf("[%s] index document ID=%d", res.Status(), post.ID))
	}

	return nil
}

func indexPost(esClient *elasticsearch.Client, post models.Post) (*esapi.Response, error) {
	body, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}
	request := esapi.IndexRequest{
		Index: elasticPostIndex,
		DocumentID: strconv.Itoa(post.ID),
		Refresh: "true",
		Body: bytes.NewBuffer(body),
	}
	return request.Do(context.Background(), esClient)
}

func deleteFromIndex(esClient *elasticsearch.Config, post models.Post) (*esapi.Response, error) {

}

func (db Database) UpdatePost(post models.Post) error {
	query := "UPDATE posts SET title=$1, body=$2 WHERE id=$3"
	_, err := db.Conn.Exec(query, post.Title, post.Body, post.ID)
	if err != nil {
		return err
	}
	if res, err := indexPost(db.esClient, post); err != nil {
		db.Logger.Err(err).Msg(fmt.Sprintf("could not update document ID=%d", post.ID))
		return err
	} else {
		db.Logger.Info().Msg(fmt.Sprintf("[%s] updated index for document ID=%d", res.Status(), post.ID))
	}
	return nil
}

func (db Database) DeletePost(post models.Post) error {
	query := ""
	_, err := db.Conn.Exec(query, post.ID)
	if err != nil {
		return err
	}
	if res, err := deleteFromIndex(db.esClient, post); err != nil {
		db.Logger.Err(err).Msg(fmt.Sprintf("could not delete document ID=%d from index", post.ID))
		return err
	} else {
		db.Logger.Info().Msg(fmt.Sprintf("[%s] deleted document ID=%d from index", res.Status(), post.ID))
	}
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

