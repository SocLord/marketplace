package handlers

import (
	"encoding/json"
	"fmt"
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

func ListPublicProducts(w http.ResponseWriter, r *http.Request) {
	query := database.DB.Model(&models.Product{}).
		Preload("Shop").
		Preload("Category").
		Where("status = ?", "approved")

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		likeTerm := "%" + q + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", likeTerm, likeTerm)
	}

	categoryParam := r.URL.Query().Get("category")
	categoryIDParam := r.URL.Query().Get("category_id")
	if categoryIDParam != "" {
		if catID, err := strconv.ParseInt(categoryIDParam, 10, 64); err == nil && catID > 0 {
			query = query.Where("category_id = ?", catID)
		}
	} else if categoryParam != "" {
		var cat models.Category
		if err := database.DB.Where("slug = ?", categoryParam).First(&cat).Error; err == nil {
			query = query.Where("category_id = ?", cat.ID)
		}
	}

	shopIDParam := r.URL.Query().Get("shop_id")
	shopSlugParam := r.URL.Query().Get("shop")
	if shopIDParam != "" {
		if sID, err := strconv.ParseInt(shopIDParam, 10, 64); err == nil && sID > 0 {
			query = query.Where("shop_id = ?", sID)
		}
	} else if shopSlugParam != "" {
		var sh models.Shop
		if err := database.DB.Where("slug = ?", shopSlugParam).First(&sh).Error; err == nil {
			query = query.Where("shop_id = ?", sh.ID)
		}
	}

	minPrice := r.URL.Query().Get("min_price")
	if minPrice != "" {
		if minP, err := strconv.ParseFloat(minPrice, 64); err == nil {
			query = query.Where("price >= ?", minP)
		}
	}

	maxPrice := r.URL.Query().Get("max_price")
	if maxPrice != "" {
		if maxP, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			query = query.Where("price <= ?", maxP)
		}
	}

	sortParam := r.URL.Query().Get("sort")
	switch sortParam {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "rating":
		query = query.Order("rating DESC, reviews_count DESC")
	case "popular":
		query = query.Order("sales_count DESC, rating DESC")
	case "newest":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии дарёфти молҳо: "+err.Error())
		return
	}

	utils.Success(w, products)
}

func GetPublicProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrSlug := vars["idOrSlug"]

	query := database.DB.Model(&models.Product{}).
		Preload("Shop").
		Preload("Category").
		Preload("Reviews.User")

	var product models.Product
	var err error
	if id, parseErr := strconv.ParseInt(idOrSlug, 10, 64); parseErr == nil {
		err = query.Where("id = ? AND status != ?", id, "deleted").First(&product).Error
	} else {
		err = query.Where("slug = ? AND status != ?", idOrSlug, "deleted").First(&product).Error
	}

	if err != nil {
		utils.Error(w, http.StatusNotFound, "Мол ёфт нашуд (Товар не найден)")
		return
	}

	utils.Success(w, product)
}

func SellerListProducts(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)
	userRole, _ := r.Context().Value(middleware.ContextRole).(string)

	query := database.DB.Model(&models.Product{}).
		Preload("Category").
		Preload("Shop").
		Where("status != ?", "deleted").
		Order("created_at DESC")

	if userRole != "admin" && userRole != "moderator" {

		var shops []models.Shop
		database.DB.Where("user_id = ?", userID).Find(&shops)
		if len(shops) == 0 {
			utils.Success(w, []models.Product{})
			return
		}
		shopIDs := make([]int64, len(shops))
		for i, s := range shops {
			shopIDs[i] = s.ID
		}
		query = query.Where("shop_id IN ?", shopIDs)
	}

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии дарёфти молҳо: "+err.Error())
		return
	}

	utils.Success(w, products)
}

