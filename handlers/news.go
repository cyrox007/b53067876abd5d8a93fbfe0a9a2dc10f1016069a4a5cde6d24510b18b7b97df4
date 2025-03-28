package handlers

import (
	"database/sql"
	"strconv"
	"test/database"
	"test/models"

	"github.com/gofiber/fiber/v2"
)

func EditNews(c *fiber.Ctx) error {
	// Извлекаем параметр Id из маршрута
	newsIDStr := c.Params("Id")

	// Преобразуем строку в uint
	newsIDUint64, err := strconv.ParseUint(newsIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Success": false,
			"Message": "Неверный формат ID новости",
		})
	}

	// Преобразуем uint64 в uint
	newsID := uint(newsIDUint64)

	var req models.News

	// Парсим тело запроса
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Success": false,
			"Message": "Неверный формат запроса",
		})
	}

	// Валидация полей
	if req.Title == "" || req.Content == "" {
		return c.Status(400).JSON(fiber.Map{
			"Success": false,
			"Message": "Заголовок и содержимое обязательны",
		})
	}

	// Начинаем транзакцию
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Обновляем новость
	result := tx.Model(&models.News{}).Where("id = ?", newsID).Updates(models.News{
		Title:   req.Title,
		Content: req.Content,
	})
	if result.Error != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Ошибка обновления новости",
		})
	}

	// Удаляем старые категории
	if err := tx.Where("news_id = ?", newsID).Delete(&models.NewsCategory{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Ошибка удаления категорий",
		})
	}

	// Добавляем новые категории
	for _, categoryID := range req.Categories {
		if err := tx.Create(&models.NewsCategory{
			NewsId:     newsID, // Теперь это uint
			CategoryId: categoryID,
		}).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"Success": false,
				"Message": "Ошибка сохранения категории",
			})
		}
	}

	// Фиксируем транзакцию
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Ошибка фиксации транзакции",
		})
	}

	return c.JSON(fiber.Map{
		"Success": true,
		"Message": "Новость успешно обновлена",
	})
}

func GetNewsList(c *fiber.Ctx) error {
	// Пагинация
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	// Запрос для получения новостей и их категорий
	query := `
        SELECT n.id, n.title, n.content, nc.category_id
        FROM news n
        LEFT JOIN news_categories nc ON n.id = nc.news_id
        ORDER BY n.id
        LIMIT ? OFFSET ?
    `

	rows, err := database.DB.Raw(query, limit, offset).Rows()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Ошибка выполнения запроса к базе данных",
		})
	}
	defer rows.Close()

	// Карта для группировки новостей и их категорий
	newsMap := make(map[uint]*models.News)

	for rows.Next() {
		var newsID uint
		var title, content string
		var categoryID sql.NullInt64

		err := rows.Scan(&newsID, &title, &content, &categoryID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"Success": false,
				"Message": "Ошибка сканирования данных",
			})
		}

		// Если новость уже есть в карте, добавляем категорию
		if _, exists := newsMap[newsID]; !exists {
			newsMap[newsID] = &models.News{
				Id:         newsID,
				Title:      title,
				Content:    content,
				Categories: []uint{},
			}
		}

		// Добавляем категорию, если она существует
		if categoryID.Valid {
			newsMap[newsID].Categories = append(newsMap[newsID].Categories, uint(categoryID.Int64))
		}
	}

	// Преобразуем карту в массив новостей
	var newsList []models.News
	for _, news := range newsMap {
		newsList = append(newsList, *news)
	}

	// Возвращаем результат в требуемом формате
	return c.JSON(fiber.Map{
		"Success": true,
		"News":    newsList,
	})
}
