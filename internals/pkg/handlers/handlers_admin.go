package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/helpers"
	"github.com/ishanshre/GoFileServerAPI/utils"
)

func (h *handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	users, err := h.mg.GetUsers(limit, page)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Message{
			MessageStatus: "error",
			Message:       err.Error(),
		})
		return
	}
	helpers.WriteJson(w, http.StatusOK, helpers.Message{
		MessageStatus: "success",
		Data:          users,
	})
}

func (h *handlers) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if err := h.mg.UsernameExists(username); err == nil {
		helpers.StatusBadRequest(w, "user does not exists")
		return
	}
	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
	if username == tokenDetail.Username {
		if err := h.redisClient.Del(h.ctx, username).Err(); err != nil {
			helpers.StatusInternalServerError(w, err.Error())
			return
		}
	}

	if err := h.mg.DeleteUser(username); err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.StatusOk(w, "user deleted")
}

func (h *handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.mg.GetUserByUsername(username)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.StatusOkData(w, user)
}
