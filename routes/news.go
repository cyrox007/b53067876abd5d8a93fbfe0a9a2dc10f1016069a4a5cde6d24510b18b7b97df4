package routes

import (
	"test/handlers"
	"test/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterNewsRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Защищенные маршруты (требуют JWT-токен)
	protected := api.Group("/", middleware.AuthMiddleware)

	protected.Post("/edit/:Id", handlers.EditNews)
	protected.Get("/list", handlers.GetNewsList)
}
