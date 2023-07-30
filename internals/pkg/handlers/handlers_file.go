package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/helpers"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/models"
	"github.com/ishanshre/GoFileServerAPI/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	root_path = "./media"
)

func (h *handlers) UploadSingleFile(w http.ResponseWriter, r *http.Request) {
	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)

	if err := r.ParseMultipartForm(1 * 1024 * 1024); err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	file, header, err := r.FormFile("files")
	if err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	defer file.Close()
	fileName := header.Filename
	userPath := filepath.Join(root_path, tokenDetail.Username)
	if err := os.MkdirAll(userPath, os.ModePerm); err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	relativePath := filepath.Join(userPath, fileName)
	dst, err := os.Create(relativePath)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	defer dst.Close()
	if err := helpers.SetFilePermission(relativePath); err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	_, err = io.Copy(dst, file)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	id, _ := primitive.ObjectIDFromHex(tokenDetail.UserID)
	uploader := &models.UserAccess{}
	uploader.Username = tokenDetail.Username
	uploader.ID = id
	uploader.AccessLevel = tokenDetail.AccessLevel
	fileData := &models.File{}
	fileData.FilePath = userPath
	fileData.Name = fileName
	fileData.Extension = filepath.Ext(fileName)
	fileData.Size = header.Size
	fileData.Uploader = uploader
	fileData.UploadedAt = time.Now()
	fileData.ModifiedAt = time.Now()
	res, err := h.mg.InsertFileData(fileData)
	if err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	helpers.StatusCreatedData(w, res)

}
