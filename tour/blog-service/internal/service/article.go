package service

// Article related request structures used by handlers and services
type CountArticleRequest struct {
	TagID uint32 `json:"tag_id" form:"tag_id" binding:"gte=1"`
	State uint8  `json:"state" form:"state,default=1" binding:"oneof=0 1"`
}

type ArticleListRequest struct {
	TagID uint32 `json:"tag_id" form:"tag_id" binding:"gte=1"`
	State uint8  `json:"state" form:"state,default=1" binding:"oneof=0 1"`
}

type CreateArticleRequest struct {
	TagID         uint32 `json:"tag_id" form:"tag_id" binding:"required,gte=1"`
	Title         string `json:"title" form:"title" binding:"required,min=3,max=100"`
	Desc          string `json:"desc" form:"desc" binding:"max=255"`
	Content       string `json:"content" form:"content" binding:"required"`
	CoverImageUrl string `json:"cover_image_url" form:"cover_image_url" binding:"max=255"`
	CreatedBy     string `json:"created_by" form:"created_by" binding:"required,min=3,max=100"`
	State         uint8  `json:"state" form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateArticleRequest struct {
	ID            uint32 `json:"id" form:"id" binding:"required,gte=1"`
	TagID         uint32 `json:"tag_id" form:"tag_id" binding:"gte=1"`
	Title         string `json:"title" form:"title" binding:"min=3,max=100"`
	Desc          string `json:"desc" form:"desc" binding:"max=255"`
	Content       string `json:"content" form:"content" binding:"-"`
	CoverImageUrl string `json:"cover_image_url" form:"cover_image_url" binding:"max=255"`
	ModifiedBy    string `json:"modified_by" form:"modified_by" binding:"required,min=3,max=100"`
	State         uint8  `json:"state" form:"state" binding:"oneof=0 1"`
}

type DeleteArticleRequest struct {
	ID uint32 `json:"id" form:"id" binding:"required,gte=1"`
}
