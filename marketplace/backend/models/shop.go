package models

import "time"

type Shop struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"not null" json:"user_id"`
	Name        string    `gorm:"type:nvarchar(200);not null" json:"name"`
	Slug        string    `gorm:"type:nvarchar(200);uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:nvarchar(max)" json:"description"`
	LogoURL     string    `gorm:"type:nvarchar(500)" json:"logo_url"`
	BannerURL   string    `gorm:"type:nvarchar(500)" json:"banner_url"`
	Phone       string    `gorm:"type:nvarchar(50)" json:"phone"`
	Address     string    `gorm:"type:nvarchar(255)" json:"address"`
	Rating      float64   `gorm:"default:5.0" json:"rating"`
	IsVerified  bool      `gorm:"default:true" json:"is_verified"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Owner    *User      `gorm:"foreignKey:UserID" json:"owner,omitempty"`
	Products []Product  `gorm:"foreignKey:ShopID" json:"products,omitempty"`
}

func (Shop) TableName() string {
	return "shops"
}

type ShopUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	BannerURL   string `json:"banner_url"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
}
