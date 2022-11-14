package model

type Post struct {
	ID      string
	UUID    string
	Title   string
	Content string
	AuditableEntity
}
