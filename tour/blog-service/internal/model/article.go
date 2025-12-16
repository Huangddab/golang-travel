package model

import "blog-service/pkg/app"

type ArticleSwagger struct {
	List  []*Tag
	Pager *app.Pager
}

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	State         uint8  `json:"state"`
	CoverImageUrl string `json:"cover_image_url"`
}

func (a *Article) TableName() string {
	return "blog_article"
}
