package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ishanshre/GoFileServerAPI/internals/pkg/helpers"
	"github.com/ishanshre/GoFileServerAPI/utils"
)

const (
	tokenDetailKey helpers.ContextKey = "tokenDetail"
)

func (m *middlewares) JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			helpers.StatusUnauthorized(w, "unauthorized")
			return
		}
		tokenString := strings.Split(bearerToken, " ")
		if len(tokenString) != 2 {
			helpers.StatusUnauthorized(w, "invalid token format")
			return
		}
		if tokenString[0] != "Bearer" {
			helpers.StatusUnauthorized(w, "invalid token format")
			return
		}
		tokenDetail, err := utils.VerifyTokenWithClaims(tokenString[1], "access_token")
		if err != nil {
			helpers.StatusUnauthorized(w, err.Error())
			return
		}
		exists, err := m.redisClient.Exists(m.ctx, tokenDetail.Username).Result()
		if err != nil {
			helpers.StatusInternalServerError(w, err.Error())
			return
		}
		if exists == 0 {
			helpers.StatusUnauthorized(w, "token already revoked")
			return
		}
		data, err := m.redisClient.Get(m.ctx, tokenDetail.Username).Result()
		if err != nil {
			helpers.StatusInternalServerError(w, err.Error())
			return
		}
		token := &utils.Token{}
		if err := json.Unmarshal([]byte(data), token); err != nil {
			helpers.StatusInternalServerError(w, err.Error())
			return
		}
		if tokenDetail.TokenID != token.AccessToken.TokenID {
			helpers.StatusUnauthorized(w, "token already revoked")
			return
		}
		ctx := context.WithValue(r.Context(), tokenDetailKey, tokenDetail)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (m *middlewares) CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
		if tokenDetail.AccessLevel != 0 {
			helpers.StatusUnauthorized(w, "unauthorized access")
			return
		}
		next.ServeHTTP(w, r)
	})
}
