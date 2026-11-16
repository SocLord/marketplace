package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"marketplace/database"
	"marketplace/models"
	"marketplace/utils"

	"github.com/gorilla/mux"
)

func ListCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	if err := database.DB.Order("sort_order ASC, name ASC").Find(&categories).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогӣ: "+err.Error())
		return
	}

	for i := range categories {
		var count int64
		database.DB.Model(&models.Product{}).
			Where("category_id = ? AND status = ?", categories[i].ID, "approved").
			Count(&count)
		categories[i].ProductsCount = count
	}

	utils.Success(w, categories)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var cat models.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст")
		return
	}

	cat.Name = strings.TrimSpace(cat.Name)
	if cat.Name == "" {
		utils.Error(w, http.StatusBadRequest, "Номи категория ҳатмӣ аст")
		return
	}

	if cat.Slug == "" {
		cat.Slug = utils.Slugify(cat.Name)
	}

	if err := database.DB.Create(&cat).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии сабт: "+err.Error())
		return
	}

	utils.Success(w, cat, "Категория бомуваффақият илова шуд")
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "ID-и нодуруст")
		return
	}

	var cat models.Category
	if err := database.DB.First(&cat, id).Error; err != nil {
		utils.Error(w, http.StatusNotFound, "Категория ёфт нашуд")
		return
	}

	var req models.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Маълумоти нодуруст")
		return
	}

	if req.Name != "" {
		cat.Name = req.Name
	}
	if req.Slug != "" {
		cat.Slug = req.Slug
	}
	cat.Icon = req.Icon
	cat.ImageURL = req.ImageURL
	cat.Description = req.Description
	cat.SortOrder = req.SortOrder

	if err := database.DB.Save(&cat).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогӣ: "+err.Error())
		return
	}

	utils.Success(w, cat, "Категория нав карда шуд")
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "ID-и нодуруст")
		return
	}

	var count int64
	database.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		utils.Error(w, http.StatusBadRequest, fmt.Sprintf("Дар ин категория %d мол мавҷуд аст. Аввал онҳоро интиқол ё нест кунед.", count))
		return
	}

	if err := database.DB.Delete(&models.Category{}, id).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии несткунӣ")
		return
	}

	utils.Success(w, nil, "Категория нест карда шуд")
}
