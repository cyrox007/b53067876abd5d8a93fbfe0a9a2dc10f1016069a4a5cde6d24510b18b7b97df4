package handlers

import (
	"time"

	"test/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// LoginHandler обрабатывает запрос на получение JWT-токена
func LoginHandler(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.WithError(err).Warn("Ошибка парсинга тела запроса")
		return c.Status(400).JSON(fiber.Map{
			"Success": false,
			"Message": "Bad Request: Invalid JSON",
		})
	}

	// Получаем тестовые учетные данные из переменных окружения
	testUsername := viper.GetString("TEST_USERNAME")
	testPassword := viper.GetString("TEST_PASSWORD")

	// Проверяем учетные данные
	if req.Username != testUsername || req.Password != testPassword {
		logger.Logger.Warn("Неверные учетные данные")
		return c.Status(401).JSON(fiber.Map{
			"Success": false,
			"Message": "Unauthorized: Invalid credentials",
		})
	}

	// Получаем секретный ключ из переменной окружения
	jwtSecret := viper.GetString("JWT_SECRET")
	if jwtSecret == "" {
		logger.Logger.Error("JWT_SECRET не задан в переменных окружения")
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Internal Server Error: JWT_SECRET is not set",
		})
	}

	// Создаем JWT-токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	})

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка создания JWT-токена")
		return c.Status(500).JSON(fiber.Map{
			"Success": false,
			"Message": "Internal Server Error: Failed to generate token",
		})
	}

	logger.Logger.Info("JWT-токен успешно создан")
	return c.JSON(fiber.Map{
		"Success": true,
		"Token":   tokenString,
	})
}
