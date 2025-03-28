package routes

import (
	"test/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterProductRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/products", handlers.GetNews)
	api.Post("/products", handlers.CreateNews)
}
