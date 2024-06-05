package routes

import (
	"echo-mongo-api/controllers"

	"github.com/labstack/echo/v4"
)

// Route for auth endpoint
func AuthRoute(e *echo.Echo) {
	e.POST("/signup", controllers.SignUp)
	e.POST("/signin", controllers.SignIn)
}
