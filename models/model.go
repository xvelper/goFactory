package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
}

type Repository struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
	Path string
}

func InitDB() *gorm.DB {
	dsn := "host=localhost user=go password=198771Ll dbname=goFactory port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{}, &Repository{})
	return db
}
