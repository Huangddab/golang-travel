package model

type Model struct {
	ID         uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedBy  int64  `json:"created_by"`
	ModifiedBy int64  `json:"modified_by"`
	CreatedOn  int64  `json:"created_on"`
	ModifiedOn int64  `json:"modified_on"`
	DeletedOn  int64  `json:"deleted_on"`
	IsDel      int8   `json:"is_del"`
}
