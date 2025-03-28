package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *sql.DB

func Connect() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения системы")
	}

	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		dbHost, dbPort, dbUser, dbName, dbPassword)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}

	log.Println("Успешно подключились к базе данных")
	return nil
}
