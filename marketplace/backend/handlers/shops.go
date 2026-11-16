package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"marketplace/database"
	"marketplace/middleware"
	"marketplace/models"
	"marketplace/utils"

	"github.com/gorilla/mux"
)

func ListShops(w http.ResponseWriter, r *http.Request) {
	var shops []models.Shop
	if err := database.DB.Where("is_verified = ?", true).Order("rating DESC, name ASC").Find(&shops).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогӣ: "+err.Error())
		return
	}

	utils.Success(w, shops)
}

func GetShop(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrSlug := vars["idOrSlug"]

	var shop models.Shop
	var err error
	if id, parseErr := strconv.ParseInt(idOrSlug, 10, 64); parseErr == nil {
		err = database.DB.Where("id = ?", id).First(&shop).Error
	} else {
		err = database.DB.Where("slug = ?", idOrSlug).First(&shop).Error
	}

	if err != nil {
		utils.Error(w, http.StatusNotFound, "Мағоза ёфт нашуд")
		return
	}

	var products []models.Product
	database.DB.Where("shop_id = ? AND status = ?", shop.ID, "approved").
		Preload("Category").
		Order("rating DESC, sales_count DESC").
		Find(&products)
	shop.Products = products

	utils.Success(w, shop)
}

func GetMyShop(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Мағозаи шумо ёфт нашуд")
		return
	}

	utils.Success(w, shop)
}

func UpdateMyShop(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Мағозаи шумо ёфт нашуд")
		return
	}

	var req models.ShopUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст")
		return
	}

	if req.Name != "" {
		shop.Name = req.Name
	}
	shop.Description = req.Description
	if req.LogoURL != "" {
		shop.LogoURL = req.LogoURL
	}
	if req.BannerURL != "" {
		shop.BannerURL = req.BannerURL
	}
	shop.Phone = req.Phone
	shop.Address = req.Address

	if err := database.DB.Save(&shop).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии сабти мағоза: "+err.Error())
		return
	}

	utils.Success(w, shop, "Танзимоти мағоза бомуваффақият сабт шуд")
}
