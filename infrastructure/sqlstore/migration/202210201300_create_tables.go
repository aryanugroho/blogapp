package migration

import (
	"github.com/aryanugroho/blogapp/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// CreateFraudCategoriesTable return migration object
func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "202210201300_create_posts_table",
		Migrate: func(tx *gorm.DB) error {

			type Post struct {
				ID      string
				UUID    string
				Title   string
				Content string
				model.AuditableEntity
			}

			return tx.AutoMigrate(&Post{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("posts")
		},
	})

	migrations = append(migrations, &gormigrate.Migration{
		ID: "202210201300_create_comments_table",
		Migrate: func(tx *gorm.DB) error {

			type Comment struct {
				ID      string
				UUID    string
				PostID  string
				Content string
				model.AuditableEntity
			}

			return tx.AutoMigrate(&Comment{})
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("comments")
		},
	})
}
