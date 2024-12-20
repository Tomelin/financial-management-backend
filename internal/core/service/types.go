package service

import (
	"context"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
)

func ValidateUser(ctx context.Context, user entity.IUser, email *string) *entity.ModuleError {

	if email == nil || *email == "" {
		return entity.Error("email is required", "user", "ValidateUser", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	response, err := user.GetByEmail(ctx, email)
	if response != nil {
		return entity.Error(err.Error(), "user", "ValidateUser", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if response == nil {
		return entity.Error("user not found", "user", "ValidateUser", entity.ApplicationLayerService, entity.ResponseCodeNotFound)
	}

	if response.Email != *email {
		return entity.Error("email unauthorized", "user", "ValidateUser", entity.ApplicationLayerService, entity.ResponseCodeUnauthorized)
	}

	return nil
}
