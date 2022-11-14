package app

import "github.com/aryanugroho/blogapp/model"

type PostPayload struct {
	UUID    string `json:"uuid"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Post struct {
	ID      string `json:"id"`
	UUID    string `json:"uuid"`
	Title   string `json:"title"`
	Content string `json:"content"`
	model.AuditableEntity
}

type CommentPayload struct {
	UUID    string `json:"uuid"`
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

type Comment struct {
	ID      string `json:"id"`
	UUID    string `json:"uuid"`
	PostID  string `json:"post_id"`
	Content string `json:"content"`
	model.AuditableEntity
}