func SellerCreateProduct(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(w, http.StatusForbidden, "Шумо соҳиби мағоза нестед (У вас нет магазина)")
		return
	}

	var req models.ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст (Неверные данные)")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" || req.Price <= 0 || req.CategoryID <= 0 {
		utils.Error(w, http.StatusBadRequest, "Номи мол, категория ва нарх ҳатмӣ мебошанд (Название, категория и цена обязательны)")
		return
	}

	if req.ImageURL == "" {
		req.ImageURL = "https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=600"
	}

	additionalJSON := ""
	if len(req.AdditionalImages) > 0 {
		b, _ := json.Marshal(req.AdditionalImages)
		additionalJSON = string(b)
	}

	slug := fmt.Sprintf("%s-%d", utils.Slugify(req.Name), time.Now().Unix())

	product := models.Product{
		ShopID:            shop.ID,
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Slug:              slug,
		Description:       req.Description,
		Price:             req.Price,
		OldPrice:          req.OldPrice,
		StockQuantity:     req.StockQuantity,
		ImageURL:          req.ImageURL,
		AdditionalImages:  additionalJSON,
		Status:            "pending",
		ModerationComment: "Интизори санҷиши модератор (Ожидает модерации)",
		Rating:            5.0,
		ReviewsCount:      0,
		SalesCount:        0,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии сабти мол: "+err.Error())
		return
	}

	database.DB.Preload("Category").Preload("Shop").First(&product, product.ID)

	utils.Success(w, product, "Мол муваффақона илова шуд ва ба санҷиши модератор фиристода шуд")
}

func SellerUpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "ID-и нодуруст (Неверный ID)")
		return
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND status != ?", id, "deleted").First(&product).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Мол ёфт нашуд (Товар с таким ID не существует)")
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	userRole, _ := r.Context().Value(middleware.ContextRole).(string)

	isAuthorized := false
	if userRole == "admin" || userRole == "moderator" {
		isAuthorized = true
	} else {

		var shop models.Shop
		if err := database.DB.Where("id = ? AND user_id = ?", product.ShopID, userID).First(&shop).Error; err == nil {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		utils.Error(w, http.StatusForbidden, "Ин мол ба мағозаи шумо тааллуқ надорад (Этот товар принадлежит другому магазину)")
		return
	}

	var req models.ProductUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст (Неверные данные)")
		return
	}

	if strings.TrimSpace(req.Name) != "" {
		product.Name = strings.TrimSpace(req.Name)
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}
	product.Description = req.Description
	if req.Price > 0 {
		product.Price = req.Price
	}
	product.OldPrice = req.OldPrice
	product.StockQuantity = req.StockQuantity

	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if len(req.AdditionalImages) > 0 {
		b, _ := json.Marshal(req.AdditionalImages)
		product.AdditionalImages = string(b)
	}

	product.Status = "pending"
	product.ModerationComment = "Таҳрир шуд, интизори санҷиши модератор (Отредактирован, ожидает проверки)"

	if err := database.DB.Save(&product).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии навсозӣ: "+err.Error())
		return
	}

	database.DB.Preload("Category").Preload("Shop").First(&product, product.ID)

	utils.Success(w, product, "Маълумоти мол бомуваффақият навсозӣ шуд ва барои санҷиш равона гардид")
}

func SellerDeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "ID-и нодуруст (Неверный ID)")
		return
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND status != ?", id, "deleted").First(&product).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Мол ёфт нашуд (Товар не найден)")
		return
	}

	userID := middleware.GetUserIDFromContext(r)
	userRole, _ := r.Context().Value(middleware.ContextRole).(string)

	isAuthorized := false
	if userRole == "admin" || userRole == "moderator" {
		isAuthorized = true
	} else {
		var shop models.Shop
		if err := database.DB.Where("id = ? AND user_id = ?", product.ShopID, userID).First(&shop).Error; err == nil {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		utils.Error(w, http.StatusForbidden, "Ин мол ба мағозаи шумо тааллуқ надорад (Этот товар принадлежит другому магазину)")
		return
	}

	var orderItemsCount int64
	database.DB.Model(&models.OrderItem{}).Where("product_id = ?", id).Count(&orderItemsCount)

	if orderItemsCount > 0 {

		product.Status = "deleted"
		product.StockQuantity = 0
		if err := database.DB.Save(&product).Error; err != nil {
			utils.Error(w, http.StatusInternalServerError, "Хатогии несткунӣ: "+err.Error())
			return
		}
	} else {
		database.DB.Where("product_id = ?", id).Delete(&models.Review{})
		if err := database.DB.Delete(&product).Error; err != nil {
			utils.Error(w, http.StatusInternalServerError, "Хатогии несткунӣ: "+err.Error())
			return
		}
	}

	utils.Success(w, nil, "Мол бомуваффақият нест карда шуд (Товар успешно удален)")
}
