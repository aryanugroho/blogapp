package model

import "context"

type BlogSQLStore interface {
	Create(ctx context.Context, model *Post) (*Post, error)
	Update(ctx context.Context, model *Post) (*Post, error)
	Delete(ctx context.Context, model *Post) error
	FindByID(ctx context.Context, id string) (*Post, error)
	FindAll(ctx context.Context) ([]*Post, error)
}

type CommentSQLStore interface {
	Create(ctx context.Context, model *Comment) (*Comment, error)
	Update(ctx context.Context, model *Comment) (*Comment, error)
	Delete(ctx context.Context, model *Comment) error
	FindByID(ctx context.Context, id string) (*Comment, error)
	FindByPostID(ctx context.Context, postID string) (*Comment, error)
}
