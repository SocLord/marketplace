package models

import "time"

type Category struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:nvarchar(150);not null" json:"name"`
	Slug        string    `gorm:"type:nvarchar(150);uniqueIndex;not null" json:"slug"`
	Icon        string    `gorm:"type:nvarchar(100)" json:"icon"`
	ImageURL    string    `gorm:"type:nvarchar(500)" json:"image_url"`
	Description string    `gorm:"type:nvarchar(max)" json:"description"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	ProductsCount int64 `gorm:"-" json:"products_count,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
