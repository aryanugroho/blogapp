package sqlstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aryanugroho/blogapp/internal/contextprop"
	"github.com/aryanugroho/blogapp/internal/db"
	"github.com/aryanugroho/blogapp/internal/logger"
	"github.com/aryanugroho/blogapp/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	Tx contextprop.ContextKey = "tx"
)

// Store is the wrapper for dto.
type Store interface {
	BeginTx(ctx context.Context) context.Context
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error

	BlogStore() model.BlogSQLStore
}

type SQLStore struct {
	db *db.GormDBWrapper

	blog model.BlogSQLStore
}

func NewSQLStore(ctx context.Context, dbConfigMaster, dbConfigSlave db.Config) (Store, error) {
	connectionTimeout := 3 * time.Second

	logger.Info(ctx,
		fmt.Sprintf("database configs master, max open %d, max idle %d, max lifetime %d",
			dbConfigMaster.MaxOpen, dbConfigMaster.MaxIdle, dbConfigMaster.MaxLifetime))

	logger.Info(ctx,
		fmt.Sprintf("database configs slave, max open %d, max idle %d, max lifetime %d",
			dbConfigSlave.MaxOpen, dbConfigSlave.MaxIdle, dbConfigSlave.MaxLifetime))

	sqlDBMaster, err := db.DB(dbConfigMaster, connectionTimeout)
	if err != nil {
		return nil, err
	}
	gormDBMaster, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDBMaster,
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDBSlave, err := db.DB(dbConfigMaster, connectionTimeout)
	if err != nil {
		return nil, err
	}
	gormDBSlave, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDBSlave,
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	masterWrapper := db.NewGormDBWrapper(gormDBMaster, dbConfigMaster.Driver, connectionTimeout)
	slaveWrapper := db.NewGormDBWrapper(gormDBSlave, dbConfigSlave.Driver, connectionTimeout)

	return &SQLStore{
		db:   masterWrapper,
		blog: NewBlogStore(masterWrapper, slaveWrapper),
	}, nil
}

func openDBWrapper(dbConfig db.Config) (*gorm.DB, error) {
	masterDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)
	sqlDB, err := gorm.Open(mysql.Open(masterDsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}
	return sqlDB, nil
}

func (s *SQLStore) BeginTx(ctx context.Context) context.Context {
	tx := s.db.Begin()
	ctx = context.WithValue(ctx, Tx, tx)
	return ctx
}

func (s *SQLStore) CommitTx(ctx context.Context) error {
	tx, ok := ctx.Value(Tx).(*gorm.DB)
	if !ok {
		return errors.New("failed to commit on non transaction mode")
	}

	return tx.Commit().Error
}

func (s *SQLStore) RollbackTx(ctx context.Context) error {
	tx, ok := ctx.Value(Tx).(*gorm.DB)
	if !ok {
		return errors.New("failed to rollback on non transaction mode")
	}
	tx.Rollback()
	return nil
}

func (s *SQLStore) BlogStore() model.BlogSQLStore {
	return s.blog
}
