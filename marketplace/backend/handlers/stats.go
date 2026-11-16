package handlers

import (
	"fmt"
	"net/http"
	"time"

	"marketplace/database"
	"marketplace/middleware"
	"marketplace/models"
	"marketplace/utils"
)

type MonthlyStat struct {
	Year        int     `json:"year"`
	Month       int     `json:"month"`
	MonthName   string  `json:"month_name"`
	Revenue     float64 `json:"revenue"`
	OrdersCount int     `json:"orders_count"`
	ItemsSold   int     `json:"items_sold"`
}

type TopProductStat struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	ImageURL    string  `json:"image_url"`
	UnitsSold   int     `json:"units_sold"`
	TotalEarned float64 `json:"total_earned"`
}

type SellerDashboardStats struct {
	ShopID              int64            `json:"shop_id"`
	ShopName            string           `json:"shop_name"`
	TotalRevenue        float64          `json:"total_revenue"`
	CurrentMonthRevenue float64          `json:"current_month_revenue"`
	TotalOrdersCount    int64            `json:"total_orders_count"`
	TotalProductsCount  int64            `json:"total_products_count"`
	PendingProducts     int64            `json:"pending_products"`
	MonthlyStats        []MonthlyStat    `json:"monthly_stats"`
	TopProducts         []TopProductStat `json:"top_products"`
}

var monthNamesTj = []string{
	"", "Январ", "Феврал", "Март", "Апрел", "Май", "Июн",
	"Июл", "Август", "Сентябр", "Октябр", "Ноябр", "Декабр",
}

func GetSellerMonthlyStats(w http.ResponseWriter, r *http.Request) {
	shopID, ok := middleware.GetShopIDFromContext(r)
	if !ok || shopID <= 0 {
		userID := middleware.GetUserIDFromContext(r)
		var shop models.Shop
		if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
			utils.Error(w, http.StatusForbidden, "Мағозаи шумо ёфт нашуд")
			return
		}
		shopID = shop.ID
	}

	var shop models.Shop
	database.DB.First(&shop, shopID)

	var totalRev float64
	database.DB.Model(&models.SubOrder{}).
		Where("shop_id = ? AND status != ?", shopID, "cancelled").
		Select("COALESCE(SUM(subtotal_amount), 0)").
		Scan(&totalRev)

	var totalOrders int64
	database.DB.Model(&models.SubOrder{}).
		Where("shop_id = ?", shopID).
		Count(&totalOrders)

	var totalProducts, pendingProducts int64
	database.DB.Model(&models.Product{}).Where("shop_id = ?", shopID).Count(&totalProducts)
	database.DB.Model(&models.Product{}).Where("shop_id = ? AND status = ?", shopID, "pending").Count(&pendingProducts)

	now := time.Now()
	var monthlyStats []MonthlyStat

	for i := 5; i >= 0; i-- {
		targetDate := now.AddDate(0, -i, 0)
		year := targetDate.Year()
		month := int(targetDate.Month())

		startOfMonth := time.Date(year, targetDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, 0)

		var monthRev float64
		database.DB.Model(&models.SubOrder{}).
			Where("shop_id = ? AND status != ? AND created_at >= ? AND created_at < ?", shopID, "cancelled", startOfMonth, endOfMonth).
			Select("COALESCE(SUM(subtotal_amount), 0)").
			Scan(&monthRev)

		var monthOrders int64
		database.DB.Model(&models.SubOrder{}).
			Where("shop_id = ? AND created_at >= ? AND created_at < ?", shopID, startOfMonth, endOfMonth).
			Count(&monthOrders)

		var itemsSold int64
		database.DB.Model(&models.OrderItem{}).
			Where("shop_id = ? AND created_at >= ? AND created_at < ?", shopID, startOfMonth, endOfMonth).
			Select("COALESCE(SUM(quantity), 0)").
			Scan(&itemsSold)

		monthName := fmt.Sprintf("%s %d", monthNamesTj[month], year)

		monthlyStats = append(monthlyStats, MonthlyStat{
			Year:        year,
			Month:       month,
			MonthName:   monthName,
			Revenue:     monthRev,
			OrdersCount: int(monthOrders),
			ItemsSold:   int(itemsSold),
		})
	}

	curMonthRev := monthlyStats[len(monthlyStats)-1].Revenue

	type TopProdRow struct {
		ProductID   int64
		ProductName string
		ImageURL    string
		UnitsSold   int
		TotalEarned float64
	}
	var topRows []TopProdRow
	database.DB.Model(&models.OrderItem{}).
		Select("product_id, product_name, MAX(image_url) as image_url, SUM(quantity) as units_sold, SUM(total_price) as total_earned").
		Where("shop_id = ?", shopID).
		Group("product_id, product_name").
		Order("units_sold DESC").
		Limit(5).
		Scan(&topRows)

	var topProducts []TopProductStat
	for _, r := range topRows {
		topProducts = append(topProducts, TopProductStat{
			ProductID:   r.ProductID,
			ProductName: r.ProductName,
			ImageURL:    r.ImageURL,
			UnitsSold:   r.UnitsSold,
			TotalEarned: r.TotalEarned,
		})
	}

	resp := SellerDashboardStats{
		ShopID:              shopID,
		ShopName:            shop.Name,
		TotalRevenue:        totalRev,
		CurrentMonthRevenue: curMonthRev,
		TotalOrdersCount:    totalOrders,
		TotalProductsCount:  totalProducts,
		PendingProducts:     pendingProducts,
		MonthlyStats:        monthlyStats,
		TopProducts:         topProducts,
	}

	utils.Success(w, resp)
}

type ModeratorStatsResponse struct {
	TotalUsers      int64   `json:"total_users"`
	TotalSellers    int64   `json:"total_sellers"`
	TotalBuyers     int64   `json:"total_buyers"`
	TotalShops      int64   `json:"total_shops"`
	TotalProducts   int64   `json:"total_products"`
	PendingProducts int64   `json:"pending_products"`
	TotalOrders     int64   `json:"total_orders"`
	TotalTurnover   float64 `json:"total_turnover"`
}

func GetModeratorStats(w http.ResponseWriter, r *http.Request) {
	var totalUsers, totalSellers, totalBuyers, totalShops, totalProducts, pendingProducts, totalOrders int64
	var totalTurnover float64

	database.DB.Model(&models.User{}).Count(&totalUsers)
	database.DB.Model(&models.User{}).Where("role = ?", "seller").Count(&totalSellers)
	database.DB.Model(&models.User{}).Where("role = ?", "buyer").Count(&totalBuyers)
	database.DB.Model(&models.Shop{}).Count(&totalShops)
	database.DB.Model(&models.Product{}).Count(&totalProducts)
	database.DB.Model(&models.Product{}).Where("status = ?", "pending").Count(&pendingProducts)
	database.DB.Model(&models.Order{}).Count(&totalOrders)
	database.DB.Model(&models.Order{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalTurnover)

	utils.Success(w, ModeratorStatsResponse{
		TotalUsers:      totalUsers,
		TotalSellers:    totalSellers,
		TotalBuyers:     totalBuyers,
		TotalShops:      totalShops,
		TotalProducts:   totalProducts,
		PendingProducts: pendingProducts,
		TotalOrders:     totalOrders,
		TotalTurnover:   totalTurnover,
	})
}
