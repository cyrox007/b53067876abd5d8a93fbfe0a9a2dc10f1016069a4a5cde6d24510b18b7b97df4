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

	protected.Post("/create", handlers.CreateNews)       // Создание новости
	protected.Post("/edit/:Id", handlers.EditNews)       // Редактирование новости
	protected.Delete("/delete/:Id", handlers.DeleteNews) // Удаление новости
	protected.Get("/list", handlers.GetNewsList)
}
