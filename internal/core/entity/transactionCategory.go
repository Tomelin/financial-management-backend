package entity

import (
	"context"
	"reflect"
	"strconv"

	"github.com/google/uuid"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type ITransactionCategoryRepository interface {
	Create(ctx context.Context, category *TransactionCategory) (*TransactionCategory, *ModuleError)
	Get(ctx context.Context, walletID *string) ([]TransactionCategory, *ModuleError)
	GetById(ctx context.Context, id *string) (*TransactionCategory, *ModuleError)
	Update(ctx context.Context, category *TransactionCategory) (*TransactionCategory, *ModuleError)
	Delete(ctx context.Context, id *string) *ModuleError
	GetByFilterMany(ctx context.Context, filter []QueryDBClause) ([]TransactionCategory, *ModuleError)
	GetByFilterOne(ctx context.Context, filter []QueryDB) (*TransactionCategory, *ModuleError)
}

type ITransactionCategory interface {
	Create(ctx context.Context, email *string, category *TransactionCategory) (*TransactionCategory, *ModuleError)
	Get(ctx context.Context, email *string, walletID *string) ([]TransactionCategory, *ModuleError)
	GetById(ctx context.Context, email *string, id *string) (*TransactionCategory, *ModuleError)
	Update(ctx context.Context, email *string, category *TransactionCategory) (*TransactionCategory, *ModuleError)
	Delete(ctx context.Context, email *string, id *string) *ModuleError
	GetByFilterMany(ctx context.Context, email *string, filter []QueryDB) ([]TransactionCategory, *ModuleError)
	GetByFilterOne(ctx context.Context, email *string, filter []QueryDB) (*TransactionCategory, *ModuleError)
}

// TransactionCategory represents the response of a transaction category
type TransactionCategory struct {
	ID       string `json:"id" firestore:"id"`
	Name     string `json:"name" binding:"required" firestore:"name"`
	Default  string `json:"default" binding:"required" firestore:"default"`
	TenantID string `json:"tenant_id" firestore:"tenant_id"`
	WalletID string `json:"wallet_id" firestore:"wallet_id"`
}

// NewTransactionCategory creates a new transaction category
func NewTransactionCategory(category *TransactionCategory) (*TransactionCategory, *ModuleError) {

	id, err := uuid.NewV7()
	if err != nil {
		return nil, &ModuleError{
			Module: "transaction category",
			Method: "new",
			Layer:  string(ApplicationLayerEntity),
			Code:   ResponseCodeBadRequest,
			Err:    err.Error(),
		}
	}

	t := &TransactionCategory{
		ID:       id.String(),
		Name:     category.Name,
		Default:  category.Default,
		TenantID: category.TenantID,
		WalletID: category.WalletID,
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	return t, nil
}

// Validate validates the transaction category
func (t *TransactionCategory) Validate() *ModuleError {

	walletDefault, err := strconv.ParseBool(t.Default)
	if err != nil {
		return &ModuleError{
			Module: "transaction_category",
			Method: "validate",
			Layer:  string(ApplicationLayerEntity),
			Code:   ResponseCodeBadRequest,
			Err:    err.Error(),
		}
	}

	if t.IsEmpty(t) {
		return &ModuleError{
			Module: "transaction_category",
			Method: "validate",
			Layer:  string(ApplicationLayerEntity),
			Code:   ResponseCodeBadRequest,
			Err:    "transaction category is required",
		}
	}

	if t.Name == "" {
		return &ModuleError{
			Module: "transaction_category",
			Method: "validate",
			Layer:  string(ApplicationLayerEntity),
			Code:   ResponseCodeBadRequest,
			Err:    "name is required",
		}
	}

	if err := utils.ValidateUUID(&t.ID); err != nil {
		return &ModuleError{
			Module: "transaction_category",
			Method: "validate",
			Layer:  string(ApplicationLayerEntity),
			Code:   ResponseCodeBadRequest,
			Err:    err.Error(),
		}
	}

	if !walletDefault && (t.WalletID != "" || t.TenantID != "") {
		if t.TenantID != "" {
			if err := utils.ValidateUUID(&t.TenantID); err != nil {
				return &ModuleError{
					Module: "transaction_category",
					Method: "validate",
					Layer:  string(ApplicationLayerEntity),
					Code:   ResponseCodeBadRequest,
					Err:    err.Error(),
				}
			}
		}

		if t.WalletID != "" {
			if err := utils.ValidateUUID(&t.WalletID); err != nil {
				return &ModuleError{
					Module: "transaction_category",
					Method: "validate",
					Layer:  string(ApplicationLayerEntity),
					Code:   ResponseCodeBadRequest,
					Err:    err.Error(),
				}
			}
		}

		if (t.WalletID != "" && t.TenantID == "") || (t.WalletID == "" && t.TenantID != "") {
			return &ModuleError{
				Module: "transaction_category",
				Method: "validate",
				Layer:  string(ApplicationLayerEntity),
				Code:   ResponseCodeBadRequest,
				Err:    "wallet and tenant ID cannot be empty",
			}
		}
	}
	return nil
}

func (w *TransactionCategory) IsEmpty(data *TransactionCategory) bool {
	return data == nil || reflect.DeepEqual(*data, TransactionCategory{})
}
