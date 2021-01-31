package models

type Post struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
}

