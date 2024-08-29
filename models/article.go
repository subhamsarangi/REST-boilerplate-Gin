package models

import "gorm.io/gorm"

type GoArticle struct {
	gorm.Model
	Title   string
	Content string
	UserID  uint
	User    GoUser `gorm:"foreignKey:UserID"`
}
