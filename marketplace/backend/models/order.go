package models

import "time"

type Order struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber     string    `gorm:"type:nvarchar(60);uniqueIndex;not null" json:"order_number"`
	UserID          int64     `gorm:"not null" json:"user_id"`
	TotalAmount     float64   `gorm:"type:decimal(14,2);not null;default:0" json:"total_amount"`
	Status          string    `gorm:"type:nvarchar(50);not null;default:'processing'" json:"status"` // processing, completed, cancelled
	RecipientName   string    `gorm:"type:nvarchar(150);not null" json:"recipient_name"`
	RecipientPhone  string    `gorm:"type:nvarchar(50);not null" json:"recipient_phone"`
	ShippingAddress string    `gorm:"type:nvarchar(500);not null" json:"shipping_address"`
	PaymentMethod   string    `gorm:"type:nvarchar(50);not null;default:'cash_on_delivery'" json:"payment_method"`
	Notes           string    `gorm:"type:nvarchar(max)" json:"notes,omitempty"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SubOrders []SubOrder `gorm:"foreignKey:OrderID" json:"sub_orders,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

type SubOrder struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID          int64     `gorm:"not null;index" json:"order_id"`
	ShopID           int64     `gorm:"not null;index" json:"shop_id"`
	SubtotalAmount   float64   `gorm:"type:decimal(14,2);not null;default:0" json:"subtotal_amount"`
	Status           string    `gorm:"type:nvarchar(50);not null;default:'pending'" json:"status"`
	ShippingTracking string    `gorm:"type:nvarchar(100)" json:"shipping_tracking,omitempty"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Order *Order      `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Shop  *Shop       `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	Items []OrderItem `gorm:"foreignKey:SubOrderID" json:"items,omitempty"`
}

func (SubOrder) TableName() string {
	return "sub_orders"
}

type OrderItem struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      int64     `gorm:"not null" json:"order_id"`
	SubOrderID   int64     `gorm:"not null;index" json:"sub_order_id"`
	ShopID       int64     `gorm:"not null;index" json:"shop_id"`
	ProductID    int64     `gorm:"not null" json:"product_id"`
	ProductName  string    `gorm:"type:nvarchar(300);not null" json:"product_name"`
	ProductPrice float64   `gorm:"type:decimal(14,2);not null" json:"product_price"`
	Quantity     int       `gorm:"not null;default:1" json:"quantity"`
	TotalPrice   float64   `gorm:"type:decimal(14,2);not null" json:"total_price"`
	ImageURL     string    `gorm:"type:nvarchar(500)" json:"image_url"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

type CartCheckoutItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CheckoutRequest struct {
	Items           []CartCheckoutItem `json:"items"`
	RecipientName   string             `json:"recipient_name"`
	RecipientPhone  string             `json:"recipient_phone"`
	ShippingAddress string             `json:"shipping_address"`
	PaymentMethod   string             `json:"payment_method"`
	Notes           string             `json:"notes,omitempty"`
}

type UpdateSubOrderStatusRequest struct {
	Status           string `json:"status"`
	ShippingTracking string `json:"shipping_tracking,omitempty"`
}
