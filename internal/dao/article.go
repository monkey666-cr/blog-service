package dao

import "blog-service/internal/model"

func (d *Dao) GetArticle(id uint32, state uint8) (model.Article, error) {
	article := model.Article{Model: model.Model{ID: id}, State: state}
	return article.Get(d.engine)
}
