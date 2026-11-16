package models

import "time"

type Review struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID int64     `gorm:"not null;index" json:"product_id"`
	UserID    int64     `gorm:"not null" json:"user_id"`
	OrderID   *int64    `json:"order_id,omitempty"`
	Rating    int       `gorm:"not null" json:"rating"` // 1 - 5
	Comment   string    `gorm:"type:nvarchar(max);not null" json:"comment"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (Review) TableName() string {
	return "reviews"
}

type CreateReviewRequest struct {
	ProductID int64  `json:"product_id"`
	Rating    int    `json:"rating"` // 1-5
	Comment   string `json:"comment"`
}
