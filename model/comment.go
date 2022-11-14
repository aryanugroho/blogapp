package model

type Comment struct {
	ID      string
	UUID    string
	PostID  string
	Content string
	AuditableEntity
}
