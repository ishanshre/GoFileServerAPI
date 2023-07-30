package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/helpers"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/models"
	"github.com/ishanshre/GoFileServerAPI/utils"
)

func (m *middlewares) FileAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
		filename := filepath.Base(r.URL.Path)
		file, err := m.repository.GetFileByFileName(filename)
		if err != nil {
			helpers.StatusNotFound(w, "Not Found")
			return
		}
		if tokenDetail.UserID != file.Uploader.ID.Hex() || tokenDetail.Username != file.Uploader.Username || tokenDetail.AccessLevel != file.Uploader.AccessLevel {
			log.Println("inisied")
			helpers.StatusUnauthorized(w, "not authorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

const fileKey helpers.ContextKey = "fileKey"

func (m *middlewares) FileOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
		fileData := &models.FileName{}
		if err := json.NewDecoder(r.Body).Decode(&fileData); err != nil {
			helpers.StatusBadRequest(w, "no data/error in parsing data")
			return
		}
		if err := m.repository.FileNameExists(fileData.Name); err == nil {
			helpers.StatusBadRequest(w, fmt.Sprintf("%s does not exists", fileData.Name))
			return
		}
		file, err := m.repository.GetFileByFileName(fileData.Name)
		if err != nil {
			helpers.StatusNotFound(w, "Not Found")
			return
		}
		if tokenDetail.UserID != file.Uploader.ID.Hex() || tokenDetail.Username != file.Uploader.Username || tokenDetail.AccessLevel != file.Uploader.AccessLevel {
			log.Println("inisied")
			helpers.StatusUnauthorized(w, "not authorized")
			return
		}
		ctx := context.WithValue(r.Context(), fileKey, file)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
