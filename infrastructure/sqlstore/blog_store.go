package sqlstore

import (
	"context"

	"github.com/aryanugroho/blogapp/internal/db"
	"github.com/aryanugroho/blogapp/model"
	"gorm.io/gorm"
)

// implementation of model store
type BlogStore struct {
	master *db.GormDBWrapper
	slave  *db.GormDBWrapper
}

func NewBlogStore(master *db.GormDBWrapper, slave *db.GormDBWrapper) *BlogStore {
	return &BlogStore{
		master: master,
		slave:  slave,
	}
}

func (u *BlogStore) Create(ctx context.Context, model *model.Post) (*model.Post, error) {
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

func (u *BlogStore) Update(ctx context.Context, model *model.Post) (*model.Post, error) {
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

func (u *BlogStore) Delete(ctx context.Context, model *model.Post) error {
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

func (u *BlogStore) FindByID(ctx context.Context, id string) (*model.Post, error) {
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

func (u *BlogStore) FindAll(ctx context.Context) ([]*model.Post, error) {
	dbs := u.slave
	tx, ok := ctx.Value(Tx).(*db.GormDBWrapper)
	if ok {
		dbs = tx
	}

	var model []*model.Post
	err := dbs.WithContext(ctx).Select("order by created_at DESC").Find(model).Error
	if err != nil {
		return nil, err
	}

	return model, nil
}
