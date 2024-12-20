package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type IWalletSvc interface {
	entity.IWallet
	SetBalance(id *string, balance *float64) error
	PlanValidate(id *string) error
}

type WalletSvc struct {
	repo   entity.IWallet
	tenant entity.ITenant
	owner  entity.IUser
}

// NewWalletSvc creates a new WalletSvc
// It requires a repository, a tenant, and a user
// It returns a WalletSvc and an error
func NewWalletSvc(repo entity.IWallet, tenant entity.ITenant, user entity.IUser) (entity.IWallet, *entity.ModuleError) {
	if repo == nil {
		return nil, entity.Error("repo is required", "wallet", "inicialization", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if tenant == nil {
		return nil, entity.Error("tenant is required", "wallet", "inicialization", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	return &WalletSvc{
		repo:   repo,
		tenant: tenant,
		owner:  user,
	}, nil
}

func (w *WalletSvc) Create(ctx context.Context, userId *string, wallet *entity.WalletResponse) (*entity.WalletResponse, *entity.ModuleError) {

	// Validate if the wallet is empty
	if wallet == nil {
		return nil, entity.Error("wallet required", "wallet", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if mErr := wallet.Validate(); mErr != nil {
		return nil, mErr
	}

	if userId == nil || *userId == "" {
		return nil, entity.Error("required email of user", "wallet", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	// Validate if the owner ID is valid
	ownerId := wallet.OwnerID
	if err := utils.ValidateUUID(&ownerId); err != nil {
		return nil, entity.Error(err.Error(), "wallet", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	// Validate if the tenant ID is valid
	tenantId := wallet.TenantID
	if err := utils.ValidateUUID(&tenantId); err != nil {
		return nil, entity.Error(err.Error(), "wallet", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	// Validate if the tenant exists
	tenant, err := w.tenant.GetById(&tenantId)
	if err != nil || tenant == nil {
		return nil, entity.Error(err.Error(), "wallet", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	// Validate if the tenant exists
	owner, err := w.owner.GetById(ctx, &ownerId)
	if err != nil || owner == nil {
		return nil, entity.Error(err.Error(), "wallet", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	// Validate if the user is authorized to create a wallet
	if tenant.OwnerID != ownerId && tenant.ID != tenantId && *userId != tenant.OwnerID && owner.ID != *userId {
		return nil, entity.Error("unathorized user", "wallet", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	// Create a query to filter the wallets by owner ID and tenant ID
	queries := []entity.QueryDB{
		{Key: "owner_id", Value: wallet.OwnerID, Condition: "=="},
		{Key: "tenant_id", Value: wallet.TenantID, Condition: "=="},
	}

	// Get the wallets by filter
	data, mErr := w.repo.GetByFilterMany(ctx, userId, queries)
	if mErr != nil {
		return nil, mErr
	}

	dataCount := len(data)
	if dataCount > 0 {
		if (dataCount == 1 && tenant.Plan.Name == "bronze") || (dataCount == 2 && tenant.Plan.Name == "silver") {
			return nil, entity.Error("limit of wallets reached", "wallet", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
		}

		for _, getWallet := range data {
			if wallet.Name == getWallet.Name {
				return nil, entity.Error("limit of wallets reached", "wallet", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
			}
		}
	}

	// Create the wallet
	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	result, mErr := w.repo.Create(ctx, userId, wallet)
	if mErr != nil {
		return nil, mErr
	}

	return result, nil
}

func (w *WalletSvc) Get(ctx context.Context, userId *string) ([]entity.WalletResponse, *entity.ModuleError) {
	return w.repo.Get(ctx, userId)
}

func (w *WalletSvc) GetWalletByIdAndUserID(ctx context.Context, email *string, walletId *string) (*entity.WalletResponse, *entity.ModuleError) {

	if walletId == nil || *walletId == "" {
		return nil, entity.Error("id cannot be empty", "wallet", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if err := utils.ValidateUUID(walletId); err != nil {
		return nil, entity.Error(err.Error(), "wallet", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if ok := utils.IsValidEmail(*email); !ok {
		return nil, entity.Error("email is invalid", "wallet", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	user, err := w.owner.GetByEmail(ctx, email)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if user == nil || user.Email == "" || user.Email != *email {
		return nil, entity.Error("user not found", "wallet", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	return w.repo.GetWalletByIdAndUserID(ctx, &user.ID, walletId)
}

func (w *WalletSvc) Update(ctx context.Context, email *string, data *entity.WalletResponse) (*entity.WalletResponse, *entity.ModuleError) {

	if data == nil {
		return nil, entity.Error("wallet cannot be empty", "wallet", "Update", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if mErr := data.Validate(); mErr != nil {
		return nil, mErr
	}

	user, mErr := w.getUser(ctx, email)
	if mErr != nil {
		return nil, mErr
	}

	if user.ID != data.OwnerID {
		return nil, entity.Error("unauthorized user", "wallet", "Update", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	tenantID := data.TenantID
	tenant, err := w.tenant.GetById(&tenantID)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if tenant == nil {
		return nil, entity.Error("tenant not found", "wallet", "Update", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	owner, err := w.owner.GetById(ctx, &data.OwnerID)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if owner == nil {
		return nil, entity.Error("owner not found", "wallet", "Update", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	return w.repo.Update(ctx, &user.ID, data)
}

func (w *WalletSvc) Delete(ctx context.Context, email *string, id *string) *entity.ModuleError {

	if err := utils.ValidateUUID(id); err != nil {
		return entity.Error(err.Error(), "wallet", "Delete", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	user, mErr := w.getUser(ctx, email)
	if mErr != nil {
		return mErr
	}

	wallet, mErr := w.GetWalletByIdAndUserID(ctx, email, id)
	if mErr != nil {
		return mErr
	}

	// if wallet == nil || *wallet == (entity.WalletResponse{}) {
	if wallet == nil {
		return entity.Error("wallet not found", "wallet", "Delete", entity.ApplicationLayerService, entity.ResponseCodeNotFound)
	}

	return w.repo.Delete(ctx, &user.ID, id)
}

func (w *WalletSvc) GetByFilterMany(ctx context.Context, email *string, filter []entity.QueryDB) ([]entity.WalletResponse, *entity.ModuleError) {

	if len(filter) == 0 {
		return nil, entity.Error("filter cannot be empty", "wallet", "GetByFilterMany", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	user, mErr := w.getUser(ctx, email)
	if mErr != nil {
		return nil, mErr
	}

	if user == nil {
		return nil, entity.Error("user not found", "wallet", "GetByFilterMany", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	return w.repo.GetByFilterMany(ctx, &user.ID, filter)
}

func (w *WalletSvc) GetByFilterOne(ctx context.Context, email *string, filter []entity.QueryDB) (*entity.WalletResponse, *entity.ModuleError) {

	if len(filter) == 0 {
		return nil, entity.Error("filter cannot be empty", "wallet", "GetByFilterOne", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	user, mErr := w.getUser(ctx, email)
	if mErr != nil {
		return nil, mErr
	}

	if user == nil {
		return nil, entity.Error("user not found", "wallet", "GetByFilterOne", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	return w.repo.GetByFilterOne(ctx, &user.ID, filter)
}

func (w *WalletSvc) GetByID(ctx context.Context, walletId *string) (*entity.WalletResponse, *entity.ModuleError) {

	if walletId == nil || *walletId == "" {
		return nil, entity.Error("id cannot be empty", "wallet", "GetWalletByID", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if err := utils.ValidateUUID(walletId); err != nil {
		return nil, entity.Error(err.Error(), "wallet", "GetWalletByID", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	result, mErr := w.repo.GetByID(ctx, walletId)
	if mErr != nil {
		return nil, mErr
	}

	// if result == nil || *result == (entity.WalletResponse{}) {
	if result == nil {
		return nil, entity.Error("wallet not found", "wallet", "GetWalletByID", entity.ApplicationLayerService, entity.ResponseCodeNotFound)
	}

	return result, nil
}

func (w *WalletSvc) UpdateBalance(ctx context.Context, walletID *string, balance *float64) (*entity.WalletResponse, *entity.ModuleError) {
	return nil, nil
}

func (w *WalletSvc) userIsValid(ctx context.Context, userId, requesterId *string) error {

	if utils.ValidateUUID(userId) != nil {
		return errors.New("user ID is invalid")
	}

	if utils.ValidateUUID(requesterId) != nil {
		return errors.New("user ID is invalid")
	}

	if requesterId != nil && *requesterId != "" && userId != nil && *userId != "" {
		if *userId != *requesterId {
			return errors.New("unauthorized user")
		}
	}

	user, err := w.owner.GetById(ctx, userId)
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("user %s", entity.ErrNotFound)
	}

	return nil
}

func (w *WalletSvc) getUser(ctx context.Context, email *string) (*entity.AccountUser, *entity.ModuleError) {

	if email == nil || *email == "" {
		return nil, entity.Error("email cannot be empty", "wallet", "Generic", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if ok := utils.IsValidEmail(*email); !ok {
		return nil, entity.Error("email is invalid", "wallet", "Generic", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	user, err := w.owner.GetByEmail(ctx, email)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "Generic", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if user == nil || user.Email == "" || user.Email != *email {
		return nil, entity.Error("user not found", "wallet", "Generic", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	return user, nil

}
