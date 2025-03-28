package handlers

import (
	"test/database"
	"test/models"

	"github.com/gofiber/fiber/v2"
)

func GetNews(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT id, name, description, price, stock, image_url FROM products")
	if err != nil {
		return c.Status(500).SendString("Ошибка выполнения запроса к базе данных")
	}
	defer rows.Close()

	var products []models.News
	for rows.Next() {
		var product models.News
		err := rows.Scan(&product.ID, &product.Name, &product.Description)
		if err != nil {
			return c.Status(500).SendString("Ошибка сканирования данных")
		}
		products = append(products, product)
	}

	return c.JSON(products)
}

func CreateNews(c *fiber.Ctx) error {
	news := new(models.News)
	if err := c.BodyParser(news); err != nil {
		return c.Status(400).SendString("Неверный формат запроса")
	}

	_, err := database.DB.Exec("INSERT INTO news (name, description) VALUES ($1, $2)",
		news.Name, news.Description)
	if err != nil {
		return c.Status(500).SendString("Ошибка вставки данных в базу")
	}

	return c.Status(201).SendString("Продукт успешно создан")
}
