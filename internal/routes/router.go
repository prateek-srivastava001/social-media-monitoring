package routes

import "github.com/labstack/echo/v4"
"controllers"

func ScraperRoutes(incomingRoutes *echo.Echo) {
	incomingRoutes.GET("/search", controllers.TwitterScraper)
}
