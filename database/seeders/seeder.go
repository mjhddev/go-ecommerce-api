package seeders

import "gorm.io/gorm"

func Seed(db *gorm.DB) error {
	if err := seedUsers(db); err != nil {
		return err
	}

	if err := seedCategories(db); err != nil {
		return err
	}

	if err := seedProducts(db); err != nil {
		return err
	}

	return nil
}
