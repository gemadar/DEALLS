package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Get MongoDB URI from .env file
func MongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGOURI")
}
