package models

// News представляет таблицу новостей
type News struct {
	Id         uint   `json:"Id"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	Categories []uint `json:"Categories"` // Список ID категорий
}

// NewsCategory представляет связь между новостями и категориями
type NewsCategory struct {
	NewsId     uint `gorm:"primaryKey"`
	CategoryId uint `gorm:"primaryKey"`
}
