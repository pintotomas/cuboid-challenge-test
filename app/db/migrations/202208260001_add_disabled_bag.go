package migrations

import (
	"cuboid-challenge/app/models"
	"fmt"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "202208260001",
		Migrate: func(transaction *gorm.DB) error {
			fmt.Println("Running migration add_disabled_bag")
			type Bag struct {
				models.Model
				Title    string
				Volume   uint
				Disabled bool
			}

			return transaction.AutoMigrate(&Bag{})
		},
		Rollback: func(transaction *gorm.DB) error {
			fmt.Println("Rollback migration add_disabled_bag")
			type Bag struct {
				models.Model
				Title  string
				Volume uint
			}

			return transaction.Migrator().DropTable(&Bag{})
		},
	})
}
