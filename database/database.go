package database

import (
	"fmt"
	"log"
	"time"

	"test/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		"localhost", "5432", "username", "newsdb", "yourpassword")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

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
