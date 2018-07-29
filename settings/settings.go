package settings

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT         string
	DATABASE_URL string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT = os.Getenv("PORT")
}
