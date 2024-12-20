package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/Tomelin/financial-management-backend/internal/core/entity"
	"github.com/Tomelin/financial-management-backend/pkg/db"
	"github.com/Tomelin/financial-management-backend/pkg/logger"
)

type IAuthorizationRepo interface {
	entity.IAuthorization
}
type AuthorizationRepo struct {
	db  db.FirebaseDatabaseInterface
	log logger.Logger
}

func NewAuthorizationRepo(db db.FirebaseDatabaseInterface, l logger.Logger) (IAuthorizationRepo, error) {

	if db == nil {
		return nil, errors.New("db é obrigatório")
	}

	return &AuthorizationRepo{
		db:  db,
		log: l,
	}, nil
}

func (a *AuthorizationRepo) GenerateTokenJWT(ctx context.Context, token *entity.AuthorizationClaims, user *entity.AccountUser) (*string, error) {

	docRef := a.db.Collection("refresh_tokens").Doc(user.ID)
	_, err := docRef.Set(ctx, *user)
	if err != nil {
		return nil, a.log.Error(&logger.Message{
			Body: fmt.Sprintf("erro ao salvar o refresh token no Firestore GenerateTokenJWT: %s", err.Error()),
			Code: logger.ResponseCodeInternalServer})

	}
	return &token.StandardClaims.Id, nil
}

func (a *AuthorizationRepo) ValidateTokenJWT(ctx context.Context, token string) (*entity.AccountUser, error) {
	return nil, nil
}

func (a *AuthorizationRepo) RevokeTokenJWT(ctx context.Context, tokenId *string) error {
	docRef := a.getToken(ctx, tokenId)
	if docRef == nil {
		log.Println("docRef is nil", *tokenId)
	} else {
		log.Println("docRef is not nil", *docRef)
	}
	var refreshToken entity.AuthorizationClaims

	err := docRef.DataTo(&refreshToken)
	if err != nil {
		return a.log.Error(&logger.Message{
			Body: fmt.Sprintf("erro ao buscar refresh token: %s", err.Error()),
			Code: logger.ResponseCodeInternalServer})
	}

	if docRef == nil {
		return a.log.Error(&logger.Message{
			Body: fmt.Sprintf("erro ao buscar refresh token: %s", err.Error()),
			Code: logger.ResponseCodeInternalServer})
	}

	if refreshToken.Email == "" {
		return a.log.Error(&logger.Message{
			Body: "refresh token não encontrado",
			Code: logger.ResponseCodeInternalServer})
	}

	refreshToken.ExpiresAt = time.Now().Unix() - 10
	refreshToken.IsRevoked = true
	_, err = docRef.Ref.Update(ctx, []firestore.Update{
		{
			Path:  "is_revoked",
			Value: refreshToken.IsRevoked,
		},
		{
			Path:  "expiresAt",
			Value: refreshToken.ExpiresAt,
		},
	})
	if err != nil {
		return a.log.Error(&logger.Message{
			Body: fmt.Sprintf("erro ao buscar refresh token: %w", err.Error()),
			Code: logger.ResponseCodeInternalServer})
	}

	return nil
}

func (a *AuthorizationRepo) RefreshTokenJWT(ctx context.Context, tokenId *string) (*string, error) {

	docRef := a.db.Collection("refresh_tokens").Where("token", "==", &tokenId).Documents(ctx)
	docsnap, err := docRef.GetAll()
	if err != nil {
		return nil, a.log.Error(&logger.Message{
			Body: fmt.Sprintf("erro ao buscar refresh token: %s", err.Error()),
			Code: logger.ResponseCodeInternalServer})
	}

	if len(docsnap) == 0 {
		return nil, a.log.Error(&logger.Message{
			Body: "refresh token não encontrado",
			Code: logger.ResponseCodeInternalServer})
	}

	var refreshToken entity.AuthorizationClaims
	err = docsnap[0].DataTo(&refreshToken)
	if err != nil {
		return nil, a.log.Error(&logger.Message{
			Body: fmt.Sprintf("erro ao buscar refresh token: %s", err.Error()),
			Code: logger.ResponseCodeInternalServer})
	}
	if refreshToken.IsRevoked {
		return nil, a.log.Error(&logger.Message{
			Body: "refresh token revogado",
			Code: logger.ResponseCodeInternalServer})
	}
	if refreshToken.StandardClaims.VerifyNotBefore(time.Now().Unix(), true) {
		return nil, a.log.Error(&logger.Message{
			Body: "refresh token expirado",
			Code: logger.ResponseCodeInternalServer})
	}
	return nil, nil
}

func (a *AuthorizationRepo) ParseTokenJWT(ctx context.Context, token string) error {
	return nil
}

func (a *AuthorizationRepo) getToken(ctx context.Context, tokenId *string) *firestore.DocumentSnapshot {
	docRef, _ := a.db.Collection("refresh_tokens").Where("id", "==", *tokenId).Documents(ctx).Next()

	return docRef
}

func (a *AuthorizationRepo) StoreTokenJWT(ctx context.Context, token []byte, userId *string) error {
	docRef := a.db.Collection("refresh_tokens").Doc(*userId)
	_, err := docRef.Set(ctx, string(token))
	if err != nil {
		return a.log.Error(&logger.Message{
			Body: fmt.Sprintf("erro ao salvar o refresh token no Firestore StoreTokenJWT: %s", err.Error()),
			Code: logger.ResponseCodeInternalServer})
	}
	return nil
}
