package database

import (
	"fmt"
	"log"
	"time"

	"test/models"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	// Настройка Viper для загрузки .env файла
	viper.SetConfigFile(".env") // Указываем файл .env
	viper.AutomaticEnv()        // Разрешаем использование переменных окружения

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения системы")
	}

	// Получение значений переменных окружения через Viper
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")

	// Формирование строки подключения
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		dbHost, dbPort, dbUser, dbName, dbPassword)

	// Подключение к базе данных
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	// Получение пула соединений
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("ошибка получения пула соединений: %v", err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Успешно подключились к базе данных")
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(&models.News{}, &models.NewsCategory{}) // Используем полные пути к моделям
	if err != nil {
		return fmt.Errorf("ошибка при выполнении миграций: %v", err)
	}
	log.Println("Миграции успешно выполнены")
	return nil
}
