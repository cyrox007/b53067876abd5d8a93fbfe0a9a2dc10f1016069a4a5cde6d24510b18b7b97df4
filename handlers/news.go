package handlers

import (
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
	if err := tx.Model(&models.News{}).Where("id = ?", newsID).Updates(models.News{
		Title:   req.Title,
		Content: req.Content,
	}).Error; err != nil {
		logger.Logger.WithError(err).Error("Ошибка обновления новости")
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

	// Инициализируем массив новостей
	var newsList []models.NewsResponse

	// Выполняем запрос с использованием GORM
	err := database.DB.Model(&models.News{}).
		Select("news.id, news.title, news.content, news_categories.category_id").
		Joins("LEFT JOIN news_categories ON news.id = news_categories.news_id").
		Order("news.id").
		Limit(limit).
		Offset(offset).
		Scan(&newsList).Error

	if err != nil {
		logger.Logger.Errorf("Ошибка выполнения запроса к базе данных: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Ошибка выполнения запроса к базе данных",
		})
	}

	// Группируем данные по ID новости
	newsMap := make(map[uint]*models.NewsResponse)
	for _, news := range newsList {
		if _, exists := newsMap[news.Id]; !exists {
			newsMap[news.Id] = &models.NewsResponse{
				Id:         news.Id,
				Title:      news.Title,
				Content:    news.Content,
				Categories: []uint{},
			}
		}
		if news.Categories != nil {
			newsMap[news.Id].Categories = append(newsMap[news.Id].Categories, news.Categories...)
		}
	}

	// Преобразуем карту в массив новостей
	var result []models.NewsResponse
	for _, news := range newsMap {
		result = append(result, *news)
	}

	logger.Logger.WithField("count", len(result)).Info("Новости успешно получены")
	return c.JSON(fiber.Map{
		"Success": true,
		"News":    result,
	})
}
