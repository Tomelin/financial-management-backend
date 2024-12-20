package middleware

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
)

// OpenTokenJWT
// Parse token and return claims.
func OpenTokenJWT(token *string) (*entity.AuthorizationClaims, error) {

	if token == nil {
		return nil, errors.New("invalid token")
	}

	secretKey := []byte(entity.SecretTokenJWT)
	extractToken, err := jwt.ParseWithClaims(*token, &entity.AuthorizationClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := extractToken.Claims.(*entity.AuthorizationClaims); ok && extractToken.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetTokenJWT
// Receive a token from header and return the token as string.
func GetTokenJWT(c *gin.Context) (*string, error) {

	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, errors.New("token authorization is required")
	}

	return &token, nil
}

// GetClaimsFromToken
// Receive a token and return the claims from the token.
func GetClaimsFromToken(c *gin.Context) (*entity.AuthorizationClaims, error) {

	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, errors.New("token authorization is required")
	}

	claims, err := OpenTokenJWT(&token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// GetEmailFromToken
//
// Receive a token and return the email from the token.
//
// Return a email of user with success or error.
func GetEmailFromToken(c *gin.Context) (*string, error) {

	token, err := GetTokenJWT(c)
	if err != nil {
		return nil, err
	}

	claims, err := OpenTokenJWT(token)
	if err != nil {
		return nil, err
	}

	return &claims.Email, nil
}
