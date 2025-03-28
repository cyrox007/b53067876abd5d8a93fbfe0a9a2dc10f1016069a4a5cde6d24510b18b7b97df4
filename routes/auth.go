package routes

import (
	"test/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/login", handlers.LoginHandler)
}
