package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"marketplace/database"
	"marketplace/models"
	"marketplace/utils"

	"github.com/gorilla/mux"
)

func ListPendingProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	if err := database.DB.Where("status = ?", "pending").
		Preload("Shop").
		Preload("Category").
		Order("created_at ASC").
		Find(&products).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии дарёфти молҳои навбатӣ: "+err.Error())
		return
	}

	utils.Success(w, products)
}

func ListAllProductsModerator(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	query := database.DB.Model(&models.Product{}).
		Preload("Shop").
		Preload("Category").
		Order("updated_at DESC")

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогӣ: "+err.Error())
		return
	}

	utils.Success(w, products)
}

func ReviewProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "ID-и нодуруст")
		return
	}

	var req models.ModerationDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст")
		return
	}

	req.Action = strings.ToLower(strings.TrimSpace(req.Action))
	if req.Action != "approve" && req.Action != "reject" {
		utils.Error(w, http.StatusBadRequest, "Амал бояд 'approve' (тасдиқ) ё 'reject' (рад) бошад")
		return
	}

	var product models.Product
	if err := database.DB.Preload("Shop").Preload("Category").First(&product, id).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Мол ёфт нашуд")
		return
	}

	if req.Action == "approve" {
		product.Status = "approved"
		if req.Comment == "" {
			product.ModerationComment = "Мол аз ҷониби модератор санҷида ва тасдиқ шуд."
		} else {
			product.ModerationComment = req.Comment
		}
	} else {
		product.Status = "rejected"
		if req.Comment == "" {
			product.ModerationComment = "Мол ба талаботи маркетплейс мувофиқат намекунад."
		} else {
			product.ModerationComment = req.Comment
		}
	}

	if err := database.DB.Save(&product).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии сабти қарор: "+err.Error())
		return
	}

	message := "Мол бомуваффақият нашр карда шуд"
	if product.Status == "rejected" {
		message = "Мол рад карда шуд"
	}

	utils.Success(w, product, message)
}
