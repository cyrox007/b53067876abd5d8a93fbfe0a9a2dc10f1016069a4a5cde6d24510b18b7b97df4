package handlers

import (
	"database/sql"
	"strconv"
	"test/database"
	"test/logger"
	"test/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func EditNews(c *fiber.Ctx) error {
	// Извлекаем параметр Id из маршрута
	newsIDStr := c.Params("Id")

	logger.Logger.WithField("news_id", newsIDStr).Debug("Получен ID новости для редактирования")

	// Преобразуем строку в uint
	newsIDUint64, err := strconv.ParseUint(newsIDStr, 10, 64)
	if err != nil {
		logger.Logger.WithError(err).Warn("Неверный формат ID новости")
		return c.Status(400).JSON(fiber.Map{
			"Success": false,
			"Message": "Неверный формат ID новости",
		})
	}

	// Преобразуем uint64 в uint
	newsID := uint(newsIDUint64)

	var req models.NewsResponse

	// Парсим тело запроса
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.WithError(err).Warn("Ошибка парсинга тела запроса")
		return c.Status(400).JSON(fiber.Map{
			"Success": false,
			"Message": "Неверный формат запроса",
		})
	}

	logger.Logger.WithFields(logrus.Fields{
		"title":      req.Title,
		"content":    req.Content,
		"categories": req.Categories,
	}).Info("Тело запроса успешно распарсено")

	// Валидация полей
	if req.Title == "" || req.Content == "" {
		logger.Logger.Warn("Заголовок или содержимое новости пустые")
		return c.Status(400).JSON(fiber.Map{
			"Success": false,
			"Message": "Заголовок и содержимое обязательны",
		})
	}

	// Начинаем транзакцию
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.WithField("error", r).Error("Произошла ошибка, выполняется откат транзакции")
			tx.Rollback()
		}
	}()

	// Обновляем новость
	result := tx.Model(&models.News{}).Where("id = ?", newsID).Updates(models.News{
		Title:   req.Title,
		Content: req.Content,
	})
	if result.Error != nil {
		logger.Logger.WithError(result.Error).Error("Ошибка обновления новости")
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Ошибка обновления новости",
		})
	}

	logger.Logger.WithField("news_id", newsID).Info("Новость успешно обновлена")

	// Удаляем старые категории
	if err := tx.Where("news_id = ?", newsID).Delete(&models.NewsCategory{}).Error; err != nil {
		logger.Logger.WithError(err).Error("Ошибка удаления старых категорий")
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Ошибка удаления категорий",
		})
	}

	logger.Logger.WithField("news_id", newsID).Info("Старые категории успешно удалены")

	// Добавляем новые категории
	for _, categoryID := range req.Categories {
		if err := tx.Create(&models.NewsCategory{
			NewsId:     newsID,
			CategoryId: categoryID,
		}).Error; err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"news_id":     newsID,
				"category_id": categoryID,
				"error":       err,
			}).Error("Ошибка сохранения категории")
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"Success": false,
				"Message": "Ошибка сохранения категории",
			})
		}
	}

	logger.Logger.WithField("news_id", newsID).Info("Новые категории успешно добавлены")

	// Фиксируем транзакцию
	if err := tx.Commit().Error; err != nil {
		logger.Logger.WithError(err).Error("Ошибка фиксации транзакции")
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Ошибка фиксации транзакции",
		})
	}

	logger.Logger.WithField("news_id", newsID).Info("Транзакция успешно зафиксирована")
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

	logger.Logger.WithFields(logrus.Fields{
		"page":  page,
		"limit": limit,
	}).Info("Запрос списка новостей")

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
		logger.Logger.Errorf("Ошибка выполнения запроса к базе данных: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Ошибка выполнения запроса к базе данных",
		})
	}
	defer rows.Close()

	// Карта для группировки новостей и их категорий
	newsMap := make(map[uint]*models.NewsResponse)

	for rows.Next() {
		var newsID uint
		var title, content string
		var categoryID sql.NullInt64

		err := rows.Scan(&newsID, &title, &content, &categoryID)
		if err != nil {
			logger.Logger.Errorf("Ошибка сканирования данных: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"Success": false,
				"Message": "Ошибка сканирования данных",
			})
		}

		// Если новость еще не добавлена в карту, создаем запись
		if _, exists := newsMap[newsID]; !exists {
			newsMap[newsID] = &models.NewsResponse{
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
	var newsList []models.NewsResponse
	for _, news := range newsMap {
		newsList = append(newsList, *news)
	}

	logger.Logger.WithField("count", len(newsList)).Info("Новости успешно получены")
	return c.JSON(fiber.Map{
		"Success": true,
		"News":    newsList,
	})
}
