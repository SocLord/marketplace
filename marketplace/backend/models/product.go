package models

import "time"

type Product struct {
	ID                int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopID            int64     `gorm:"not null;index" json:"shop_id"`
	CategoryID        int64     `gorm:"not null;index" json:"category_id"`
	Name              string    `gorm:"type:nvarchar(300);not null" json:"name"`
	Slug              string    `gorm:"type:nvarchar(300);not null" json:"slug"`
	Description       string    `gorm:"type:nvarchar(max)" json:"description"`
	Price             float64   `gorm:"type:decimal(14,2);not null;default:0" json:"price"`
	OldPrice          *float64  `gorm:"type:decimal(14,2)" json:"old_price,omitempty"`
	StockQuantity     int       `gorm:"default:0" json:"stock_quantity"`
	ImageURL          string    `gorm:"type:nvarchar(500);not null" json:"image_url"`
	AdditionalImages  string    `gorm:"type:nvarchar(max)" json:"additional_images"`
	Status            string    `gorm:"type:nvarchar(50);not null;default:'pending';index" json:"status"`
	ModerationComment string    `gorm:"type:nvarchar(max)" json:"moderation_comment,omitempty"`
	Rating            float64   `gorm:"default:0.0" json:"rating"`
	ReviewsCount      int       `gorm:"default:0" json:"reviews_count"`
	SalesCount        int       `gorm:"default:0" json:"sales_count"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Shop     *Shop     `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Reviews  []Review  `gorm:"foreignKey:ProductID" json:"reviews,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

type ProductCreateRequest struct {
	CategoryID       int64    `json:"category_id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Price            float64  `json:"price"`
	OldPrice         *float64 `json:"old_price,omitempty"`
	StockQuantity    int      `json:"stock_quantity"`
	ImageURL         string   `json:"image_url"`
	AdditionalImages []string `json:"additional_images,omitempty"`
}

type ProductUpdateRequest struct {
	CategoryID       int64    `json:"category_id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Price            float64  `json:"price"`
	OldPrice         *float64 `json:"old_price,omitempty"`
	StockQuantity    int      `json:"stock_quantity"`
	ImageURL         string   `json:"image_url"`
	AdditionalImages []string `json:"additional_images,omitempty"`
}

type ModerationDecisionRequest struct {
	Action  string `json:"action"`
	Comment string `json:"comment,omitempty"`
}
