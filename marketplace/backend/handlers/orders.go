package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace/database"
	"marketplace/middleware"
	"marketplace/models"
	"marketplace/utils"

	"github.com/gorilla/mux"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	var req models.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст (Invalid request body)")
		return
	}

	if len(req.Items) == 0 {
		utils.Error(w, http.StatusBadRequest, "Сабад холӣ аст (Cart is empty)")
		return
	}

	req.RecipientName = strings.TrimSpace(req.RecipientName)
	req.RecipientPhone = strings.TrimSpace(req.RecipientPhone)
	req.ShippingAddress = strings.TrimSpace(req.ShippingAddress)

	if req.RecipientName == "" || req.RecipientPhone == "" || req.ShippingAddress == "" {
		utils.Error(w, http.StatusBadRequest, "Номи қабулкунанда, рақами телефон ва суроғаи интиқол ҳатмӣ мебошанд")
		return
	}

	if req.PaymentMethod == "" {
		req.PaymentMethod = "cash_on_delivery"
	}

	type ItemWithProduct struct {
		Product  models.Product
		Quantity int
		Subtotal float64
	}
	shopItemsMap := make(map[int64][]ItemWithProduct)
	var grandTotal float64

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			item.Quantity = 1
		}
		var prod models.Product
		if err := database.DB.Where("id = ? AND status = ?", item.ProductID, "approved").First(&prod).Error; err != nil {
			utils.Error(w, http.StatusBadRequest, fmt.Sprintf("Мол бо рамзи #%d дастрас нест ё нашр нашудааст", item.ProductID))
			return
		}

		if prod.StockQuantity < item.Quantity {
			utils.Error(w, http.StatusBadRequest, fmt.Sprintf("Барои моли '%s' танҳо %d дона дар анбор мавҷуд аст", prod.Name, prod.StockQuantity))
			return
		}

		subtotal := prod.Price * float64(item.Quantity)
		grandTotal += subtotal

		shopItemsMap[prod.ShopID] = append(shopItemsMap[prod.ShopID], ItemWithProduct{
			Product:  prod,
			Quantity: item.Quantity,
			Subtotal: subtotal,
		})
	}

	orderNumber := fmt.Sprintf("ORD-%s-%04d", time.Now().Format("20060102"), rand.Intn(9000)+1000)
	masterOrder := models.Order{
		OrderNumber:     orderNumber,
		UserID:          userID,
		TotalAmount:     grandTotal,
		Status:          "processing",
		RecipientName:   req.RecipientName,
		RecipientPhone:  req.RecipientPhone,
		ShippingAddress: req.ShippingAddress,
		PaymentMethod:   req.PaymentMethod,
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	tx := database.DB.Begin()

	if err := tx.Create(&masterOrder).Error; err != nil {
		tx.Rollback()
		utils.Error(w, http.StatusInternalServerError, "Хатогии эҷоди фармоиш: "+err.Error())
		return
	}

	var createdSubOrders []models.SubOrder
	for shopID, items := range shopItemsMap {
		var shopSubtotal float64
		for _, it := range items {
			shopSubtotal += it.Subtotal
		}

		trackingCode := fmt.Sprintf("TRK-%d-%05d", shopID, rand.Intn(90000)+10000)

		subOrder := models.SubOrder{
			OrderID:          masterOrder.ID,
			ShopID:           shopID,
			SubtotalAmount:   shopSubtotal,
			Status:           "pending",
			ShippingTracking: trackingCode,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := tx.Create(&subOrder).Error; err != nil {
			tx.Rollback()
			utils.Error(w, http.StatusInternalServerError, "Хатогии тақсими фармоиш: "+err.Error())
			return
		}

		for _, it := range items {
			orderItem := models.OrderItem{
				OrderID:      masterOrder.ID,
				SubOrderID:   subOrder.ID,
				ShopID:       shopID,
				ProductID:    it.Product.ID,
				ProductName:  it.Product.Name,
				ProductPrice: it.Product.Price,
				Quantity:     it.Quantity,
				TotalPrice:   it.Subtotal,
				ImageURL:     it.Product.ImageURL,
				CreatedAt:    time.Now(),
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				tx.Rollback()
				utils.Error(w, http.StatusInternalServerError, "Хатогии сабти моли фармоиш: "+err.Error())
				return
			}

			tx.Model(&models.Product{}).Where("id = ?", it.Product.ID).Updates(map[string]interface{}{
				"stock_quantity": it.Product.StockQuantity - it.Quantity,
				"sales_count":    it.Product.SalesCount + it.Quantity,
			})
		}

		createdSubOrders = append(createdSubOrders, subOrder)
	}

	tx.Commit()

	var fullOrder models.Order
	database.DB.Preload("SubOrders.Shop").
		Preload("SubOrders.Items").
		First(&fullOrder, masterOrder.ID)

	utils.Success(w, fullOrder, "Фармоиш бомуваффақият сабт шуд ва ба қисмҳои фурӯшандагон тақсим гардид!")
}

func GetBuyerOrders(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	var orders []models.Order
	if err := database.DB.Where("user_id = ?", userID).
		Preload("SubOrders.Shop").
		Preload("SubOrders.Items").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии дарёфти фармоишҳо: "+err.Error())
		return
	}

	utils.Success(w, orders)
}

