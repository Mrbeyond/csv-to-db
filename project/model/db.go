package model

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Current database instance
var DB *gorm.DB

// Establishes Database connection and migrations
func DbConfig() {
	var (
		dsn string
		err error
	)

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		fmt.Println(err)
		panic("Database connection error: " + err.Error())
	}

	err = db.AutoMigrate(
		&Ohcl{},
	)
	if err != nil {
		fmt.Println("Error from the migration", err.Error())
		panic("Error from the migration")
	}

	DB = db
}
