package entity

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type ITenant interface {
	Create(*TenantResponse) (*TenantResponse, error)
	Get() ([]TenantResponse, error)
	GetById(id *string) (*TenantResponse, error)
	Update(data *TenantResponse) (*TenantResponse, error)
	Delete(id *string) error
	GetByFilterMany(ctx context.Context, filter []QueryDB) ([]TenantResponse, error)
	GetByFilterOne(ctx context.Context, filter []QueryDB) (*TenantResponse, error)
}

type TenantResponse struct {
	Name      string       `json:"name" binding:"required,min=3,max=200" firestore:"name"`
	Alias     string       `json:"alias,omitempty" firestore:"alias"`
	OwnerID   string       `json:"owner_id" binding:"required" firestore:"owner_id"`
	Users     []string     `json:"users,omitempty" firestore:"users"`
	Plan      PlanResponse `json:"plan" binding:"required" firestore:"plan"`
	ID        string       `json:"id"  firestore:"id"`
	CreatedAt time.Time    `json:"created_at" firestore:"create_at"`
	UpdatedAt time.Time    `json:"updated_at" firestore:"update_at"`
	Wallets   []string     `json:"wallets,omitempty" firestore:"wallets"`
}

func NewTenant(tenant *TenantResponse) (*TenantResponse, error) {

	if tenant == nil || reflect.DeepEqual(*tenant, TenantResponse{}) {
		return nil, errors.New("tenant is required")
	}

	if tenant.ID == "" {
		id, _ := uuid.NewV7()
		tenant.ID = id.String()
	}

	t := &TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Alias:     tenant.Alias,
		OwnerID:   tenant.OwnerID,
		Plan:      tenant.Plan,
		Users:     tenant.Users,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Wallets:   tenant.Wallets,
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TenantResponse) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("tenant name is required")
	}

	if !utils.IsValidEmail(t.Name) {
		return fmt.Errorf("tenant name is invalid. Tenant name must be a valid email")
	}

	err := utils.ValidateUUID(&t.OwnerID)
	if err != nil {
		return errors.New("invalid owner id")
	}

	if t.Plan.ID != "" {
		err = utils.ValidateUUID(&t.Plan.ID)
		if err != nil {
			return errors.New("invalid plan id")
		}
	}

	if len(t.Wallets) > 0 {
		for _, w := range t.Wallets {
			err = utils.ValidateUUID(&w)
			if err != nil {
				return errors.New("invalid wallet id")
			}
		}
	}

	err = utils.ValidateUUID(&t.ID)
	if err != nil {
		return errors.New("invalid tenant id")
	}

	return nil
}

func (c *TenantResponse) IsEmpty(data *TenantResponse) bool {
	return data == nil || reflect.DeepEqual(*data, TenantResponse{})
}
