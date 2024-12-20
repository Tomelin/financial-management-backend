package service

import (
	"errors"
	"fmt"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type PlanSvc struct {
	repo entity.IPlan
}

func NewPlanService(repo entity.IPlan) (entity.IPlan, error) {

	plan := &PlanSvc{
		repo: repo,
	}
	return plan, nil
}

func (p *PlanSvc) Create(plan *entity.PlanResponse) (*entity.PlanResponse, error) {

	if plan == nil {
		return nil, fmt.Errorf("plan %s", entity.ErrRequired)
	}

	if err := plan.Validate(); err != nil {
		return nil, err
	}

	// Validar se o plano com o mesmo nome já existe
	result, err := p.GetByFilterOne("name", &plan.Name)
	if err != nil {
		return nil, err
	}

	if result != nil && result.Name != "" {
		return nil, fmt.Errorf("plan %s", entity.ErrAlreadyExists)
	}

	response, err := p.repo.Create(plan)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *PlanSvc) Get() ([]entity.PlanResponse, error) {

	return p.repo.Get()
}

func (p *PlanSvc) GetById(id *string) (*entity.PlanResponse, error) {

	if id == nil || *id == "" {
		return nil, errors.New("id cannot be empty")
	}

	if err := utils.ValidateUUID(id); err != nil {
		return nil, err
	}

	result, err := p.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("plan %v", entity.ErrNotFound)
	}

	return result, nil
}

func (p *PlanSvc) Update(data *entity.PlanResponse) (*entity.PlanResponse, error) {

	if data == nil {
		return nil, errors.New("plan cannot be empty")
	}

	err := utils.ValidateUUID(&data.ID)
	if err != nil {
		return nil, err
	}

	if err := data.Validate(); err != nil {
		return nil, err
	}

	get, err := p.GetById(&data.ID)
	if err != nil {
		return nil, err
	}

	if get == nil {
		return nil, fmt.Errorf("plan %v", entity.ErrNotFound)
	}

	data.SetUpdate()
	return p.repo.Update(data)
}

func (p *PlanSvc) Delete(id *string) error {
	if id == nil {
		return fmt.Errorf("id %s", entity.ErrCannotEmpty)
	}

	err := utils.ValidateUUID(id)
	if err != nil {
		return err
	}

	get, err := p.GetById(id)
	if err != nil {
		return err
	}

	if get == nil {
		return fmt.Errorf("plan %s", entity.ErrNotFound)
	}

	return p.repo.Delete(id)
}

func (p *PlanSvc) GetByFilterMany(key string, value *string) ([]entity.PlanResponse, error) {

	if key == "" {
		return nil, fmt.Errorf("key %s", entity.ErrCannotEmpty)
	}

	if value == nil || *value == "" {
		return nil, fmt.Errorf("value %s", entity.ErrCannotEmpty)
	}

	return p.repo.GetByFilterMany(key, value)
}

func (p *PlanSvc) GetByFilterOne(key string, value *string) (*entity.PlanResponse, error) {
	if key == "" {
		return nil, fmt.Errorf("key %s", entity.ErrCannotEmpty)
	}

	if value == nil || *value == "" {
		return nil, fmt.Errorf("value %s", entity.ErrCannotEmpty)
	}

	return p.repo.GetByFilterOne(key, value)
}
