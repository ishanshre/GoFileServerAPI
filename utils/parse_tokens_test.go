package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestVerifyTokenWithClaims(t *testing.T) {
	id := RandomAlphaNum(20)
	username := RandomAlphapet(10)
	tokenID := RandomAlphaNum(12)
	tokenDetail, err := GenerateAccessToken(id, username, tokenID, accessLevel)
	assert.NoError(t, err)
	assert.NotNil(t, tokenDetail)
	assert.Equal(t, tokenDetail.TokenID, tokenID)
	assert.Equal(t, tokenDetail.UserID, id)
	assert.Equal(t, tokenDetail.Username, username)
	assert.NotNil(t, tokenDetail.Token)

	token, err := VerifyTokenWithClaims(*tokenDetail.Token, "access_token")
	assert.NoError(t, err)
	assert.Equal(t, token.Token, tokenDetail.Token)
	assert.Equal(t, token.TokenID, tokenDetail.TokenID)
	assert.Equal(t, token.UserID, tokenDetail.UserID)
	assert.Equal(t, token.Username, tokenDetail.Username)
	assert.Equal(t, token.ExpiresAt, tokenDetail.ExpiresAt)
}

func TestVerifyTokenWithClaims_Failure(t *testing.T) {
	id := RandomAlphaNum(20)
	username := RandomAlphapet(10)
	tokenID := RandomAlphaNum(12)
	tokenDetail, err := GenerateAccessToken(id, username, tokenID, accessLevel)
	assert.NoError(t, err)
	assert.NotNil(t, tokenDetail)
	assert.Equal(t, tokenDetail.TokenID, tokenID)
	assert.Equal(t, tokenDetail.UserID, id)
	assert.Equal(t, tokenDetail.Username, username)
	assert.NotNil(t, tokenDetail.Token)

	token, err := VerifyTokenWithClaims(*tokenDetail.Token, "refresh_token")
	assert.Error(t, err)
	assert.Nil(t, token)
	token, err = VerifyTokenWithClaims("invalid token", "refresh_token")
	assert.Error(t, err)
	assert.Nil(t, token)

	newClaims := &Claims{
		Username: username,
		ID:       id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * -12)),
			Subject:   "access_token",
		},
	}
	newtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	err = ValidateToken(newtoken, newClaims, "access_token")
	assert.Error(t, err)
}
