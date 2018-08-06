package handlers

import (
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/nicolaszein/go-retro/database"
)

var (
	testDB    *database.Postgres
	connError error
	env       Env
	dbMock    *database.Mock
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("Error while trying to load .env file: ", connError)
		panic(".env load error!")
	}

	testDB, connError = database.NewPostgres(os.Getenv("TEST_DATABASE_URL"))
	if connError != nil {
		fmt.Println("Error while trying to connect with test database: ", connError)
		panic("Database Error!")
	}

	env = Env{
		DB: testDB,
	}

	dbMock = &database.Mock{}
}
