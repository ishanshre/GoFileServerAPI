package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	id := RandomAlphaNum(6)
	username := RandomAlphapet(6)
	tokenID := RandomAlphaNum(10)
	tokenDetail, err := GenerateAccessToken(id, username, tokenID)
	assert.NoError(t, err)
	assert.NotNil(t, tokenDetail)
	assert.Equal(t, id, tokenDetail.UserID)
	assert.Equal(t, username, tokenDetail.Username)
	assert.Equal(t, tokenID, tokenDetail.TokenID)
	assert.NotNil(t, tokenDetail.TokenID)
	assert.NotEmpty(t, tokenDetail.ExpiresAt)
	assert.Equal(t, "access_token", tokenDetail.Subject)
}
func TestGenerateRefreshToken_Failure(t *testing.T) {
	id := RandomAlphaNum(6)
	username := RandomAlphapet(6)
	tokenID := RandomAlphaNum(10)
	tokenDetail, err := GenerateRefreshToken(id, username, tokenID)
	assert.NoError(t, err)
	assert.NotNil(t, tokenDetail)
	assert.Equal(t, id, tokenDetail.UserID)
	assert.Equal(t, username, tokenDetail.Username)
	assert.Equal(t, tokenID, tokenDetail.TokenID)
	assert.NotNil(t, tokenDetail.TokenID)
	assert.NotEmpty(t, tokenDetail.ExpiresAt)
	assert.Equal(t, "refresh_token", tokenDetail.Subject)
}

func TestGenerateLoginResponse(t *testing.T) {
	id := RandomAlphaNum(10)
	username := RandomAlphaNum(6)
	loginResponse, token, err := GenerateLoginResponse(id, username)
	assert.NoError(t, err)
	assert.NotNil(t, loginResponse)
	assert.NotNil(t, token)
	assert.Equal(t, id, token.AccessToken.UserID)
	assert.Equal(t, id, token.RefreshToken.UserID)
	assert.Equal(t, username, token.RefreshToken.Username)
	assert.Equal(t, loginResponse.AccessToken, *token.AccessToken.Token)
	assert.Equal(t, loginResponse.RefreshToken, *token.RefreshToken.Token)
	assert.Equal(t, token.AccessToken.Subject, "access_token")
	assert.Equal(t, token.RefreshToken.Subject, "refresh_token")
}
