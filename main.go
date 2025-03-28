package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"test/database"
	"test/logger"
	"test/routes"
)

func main() {
	logger.InitLogger()

	if err := database.Connect(); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	app.Use(compress.New())
	app.Use(recover.New())
	app.Use(limiter.New())

	routes.RegisterProductRoutes(app)

	logger.Logger.Info("Приложение запущено")
	log.Fatal(app.Listen(":9000"))
}