func GetSellerSubOrders(w http.ResponseWriter, r *http.Request) {
	shopID, ok := middleware.GetShopIDFromContext(r)
	if !ok || shopID <= 0 {
		userID := middleware.GetUserIDFromContext(r)
		var shop models.Shop
		if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
			utils.Error(w, http.StatusForbidden, "Мағозаи фурӯшанда ёфт нашуд")
			return
		}
		shopID = shop.ID
	}

	statusFilter := r.URL.Query().Get("status")
	query := database.DB.Where("shop_id = ?", shopID).
		Preload("Order").
		Preload("Items").
		Order("created_at DESC")

	if statusFilter != "" && statusFilter != "all" {
		query = query.Where("status = ?", statusFilter)
	}

	var subOrders []models.SubOrder
	if err := query.Find(&subOrders).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии дарёфти фармоишҳои фурӯшанда: "+err.Error())
		return
	}

	utils.Success(w, subOrders)
}

func UpdateSubOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subOrderID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "ID-и нодуруст")
		return
	}

	shopID, ok := middleware.GetShopIDFromContext(r)
	if !ok || shopID <= 0 {
		userID := middleware.GetUserIDFromContext(r)
		var shop models.Shop
		if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
			utils.Error(w, http.StatusForbidden, "Мағоза ёфт нашуд")
			return
		}
		shopID = shop.ID
	}

	var subOrder models.SubOrder
	if err := database.DB.Where("id = ? AND shop_id = ?", subOrderID, shopID).First(&subOrder).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Ин қисми фармоиш ёфт нашуд ё ба мағозаи шумо дахл надорад")
		return
	}

	var req models.UpdateSubOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст")
		return
	}

	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
	}

	req.Status = strings.ToLower(strings.TrimSpace(req.Status))
	if !validStatuses[req.Status] {
		utils.Error(w, http.StatusBadRequest, "Ҳолати нодуруст. Ҳолатҳои дуруст: pending, processing, shipped, delivered, cancelled")
		return
	}

	subOrder.Status = req.Status
	if req.ShippingTracking != "" {
		subOrder.ShippingTracking = req.ShippingTracking
	}

	if err := database.DB.Save(&subOrder).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии навсозии ҳолат: "+err.Error())
		return
	}

	checkMasterOrderCompletion(subOrder.OrderID)

	utils.Success(w, subOrder, fmt.Sprintf("Ҳолати қисми фармоиш ба '%s' иваз карда шуд", req.Status))
}

func checkMasterOrderCompletion(orderID int64) {
	var subOrders []models.SubOrder
	database.DB.Where("order_id = ?", orderID).Find(&subOrders)

	allDelivered := true
	for _, so := range subOrders {
		if so.Status != "delivered" {
			allDelivered = false
			break
		}
	}

	if allDelivered && len(subOrders) > 0 {
		database.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "completed")
	}
}
