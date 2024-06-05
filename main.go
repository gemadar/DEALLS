package main

import (
	"echo-mongo-api/configs"
	"echo-mongo-api/controllers"
	"echo-mongo-api/routes"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// Get ClearViewTime from .env file
func ClearViewTime() time.Duration {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	res, _ := time.ParseDuration(os.Getenv("ClearView"))

	return res
}

func main() {
	e := echo.New()

	//run database
	configs.DBConn()

	//routes
	routes.UserRoute(e)
	routes.AuthRoute(e)

	// Clear Expired Logged In User
	controllers.Schedule(controllers.ClearView, ClearViewTime()*time.Hour)

	e.Logger.Fatal(e.Start(":6161"))
}
