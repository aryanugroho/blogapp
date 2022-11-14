package app

import (
	"context"

	"github.com/aryanugroho/blogapp/model"
)

func (u *App) CreatePost(ctx context.Context, payload *PostPayload) (*Post, error) {
	post, err := u.Infrastructure.SQLStore().BlogStore().Create(ctx, &model.Post{
		Title:   payload.Title,
		Content: payload.Content,
	})
	if err != nil {
		return nil, err
	}
	return &Post{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		AuditableEntity: post.AuditableEntity,
	}, nil
}
