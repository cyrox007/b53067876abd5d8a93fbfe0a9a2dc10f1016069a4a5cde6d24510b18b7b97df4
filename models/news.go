package models

// News представляет таблицу новостей (для GORM)
type News struct {
	Id      uint   `gorm:"primaryKey;autoIncrement" json:"Id"`
	Title   string `gorm:"size:255;not null" json:"Title"`
	Content string `gorm:"type:text;not null" json:"Content"`
}

// NewsCategory представляет связь между новостями и категориями
type NewsCategory struct {
	NewsId     uint `gorm:"primaryKey"`
	CategoryId uint `gorm:"primaryKey"`
}

// NewsResponse представляет ответ для клиента (JSON)
type NewsResponse struct {
	Id         uint   `json:"Id"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	Categories []uint `json:"Categories"` // Список ID категорий
}
