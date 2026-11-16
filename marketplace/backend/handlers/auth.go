package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"marketplace/database"
	"marketplace/middleware"
	"marketplace/models"
	"marketplace/utils"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст (Invalid JSON)")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" || req.Name == "" {
		utils.Error(w, http.StatusBadRequest, "Ном, почтаи электронӣ ва парол ҳатмӣ мебошанд")
		return
	}

	if req.Role == "" {
		req.Role = "buyer"
	}
	if req.Role != "buyer" && req.Role != "seller" {
		req.Role = "buyer"
	}

	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		utils.Error(w, http.StatusConflict, "Ин суроғаи почтаи электронӣ аллакай сабти ном шудааст")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии эҷоди парол")
		return
	}

	newUser := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         req.Role,
		Phone:        req.Phone,
		AvatarURL:    fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=0d9488&color=fff", utils.Slugify(req.Name)),
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии сабти корбар: "+err.Error())
		return
	}

	var createdShop *models.Shop
	if req.Role == "seller" {
		shopName := req.ShopName
		if shopName == "" {
			shopName = fmt.Sprintf("Мағозаи %s", req.Name)
		}
		shop := models.Shop{
			UserID:      newUser.ID,
			Name:        shopName,
			Slug:        fmt.Sprintf("%s-%d", utils.Slugify(shopName), newUser.ID),
			Description: "Мағозаи нави расмӣ дар маркетплейс",
			LogoURL:     "https://images.unsplash.com/photo-1472851294608-062f824d29cc?w=200",
			BannerURL:   "https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=1200",
			Phone:       req.Phone,
			Rating:      5.0,
			IsVerified:  true,
		}
		if err := database.DB.Create(&shop).Error; err == nil {
			createdShop = &shop
		}
	}

	var shopID *int64
	var shopName, shopSlug string
	if createdShop != nil {
		shopID = &createdShop.ID
		shopName = createdShop.Name
		shopSlug = createdShop.Slug
	}

	token, err := utils.GenerateToken(newUser.ID, newUser.Email, newUser.Role, shopID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии эҷоди токен")
		return
	}

	utils.Success(w, models.AuthResponse{
		Token: token,
		User: models.UserSummary{
			ID:        newUser.ID,
			Name:      newUser.Name,
			Email:     newUser.Email,
			Role:      newUser.Role,
			Phone:     newUser.Phone,
			AvatarURL: newUser.AvatarURL,
			ShopID:    shopID,
			ShopName:  shopName,
			ShopSlug:  shopSlug,
		},
	}, "Сабти ном бомуваффақият анҷом ёфт")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.Error(w, http.StatusUnauthorized, "Почтаи электронӣ ё парол нодуруст аст")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		utils.Error(w, http.StatusUnauthorized, "Почтаи электронӣ ё парол нодуруст аст")
		return
	}

	var shop models.Shop
	var shopID *int64
	var shopName, shopSlug string
	if err := database.DB.Where("user_id = ?", user.ID).First(&shop).Error; err == nil {
		shopID = &shop.ID
		shopName = shop.Name
		shopSlug = shop.Slug
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, shopID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии тавлиди токен")
		return
	}

	utils.Success(w, models.AuthResponse{
		Token: token,
		User: models.UserSummary{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			Phone:     user.Phone,
			AvatarURL: user.AvatarURL,
			ShopID:    shopID,
			ShopName:  shopName,
			ShopSlug:  shopSlug,
		},
	}, "Воридшавӣ бомуваффақият иҷро шуд")
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Корбар ёфт нашуд")
		return
	}

	var shop models.Shop
	var shopID *int64
	var shopName, shopSlug string
	if err := database.DB.Where("user_id = ?", user.ID).First(&shop).Error; err == nil {
		shopID = &shop.ID
		shopName = shop.Name
		shopSlug = shop.Slug
	}

	utils.Success(w, models.UserSummary{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Phone:     user.Phone,
		AvatarURL: user.AvatarURL,
		ShopID:    shopID,
		ShopName:  shopName,
		ShopSlug:  shopSlug,
	})
}
