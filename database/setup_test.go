package database

import (
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var (
	testDB    *Postgres
	connError error
	dbMock    *Mock
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("Error while trying to load .env file: ", connError)
		panic(".env load error!")
	}

	testDB, connError = NewPostgres(os.Getenv("TEST_DATABASE_URL"))
	if connError != nil {
		fmt.Println("Error while trying to connect with test database: ", connError)
		panic("Database Error!")
	}

	dbMock = &Mock{}
}
