package entity

import (
	"errors"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type IPlan interface {
	Create(*PlanResponse) (*PlanResponse, error)
	Get() ([]PlanResponse, error)
	GetById(id *string) (*PlanResponse, error)
	Update(data *PlanResponse) (*PlanResponse, error)
	Delete(id *string) error
	GetByFilterMany(key string, value *string) ([]PlanResponse, error)
	GetByFilterOne(key string, value *string) (*PlanResponse, error)
}

type PlanFeatures struct {
	Name        string `json:"name" binding:"required" firestore:"name"`
	Description string `json:"description" firestore:"description"`
	Count       int    `json:"count" firestore:"count"`
	Unlimited   bool   `json:"unlimited" firestore:"unlimited"`
}

type PlanResponse struct {
	ID          string         `json:"id" firestore:"id"`
	CreatedAt   time.Time      `json:"created_at" firestore:"create_at"`
	UpdatedAt   time.Time      `json:"updated_at" firestore:"update_at"`
	Name        string         `json:"name" binding:"required" firestore:"name"`
	Description string         `json:"description" firestore:"description"`
	Features    []PlanFeatures `json:"features"  firestore:"features"`
	Price       float64        `json:"price" firestore:"price"`
}

func NewPlan(plan *PlanResponse) (*PlanResponse, error) {

	id, _ := uuid.NewV7()

	if plan.IsEmpty(plan) {
		return nil, errors.New("plan is required")
	}
	plan.ID = id.String()
	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()

	if err := plan.Validate(); err != nil {
		return nil, err
	}

	return plan, nil
}

func (p *PlanResponse) Validate() error {

	if p.Name == "" {
		return errors.New("plan name is required")
	}

	if p.Name != "bronze" {
		if len(p.Features) == 0 {
			return errors.New("features are required")
		}
	}

	if err := utils.ValidateUUID(&p.ID); err != nil {
		return errors.New("invalid ID")
	}

	return nil
}

func (p *PlanResponse) IsEmpty(data *PlanResponse) bool {
	return data == nil || reflect.DeepEqual(*data, PlanResponse{})
}

func (p *PlanResponse) SetUpdate() {
	p.UpdatedAt = time.Now()
}

func (p *PlanFeatures) Validate() error {

	if p.Name == "" {
		return errors.New("plan name is required")
	}

	return nil
}

func (p *PlanFeatures) IsEmpty(data *PlanFeatures) bool {
	return data == nil || reflect.DeepEqual(*data, PlanFeatures{})
}
