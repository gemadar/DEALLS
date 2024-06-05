package routes

import (
	"echo-mongo-api/controllers"

	"github.com/labstack/echo/v4"
)

// Route for users endpoint
func UserRoute(e *echo.Echo) {
	e.GET("users/getpremium", controllers.GetPremium)
	e.GET("users/like", controllers.Like)
	e.GET("users/pass", controllers.Pass)
	e.GET("users/scroll", controllers.View)
}
