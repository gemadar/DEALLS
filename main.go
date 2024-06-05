package main

import (
	"echo-mongo-api/configs"
	"echo-mongo-api/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	//run database
	configs.DBConn()

	//routes
	routes.UserRoute(e)
	routes.AuthRoute(e)

	e.Logger.Fatal(e.Start(":6161"))
}
