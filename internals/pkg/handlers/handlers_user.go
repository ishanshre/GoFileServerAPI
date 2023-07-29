package handlers

import (
	"net/http"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/helpers"
	"github.com/ishanshre/GoFileServerAPI/utils"
)

func (h *handlers) GetMe(w http.ResponseWriter, r *http.Request) {
	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
	if tokenDetail.Username == "" {
		helpers.StatusUnauthorized(w, "user not authorized")
		return
	}
	user, err := h.mg.GetUserByUsername(tokenDetail.Username)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.StatusOkData(w, user)
}

func (h *handlers) DeleteMe(w http.ResponseWriter, r *http.Request) {
	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
	if err := h.mg.UsernameExists(tokenDetail.Username); err == nil {
		helpers.StatusBadRequest(w, "user does not exists")
		return
	}

	if err := h.redisClient.Del(h.ctx, tokenDetail.Username).Err(); err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}

	if err := h.mg.DeleteUser(tokenDetail.Username); err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.StatusOk(w, "user deleted")
}
