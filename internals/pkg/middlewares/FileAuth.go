package middlewares

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/helpers"
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
