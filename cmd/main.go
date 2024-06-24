package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prateek-srivastava001/social-media-monitoring/internal/routes"
)

func main() {
	app := echo.New()
	app.Use(middleware.Logger())

	app.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{
			"message": "pong",
			"status":  "successful start",
		})
	})

	routes.ScraperRoutes(app)
	app.Logger.Fatal(app.Start(":8080"))

}
