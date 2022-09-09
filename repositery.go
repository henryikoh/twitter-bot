package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DAO struct {
	db *gorm.DB
}

func InitDAO() *DAO {

	db, err := gorm.Open(sqlite.Open("bot.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&accessToken{})

	DAO := &DAO{
		db: db,
	}

	return DAO
}
