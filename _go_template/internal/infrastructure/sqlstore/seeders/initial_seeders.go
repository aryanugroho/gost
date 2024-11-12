package seeders

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	seeders = append(seeders, &gormigrate.Migration{
		ID: "initial_seeders",
		Migrate: func(tx *gorm.DB) error {

			return nil
		},
		Rollback: func(tx *gorm.DB) error {

			return nil
		},
	})
}
