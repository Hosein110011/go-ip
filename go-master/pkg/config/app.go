package config

import (
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Connect() {
	erre := godotenv.Load(".env") // Load .env file
	if erre != nil {
		log.Fatal("Error loading .env file", erre)
	}
	d, err := gorm.Open(postgres.New(postgres.Config{
		DSN: os.Getenv("DSN"),                 
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
