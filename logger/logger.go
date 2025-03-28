package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger экземпляр Logrus
var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	// Устанавливаем формат вывода (JSON)
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// Устанавливаем уровень логгирования
	Logger.SetLevel(logrus.DebugLevel)

	// Устанавливаем вывод в стандартный поток ошибок
	Logger.SetOutput(os.Stderr)
}
