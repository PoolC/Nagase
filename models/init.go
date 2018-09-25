package models

import "nagase/components/database"

func init() {
	database.DB.AutoMigrate(
		&Member{},
	)
}
