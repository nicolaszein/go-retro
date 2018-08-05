package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var (
	testDB    *gorm.DB
	connError error
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("Error while trying to load .env file: ", connError)
		panic(".env load error!")
	}

	testDB, connError = gorm.Open("postgres", os.Getenv("TEST_DATABASE_URL"))

	if connError != nil {
		fmt.Println("Error while trying to connect with test database: ", connError)
		panic("Database Error!")
	}
}

func cleanDatabase(db *gorm.DB) {
	db.Unscoped().Delete(&Retrospective{})
	db.Unscoped().Delete(&Team{})
}
