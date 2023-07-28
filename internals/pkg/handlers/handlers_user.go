package handlers

import (
	"net/http"
	"strconv"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/helpers"
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
