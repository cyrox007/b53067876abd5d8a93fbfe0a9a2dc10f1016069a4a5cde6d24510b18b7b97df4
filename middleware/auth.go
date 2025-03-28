package middleware

import (
	"errors"
	"strings"

	"test/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// AuthMiddleware проверяет наличие и валидность JWT-токена в заголовке Authorization
func AuthMiddleware(c *fiber.Ctx) error {
	// Извлекаем заголовок Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		logger.Logger.Warn("Заголовок Authorization отсутствует")
		return c.Status(401).JSON(fiber.Map{
			"Success": false,
			"Message": "Unauthorized: Missing Authorization header",
		})
	}

	// Проверяем формат заголовка (Bearer <token>)
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader || tokenString == "" {
		logger.Logger.Warn("Неверный формат заголовка Authorization")
		return c.Status(401).JSON(fiber.Map{
			"Success": false,
			"Message": "Unauthorized: Invalid Authorization format",
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

	// Парсим и проверяем токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		logger.Logger.WithError(err).Warn("Невалидный JWT-токен")
		return c.Status(401).JSON(fiber.Map{
			"Success": false,
			"Message": "Unauthorized: Invalid token",
		})
	}

	logger.Logger.Info("JWT-токен успешно проверен")
	return c.Next()
}
