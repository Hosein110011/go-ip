package config

import (
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Connect() {
	d, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=postgres password=h55555H55555h dbname=go port=5432 sslmode=disable TimeZone=Asia/Tehran",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}
