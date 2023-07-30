package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		helpers.StatusBadRequest(w, "No files uploaded by user")
		return
	}

	userPath := filepath.Join(root_path, tokenDetail.Username)
	if err := h.mg.FileNameExists(files[0].Filename); err != nil {
		helpers.StatusBadRequest(w, "file already exists! ignoring it")
		return
	}
	fileName, size, err := helpers.UploadFile(w, userPath, files[0])
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
	fileData.Size = size
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

// func (h *handlers) UploadSingleFile(w http.ResponseWriter, r *http.Request) {
// 	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)

// 	if err := r.ParseMultipartForm(1 * 1024 * 1024); err != nil {
// 		helpers.StatusBadRequest(w, err.Error())
// 		return
// 	}
// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		helpers.StatusBadRequest(w, err.Error())
// 		return
// 	}
// 	defer file.Close()
// 	fileName := header.Filename
// 	userPath := filepath.Join(root_path, tokenDetail.Username)
// 	if err := helpers.UploadFile(w, userPath, fileName, &file); err != nil {
// 		helpers.StatusInternalServerError(w, err.Error())
// 		return
// 	}
// 	id, _ := primitive.ObjectIDFromHex(tokenDetail.UserID)
// 	uploader := &models.UserAccess{}
// 	uploader.Username = tokenDetail.Username
// 	uploader.ID = id
// 	uploader.AccessLevel = tokenDetail.AccessLevel
// 	fileData := &models.File{}
// 	fileData.FilePath = userPath
// 	fileData.Name = fileName
// 	fileData.Extension = filepath.Ext(fileName)
// 	fileData.Size = header.Size
// 	fileData.Uploader = uploader
// 	fileData.UploadedAt = time.Now()
// 	fileData.ModifiedAt = time.Now()
// 	res, err := h.mg.InsertFileData(fileData)
// 	if err != nil {
// 		helpers.StatusBadRequest(w, err.Error())
// 		return
// 	}
// 	helpers.StatusCreatedData(w, res)
// }

func (h *handlers) UploadMultipleFile(w http.ResponseWriter, r *http.Request) {
	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
	if err := r.ParseMultipartForm(1 * 1024 * 1024); err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		helpers.StatusBadRequest(w, "no files upoaded by user")
		return
	}
	userPath := filepath.Join(root_path, tokenDetail.Username)
	fileDatas := []*models.File{}
	errors := []string{}
	for _, fileHeader := range files {
		if err := h.mg.FileNameExists(fileHeader.Filename); err != nil {
			errors = append(errors, fmt.Sprintf("file %s already exists. Ignoring it", fileHeader.Filename))
			continue
		}
		fileName, size, err := helpers.UploadFile(w, userPath, fileHeader)
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
		fileData.Size = size
		fileData.Uploader = uploader
		fileData.UploadedAt = time.Now()
		fileData.ModifiedAt = time.Now()
		res, err := h.mg.InsertFileData(fileData)
		if err != nil {
			helpers.StatusBadRequest(w, err.Error())
			return
		}
		fileDatas = append(fileDatas, res)
	}
	if len(fileDatas) == 0 {
		helpers.StatusBadRequestData(w, errors)
		return
	}
	helpers.StatusCreatedData(w, fileDatas)
}

func (h *handlers) GetAllFilesByUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}
	files, err := h.mg.AllFilesByUser(username, limit, page)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.StatusOkData(w, files)
}

const fileKey helpers.ContextKey = "fileKey"

func (h *handlers) DeleteFilesByUser(w http.ResponseWriter, r *http.Request) {
	fileData := r.Context().Value(fileKey).(*models.File)
	if err := h.db.StartTranscation(); err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}

	defer func() {
		if err := h.db.CommitTranscation(); err != nil {
			h.db.EndTranscation()
			helpers.StatusInternalServerError(w, err.Error())
			return
		}
	}()

	if err := h.mg.FileDelete(fileData.Name); err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}

	relatetivePath := filepath.Join("./", fileData.FilePath, fileData.Name)
	if err := os.Remove(relatetivePath); err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.StatusOk(w, "deleted successfull")

}
