package models

import "time"

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:nvarchar(150);not null" json:"name"`
	Email        string    `gorm:"type:nvarchar(150);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:nvarchar(255);not null" json:"-"`
	Role         string    `gorm:"type:nvarchar(50);not null;default:'buyer'" json:"role"`
	Phone        string    `gorm:"type:nvarchar(50)" json:"phone"`
	AvatarURL    string    `gorm:"type:nvarchar(500)" json:"avatar_url"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Shop *Shop `gorm:"foreignKey:UserID" json:"shop,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Phone    string `json:"phone"`
	ShopName string `json:"shop_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  UserSummary `json:"user"`
}

type UserSummary struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
	ShopID    *int64 `json:"shop_id,omitempty"`
	ShopName  string `json:"shop_name,omitempty"`
	ShopSlug  string `json:"shop_slug,omitempty"`
}
