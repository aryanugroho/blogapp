package migration

import "github.com/go-gormigrate/gormigrate/v2"

var migrations []*gormigrate.Migration

func GetMigrations() []*gormigrate.Migration {
	return migrations
}
