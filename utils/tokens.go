package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

var (
	AccessExpiresAt  = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	RefreshExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	IssuedAt         = jwt.NewNumericDate(time.Now())
	NotBefore        = jwt.NewNumericDate(time.Now())
	Secret           = os.Getenv("JWT_SECRET")
)

type Claims struct {
	Username    string
	ID          string
	AccessLevel int
	jwt.RegisteredClaims
}

type TokenDetail struct {
	UserID      string
	Username    string
	AccessLevel int
	Token       *string
	TokenID     string
	ExpiresAt   time.Time
	Subject     string
}

type Token struct {
	AccessToken  *TokenDetail
	RefreshToken *TokenDetail
}

type LoginResponse struct {
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateLoginResponse(id, username string, accessLevel int) (*LoginResponse, *Token, error) {
	tokenID := uuid.NewV4().String()
	accessTokenDetail, err := GenerateAccessToken(id, username, tokenID, accessLevel)
	if err != nil {
		return nil, nil, errors.New("error generating access token")
	}
	refreshTokenDetail, err := GenerateRefreshToken(id, username, tokenID, accessLevel)
	if err != nil {
		return nil, nil, errors.New("error generating refresh token")
	}
	return &LoginResponse{
			Username:     username,
			AccessToken:  *accessTokenDetail.Token,
			RefreshToken: *refreshTokenDetail.Token,
		}, &Token{
			AccessToken:  accessTokenDetail,
			RefreshToken: refreshTokenDetail,
		}, nil

}

func GenerateAccessToken(id, username, tokenID string, accessLevel int) (*TokenDetail, error) {
	access_claims := &Claims{
		ID:          id,
		Username:    username,
		AccessLevel: accessLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: AccessExpiresAt,
			ID:        tokenID,
			IssuedAt:  IssuedAt,
			NotBefore: NotBefore,
			Subject:   "access_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, access_claims)
	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		return nil, err
	}
	return &TokenDetail{
		UserID:      id,
		Username:    username,
		AccessLevel: accessLevel,
		Token:       &tokenString,
		TokenID:     access_claims.RegisteredClaims.ID,
		ExpiresAt:   access_claims.RegisteredClaims.ExpiresAt.Time,
		Subject:     access_claims.RegisteredClaims.Subject,
	}, nil
}
func GenerateRefreshToken(id, username, tokenID string, accessLevel int) (*TokenDetail, error) {
	refresh_token := &Claims{
		ID:          id,
		Username:    username,
		AccessLevel: accessLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: RefreshExpiresAt,
			ID:        tokenID,
			IssuedAt:  IssuedAt,
			NotBefore: NotBefore,
			Subject:   "refresh_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_token)
	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		return nil, err
	}
	return &TokenDetail{
		UserID:      id,
		Username:    username,
		AccessLevel: accessLevel,
		Token:       &tokenString,
		TokenID:     refresh_token.RegisteredClaims.ID,
		ExpiresAt:   refresh_token.RegisteredClaims.ExpiresAt.Time,
		Subject:     refresh_token.RegisteredClaims.Subject,
	}, nil
}
