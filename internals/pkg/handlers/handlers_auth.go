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

const (
	tokenDetailKey helpers.ContextKey = "tokenDetail"
)

func (h *handlers) UserRegister(w http.ResponseWriter, r *http.Request) {
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
		AccessLevel: 2,
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

func (h *handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	user := &models.LoginUser{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.StatusBadRequest(w, "error in parsing json/no data")
		return
	}
	if err := validate.Struct(user); err != nil {
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
	if err := h.mg.UsernameExists(user.Username); err == nil {
		helpers.StatusBadRequest(w, "username does not exists")
		return
	}
	getUser, err := h.mg.GetUserByUsername(user.Username)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	if err := utils.ComparePassword(getUser.Password, user.Password); err != nil {
		helpers.StatusBadRequest(w, "invalid username/password")
		return
	}
	exists, err := h.redisClient.Exists(h.ctx, getUser.Username).Result()
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	if exists == 1 {
		if err := h.redisClient.Del(h.ctx, getUser.Username).Err(); err != nil {
			helpers.StatusInternalServerError(w, err.Error())
			return
		}
	}
	loginResponse, token, err := utils.GenerateLoginResponse(getUser.ID.Hex(), getUser.Username, getUser.AccessLevel)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	tokenJson, err := json.Marshal(token)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	if err := h.redisClient.Set(h.ctx, token.AccessToken.Username, tokenJson, time.Until(token.RefreshToken.ExpiresAt)).Err(); err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	getUser.LastLogin = time.Now()
	_, _ = h.mg.UpdateUser(getUser.Username, getUser)
	helpers.StatusOkData(w, loginResponse)
}

func (h *handlers) UserLogout(w http.ResponseWriter, r *http.Request) {
	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
	if tokenDetail.Username == "" {
		helpers.StatusUnauthorized(w, "user not authorized")
		return
	}
	exists, err := h.redisClient.Exists(h.ctx, tokenDetail.Username).Result()
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	if exists == 0 {
		helpers.StatusBadRequest(w, "user not logged in")
		return
	} else {
		if err := h.redisClient.Del(h.ctx, tokenDetail.Username).Err(); err != nil {
			helpers.StatusInternalServerError(w, err.Error())
			return
		}
	}
	helpers.StatusOk(w, "logout successfull")

}
