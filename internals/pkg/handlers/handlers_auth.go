package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/helpers"
	"github.com/ishanshre/GoFileServerAPI/internals/pkg/models"
	"github.com/ishanshre/GoFileServerAPI/utils"
)

func (h *handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	newUser := &models.CreateUser{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		helpers.StatusBadRequest(w, "error in parsing json/no json data")
		return
	}
	if err := validate.Struct(newUser); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			if err.Tag() == "containsany" {
				helpers.StatusBadRequest(w, fmt.Sprintf("%s must have at least one special characters from: %v", fieldName, err.Param()))
				return
			}
			helpers.StatusBadRequest(w, fmt.Sprintf("%s must have at least one %s %v characters", fieldName, err.Tag(), err.Param()))
		}
		return
	}
	if err := h.mg.UsernameExists(newUser.Username); err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	if err := h.mg.EmailExists(newUser.Email); err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	hashPassword, err := utils.GeneratePassword(newUser.Password)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	user := &models.User{
		Username:    newUser.Username,
		Email:       newUser.Email,
		Password:    hashPassword,
		AccessLevel: 1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	getUser, err := h.mg.CreateUser(user)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.StatusCreatedData(w, getUser)
}
