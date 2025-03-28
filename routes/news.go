package routes

import (
	"test/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterProductRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/edit/:Id", handlers.EditNews)
	api.Get("/list", handlers.GetNewsList)
}
