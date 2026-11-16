package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"marketplace/utils"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.Error(w, http.StatusBadRequest, "Ҳаҷми файл хеле калон аст (макс 10MB)")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Файл ёфт нашуд")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		utils.Error(w, http.StatusBadRequest, "Танҳо файлҳои расм (jpg, jpeg, png, webp) иҷозат дода шудаанд")
		return
	}

	uploadsDir := "./uploads"
	_ = os.MkdirAll(uploadsDir, os.ModePerm)

	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), utils.Slugify(strings.TrimSuffix(handler.Filename, ext)), ext)
	dstPath := filepath.Join(uploadsDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии эҷоди файл")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Хатогии сабти файл")
		return
	}

	fileURL := fmt.Sprintf("http://localhost:8089/uploads/%s", filename)
	utils.Success(w, map[string]string{
		"url":      fileURL,
		"filename": filename,
	}, "Расм бомуваффақият боргузорӣ шуд")
}
