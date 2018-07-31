package handlers

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/nicolaszein/go-retro/models"
)

var (
	testDB    *gorm.DB
	connError error
	env       Env
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

	env = Env{
		DB: testDB,
	}
}

func cleanDatabase(db *gorm.DB) {
	db.Unscoped().Delete(&models.Team{})
}
