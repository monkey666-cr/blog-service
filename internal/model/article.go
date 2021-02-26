package model

import (
	"gorm.io/gorm"
)

type Article struct {
	Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

func (a Article) TableName() string {
	return "blog_article"
}

func (a Article) Get(db *gorm.DB) (Article, error) {
	var article Article
	db = db.Where("id = ? AND state = ? AND deleted_at = ?", a.ID, a.State, nil)
	err := db.First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound{
		return article, err
	}
	return article, nil
}
