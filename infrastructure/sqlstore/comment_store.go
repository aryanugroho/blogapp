package sqlstore

import (
	"context"

	"github.com/aryanugroho/blogapp/internal/db"
	"github.com/aryanugroho/blogapp/model"
	"gorm.io/gorm"
)

// implementation of model store
type CommentStore struct {
	master *db.GormDBWrapper
	slave  *db.GormDBWrapper
}

func NewCommentStore(master *db.GormDBWrapper, slave *db.GormDBWrapper) *BlogStore {
	return &BlogStore{
		master: master,
		slave:  slave,
	}
}

func (u *CommentStore) Create(ctx context.Context, model *model.Post) (*model.Post, error) {
	dbs := u.master
	tx, ok := ctx.Value(Tx).(*db.GormDBWrapper)
	if ok {
		dbs = tx
	}
	err := dbs.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (u *CommentStore) Update(ctx context.Context, model *model.Post) (*model.Post, error) {
	dbs := u.master
	tx, ok := ctx.Value(Tx).(*db.GormDBWrapper)
	if ok {
		dbs = tx
	}

	err := dbs.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Model(model).Updates(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (u *CommentStore) Delete(ctx context.Context, model *model.Post) error {
	dbs := u.master
	tx, ok := ctx.Value(Tx).(*db.GormDBWrapper)
	if ok {
		dbs = tx
	}

	err := dbs.WithContext(ctx).Delete(model).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *CommentStore) FindByID(ctx context.Context, id string) (*model.Post, error) {
	dbs := u.slave
	tx, ok := ctx.Value(Tx).(*db.GormDBWrapper)
	if ok {
		dbs = tx
	}

	var model *model.Post
	err := dbs.WithContext(ctx).First(model, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (u *CommentStore) FindByPostID(ctx context.Context, postID string) ([]*model.Post, error) {
	dbs := u.slave
	tx, ok := ctx.Value(Tx).(*db.GormDBWrapper)
	if ok {
		dbs = tx
	}

	var model []*model.Post
	err := dbs.WithContext(ctx).First(model, "post_id = ?", postID).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}
