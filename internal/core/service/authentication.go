package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/internal/core/repository"
	"github.com/synera-br/financial-management/src/backend/pkg/logger"
)

type AuthorizationSvc struct {
	repo repository.IAuthorizationRepo
	log  logger.Logger
}

func NewAuthorizationSvc(repo repository.IAuthorizationRepo, l logger.Logger) (entity.IAuthorization, error) {
	if repo == nil {
		return nil, l.Error(&logger.Message{
			Body: "repository is required",
			Code: logger.ResponseCodeInternalServer,
		})
	}
	return &AuthorizationSvc{
		repo: repo,
		log:  l,
	}, nil
}

func (a *AuthorizationSvc) GenerateTokenJWT(ctx context.Context, token *entity.AuthorizationClaims, user *entity.AccountUser) (*string, error) {

	claimsID, _ := uuid.NewV7()
	claims := &entity.AuthorizationClaims{
		UserID:   user.ID,
		Username: user.Name,
		Email:    user.Email,
		// Roles:     user.Roles,
		IsRevoked: false,
		StandardClaims: jwt.StandardClaims{
			// Audience  string `json:"aud,omitempty"`
			Id:        claimsID.String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "financial-app",
			Subject:   user.Email,
			ExpiresAt: time.Now().Add(time.Hour * 4).Unix(),
		},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenJWT.SignedString([]byte(entity.SecretTokenJWT))
	if err != nil {
		return nil, a.log.Error(&logger.Message{
			Body: fmt.Sprintf("error signed token: %s", err.Error()),
			Code: logger.ResponseCodeInternalServer,
		})
	}

	a.repo.GenerateTokenJWT(ctx, token, user)

	return &tokenString, nil
}

func (a *AuthorizationSvc) ValidateTokenJWT(ctx context.Context, token string) (*entity.AccountUser, error) {
	return a.repo.ValidateTokenJWT(ctx, token)
}

func (a *AuthorizationSvc) RevokeTokenJWT(ctx context.Context, token *string) error {
	return a.repo.RevokeTokenJWT(ctx, token)
}

func (a *AuthorizationSvc) RefreshTokenJWT(ctx context.Context, token *string) (*string, error) {
	return a.repo.RefreshTokenJWT(ctx, token)
}

func (a *AuthorizationSvc) ParseTokenJWT(ctx context.Context, token string) error {
	return a.repo.ParseTokenJWT(ctx, token)
}

func (a *AuthorizationSvc) StoreTokenJWT(ctx context.Context, token []byte, userId *string) error {
	return a.repo.StoreTokenJWT(ctx, token, userId)
}
