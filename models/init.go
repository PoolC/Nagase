package models

import "nagase/components/database"

func init() {
	database.DB.AutoMigrate(
		&Board{},
		&Member{},
		&Post{},
		&Comment{},
		&Vote{},
		&VoteOption{},
		&VoteSelection{},
		&Project{},
	)
}
