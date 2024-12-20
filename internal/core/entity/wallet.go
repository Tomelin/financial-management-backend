package entity

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type IWallet interface {
	Create(ctx context.Context, userId *string, wallet *WalletResponse) (*WalletResponse, *ModuleError)
	Get(ctx context.Context, userId *string) ([]WalletResponse, *ModuleError)
	GetWalletByIdAndUserID(ctx context.Context, userId *string, walletId *string) (*WalletResponse, *ModuleError)
	GetByID(ctx context.Context, walletId *string) (*WalletResponse, *ModuleError)
	Update(ctx context.Context, userId *string, data *WalletResponse) (*WalletResponse, *ModuleError)
	UpdateBalance(ctx context.Context, walletID *string, balance *float64) (*WalletResponse, *ModuleError)
	Delete(ctx context.Context, userId *string, walletId *string) *ModuleError
	GetByFilterMany(ctx context.Context, userId *string, filter []QueryDB) ([]WalletResponse, *ModuleError)
	GetByFilterOne(ctx context.Context, userId *string, filter []QueryDB) (*WalletResponse, *ModuleError)
}

type WalletResponse struct {
	mu sync.Mutex // Mutex para garantir segurança em acessos concorrentes

	ID                string    `json:"id" firestore:"id"`
	Name              string    `json:"name" firestore:"name"`
	Description       string    `json:"description" firestore:"description"`
	OwnerID           string    `json:"owner_id" binding:"required" firestore:"owner_id"`
	TenantID          string    `json:"tenant_id" binding:"required" firestore:"tenant_id"`
	Balance           float64   `json:"balance" firestore:"balance"`
	Currency          string    `json:"currency" firestore:"currency"`
	SharedWithTenants []string  `json:"shared_with_tenants" firestore:"shared_with_tenants"`
	CreatedAt         time.Time `json:"createdAt" firestore:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" firestore:"updated_at"`
}

func NewWallet(w *WalletResponse) (*WalletResponse, *ModuleError) {

	if w.IsEmpty(w) {
		return nil, Error("wallet is required", "wallet", "NewWallet", ApplicationLayerEntity, ResponseCodeBadRequest)
	}

	if w.ID == "" {
		id, _ := uuid.NewV7()
		w.ID = id.String()
	}

	wallet := &WalletResponse{
		ID:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		OwnerID:     w.OwnerID,
		TenantID:    w.TenantID,
		Balance:     0,
		Currency:    w.Currency,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := wallet.Validate(); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (w *WalletResponse) Validate() *ModuleError {

	if w.IsEmpty(w) {
		return Error("wallet cannot be empty", "wallet", "Validate", ApplicationLayerEntity, ResponseCodeBadRequest)
	}

	if w.Name == "" {
		w.Name = "MyWallet"
	}

	if w.Currency == "" {
		w.Currency = "BRL"
	}

	ownerId := w.OwnerID
	if err := utils.ValidateUUID(&ownerId); err != nil {
		return Error(err.Error(), "wallet", "Validate", ApplicationLayerEntity, ResponseCodeBadRequest)
	}

	tenantId := w.TenantID
	if err := utils.ValidateUUID(&tenantId); err != nil {
		return Error(err.Error(), "wallet", "Validate", ApplicationLayerEntity, ResponseCodeBadRequest)
	}

	if w.ID == "" {
		return Error("id is required", "wallet", "Validate", ApplicationLayerEntity, ResponseCodeBadRequest)
	}

	id := w.ID
	if err := utils.ValidateUUID(&id); err != nil {
		return Error(err.Error(), "wallet", "Validate", ApplicationLayerEntity, ResponseCodeBadRequest)
	}

	if w.UpdatedAt == (time.Time{}) {
		w.UpdatedAt = time.Now()
	}

	return nil
}

func (w *WalletResponse) SetBalance(data float64) error {

	balance, err := strconv.ParseFloat(fmt.Sprintf("%.2f", data), 64)
	if err != nil {
		return err
	}

	w.Balance = balance
	return nil
}

func (w *WalletResponse) IsEmpty(data *WalletResponse) bool {
	return data == nil || reflect.DeepEqual(*data, WalletResponse{})
}

func (w *WalletResponse) WalletByPlan(data *string) int {

	switch *data {
	case "silver":
		return 3
	case "gold":
		return -1
	default:
		return 1
	}
}

func (w *WalletResponse) SharedByPlan(data *string) int {

	switch *data {
	case "silver":
		return 2
	case "gold":
		return -1
	default:
		return 1
	}
}

func (w *WalletResponse) SetUpdate() {
	w.UpdatedAt = time.Now()
}

func (w *WalletResponse) Share(tenantID string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.SharedWithTenants = append(w.SharedWithTenants, tenantID)
}

func (w *WalletResponse) Unshare(tenantID string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i, id := range w.SharedWithTenants {
		if id == tenantID {
			w.SharedWithTenants = append(w.SharedWithTenants[:i], w.SharedWithTenants[i+1:]...)
			return
		}
	}
}

func (w *WalletResponse) IsSharedWith(tenantID string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, id := range w.SharedWithTenants {
		if id == tenantID {
			return true
		}
	}
	return false
}
