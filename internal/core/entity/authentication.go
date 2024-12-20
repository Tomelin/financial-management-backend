package entity

import (
	"context"

	"github.com/dgrijalva/jwt-go"
)

type IAuthorization interface {
	GenerateTokenJWT(ctx context.Context, token *AuthorizationClaims, user *AccountUser) (*string, error)
	ValidateTokenJWT(ctx context.Context, token string) (*AccountUser, error)
	RevokeTokenJWT(ctx context.Context, token *string) error
	RefreshTokenJWT(ctx context.Context, token *string) (*string, error)
	ParseTokenJWT(ctx context.Context, token string) error
	StoreTokenJWT(ctx context.Context, token []byte, userId *string) error
}

type AuthorizationClaims struct {
	UserID    string         `json:"user_id" firestore:"user_id"`
	Username  string         `json:"username" firestore:"username"`
	Email     string         `json:"email" firestore:"email"`
	Roles     []AccountRoles `json:"roles" firestore:"roles"`
	IsRevoked bool           `json:"is_revoked" firestore:"is_revoked"`
	jwt.StandardClaims
}

const SecretTokenJWT = "TW9uIERlYyAxNiAyMjoxNDowNyAtMDMgMjAyNAo="
