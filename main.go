package main

import (
	"log"

	"test/database"
	"test/logger"
	"test/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Инициализация логгера
	logger.InitLogger()

	// Подключение к базе данных
	if err := database.Connect(); err != nil {
		logger.Logger.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Выполнение миграций
	if err := database.Migrate(); err != nil {
		logger.Logger.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	app.Use(recover.New())

	// Регистрация маршрутов
	routes.RegisterAuthRoutes(app) // Маршруты аутентификации
	routes.RegisterNewsRoutes(app) // Защищенные маршруты

	logger.Logger.Info("Приложение запущено")
	log.Fatal(app.Listen(":9000"))
}
