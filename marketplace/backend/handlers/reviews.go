package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"marketplace/database"
	"marketplace/middleware"
	"marketplace/models"
	"marketplace/utils"

	"github.com/gorilla/mux"
)

func CreateReview(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r)

	var req models.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст")
		return
	}

	req.Comment = strings.TrimSpace(req.Comment)
	if req.ProductID <= 0 || req.Rating < 1 || req.Rating > 5 || req.Comment == "" {
		utils.Error(w, http.StatusBadRequest, "Интихоби мол, баҳо (1 то 5) ва матни тақриз ҳатмӣ аст")
		return
	}

	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Мол ёфт нашуд")
		return
	}

	review := models.Review{
		ProductID: req.ProductID,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := database.DB.Create(&review).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии сабти фикру мулоҳиза: "+err.Error())
		return
	}

	type RatingStats struct {
		AvgRating float64
		Count     int
	}
	var stats RatingStats
	database.DB.Model(&models.Review{}).
		Select("AVG(CAST(rating AS FLOAT)) as avg_rating, COUNT(id) as count").
		Where("product_id = ?", req.ProductID).
		Scan(&stats)

	roundedRating := math.Round(stats.AvgRating*10) / 10
	database.DB.Model(&models.Product{}).
		Where("id = ?", req.ProductID).
		Updates(map[string]interface{}{
			"rating":        roundedRating,
			"reviews_count": stats.Count,
		})

	database.DB.Preload("User").First(&review, review.ID)

	utils.Success(w, review, "Тақризи шумо қабул гардид ва баҳо ба ҳисоб гирифта шуд")
}

func GetProductReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "ID-и нодуруст")
		return
	}

	var reviews []models.Review
	if err := database.DB.Where("product_id = ?", productID).
		Preload("User").
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии дарёфти тақризҳо")
		return
	}

	utils.Success(w, reviews)
}
