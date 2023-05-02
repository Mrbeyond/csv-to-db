package model

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

// Current database instance
var DB *gorm.DB

// Establishes Database connection and migrations
// @instanceType  indicates the type of DB to use
func DbConfig(instanceType string) {
	var (
		db  *gorm.DB
		err error
	)
	if instanceType == "LIVE_CONNECTION" {
		db, err = PostgressInstance()
	} else {
		db, err = SQLiteInstance()
	}

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

func PostgressInstance() (*gorm.DB, error) {

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbname, port)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
}

func SQLiteInstance() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		PrepareStmt: true,
	})
}
