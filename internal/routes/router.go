package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/prateek-srivastava001/social-media-monitoring/internal/controllers"
)

func ScraperRoutes(incomingRoutes *echo.Echo) {
	incomingRoutes.GET("/search", controllers.TwitterScraper)
}
