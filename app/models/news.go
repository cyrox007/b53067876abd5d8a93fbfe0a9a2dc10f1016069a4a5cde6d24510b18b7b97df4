package models

import (
	_ "gopkg.in/reform.v1"
)

// News represents the News table.
type News struct {
	Id      int64  `reform:"Id,pk"`
	Title   string `reform:"Title"`
	Content string `reform:"Content"`
}

// NewsCategory represents the NewsCategories table.
type NewsCategory struct {
	NewsId     int64 `reform:"NewsId,pk"`
	CategoryId int64 `reform:"CategoryId,pk"`
}
