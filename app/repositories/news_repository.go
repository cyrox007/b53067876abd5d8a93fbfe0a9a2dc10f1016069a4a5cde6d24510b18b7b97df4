package repositories

import (
	"errors"
	"fmt"

	"test/app/models"
)

type NewsRepository struct {
	db *DB
}

func NewNewsRepository(db *DB) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) UpdateNews(id int64, title, content *string, categories []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Находим новость по ID
	news := &models.News{Id: id}
	if err := tx.FindByPrimaryKeyTo(news, id); err != nil {
		return errors.New("news not found")
	}

	// Обновляем поля новости, если они переданы
	if title != nil {
		news.Title = *title
	}
	if content != nil {
		news.Content = *content
	}

	// Сохраняем изменения в новости
	if err := tx.Save(news); err != nil {
		return fmt.Errorf("failed to save news: %w", err)
	}

	// Удаляем старые категории для новости
	const newsCategoriesTable = "NewsCategories"
	if _, err := tx.DeleteFrom(newsCategoriesTable, "NewsId = $1", id); err != nil {
		return fmt.Errorf("failed to delete old categories: %w", err)
	}

	// Добавляем новые категории
	for _, catID := range categories {
		category := &models.NewsCategory{NewsId: id, CategoryId: catID}
		if err := tx.Insert(category); err != nil {
			return fmt.Errorf("failed to insert category: %w", err)
		}
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *NewsRepository) ListNews(limit, offset int) ([]*models.News, error) {
	var news []*models.News
	query := "ORDER BY Id DESC LIMIT $1 OFFSET $2"
	err := r.db.SelectAll(&news, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list news: %w", err)
	}
	return news, nil
}
