package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Tomelin/financial-management-backend/internal/core/entity"
	"github.com/Tomelin/financial-management-backend/internal/core/repository"
	"github.com/Tomelin/financial-management-backend/pkg/utils"
)

type ITenantService interface {
	entity.ITenant
	GetPlan(id *string) (*entity.PlanResponse, error)
	SetPlan(id *string, plan *entity.PlanResponse) error
}

type TenantSvc struct {
	repo repository.ITenantRepo
	plan entity.IPlan
}

func NewTenantService(repo repository.ITenantRepo, plan entity.IPlan) (ITenantService, error) {

	if repo == nil {
		return nil, errors.New("repository is required")
	}

	if plan == nil {
		return nil, errors.New("plan is required")
	}

	return &TenantSvc{repo: repo, plan: plan}, nil
}

func (u *TenantSvc) Create(tenant *entity.TenantResponse) (*entity.TenantResponse, error) {

	if tenant == nil {
		return nil, errors.New("tenant is required")
	}

	if tenant.Plan.IsEmpty(&tenant.Plan) {
		name := "bronze"
		resultPlan, _ := u.plan.GetByFilterOne("name", &name)
		log.Println(resultPlan)
	}
	if err := tenant.Validate(); err != nil {
		return nil, err
	}

	ctx := context.Background()
	filter := entity.QueryDB{
		Key:       "name",
		Value:     tenant.Name,
		Condition: string(entity.QueryFirebaseEqual),
	}
	res, err := u.GetByFilterOne(ctx, []entity.QueryDB{filter})
	if err != nil && err.Error() != "not found" {
		return nil, err
	}

	// user, err := u.user.GetById(ctx, &tenant.OwnerID)
	// if err != nil {
	// 	return nil, err
	// }

	// if user == nil || user.Email == "" {
	// 	return nil, errors.New("user not found")
	// }
	if !tenant.Plan.IsEmpty(&tenant.Plan) {
		if err := tenant.Plan.Validate(); err != nil {
			return nil, err
		}
	}

	if res.IsEmpty(res) {
		response, err := u.repo.Create(tenant)
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	return nil, errors.New("tenant already exists")
}

func (u *TenantSvc) Get() ([]entity.TenantResponse, error) {

	return u.repo.Get()
}

func (u *TenantSvc) GetById(id *string) (*entity.TenantResponse, error) {
	if id == nil || *id == "" {
		return nil, errors.New("id cannot be empty")
	}

	if err := utils.ValidateUUID(id); err != nil {
		return nil, err
	}

	data, err := u.repo.GetById(id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}
	return data, err
}

func (u *TenantSvc) Update(data *entity.TenantResponse) (*entity.TenantResponse, error) {
	if data == nil {
		return nil, errors.New("tenant cannot be empty")
	}

	err := utils.ValidateUUID(&data.ID)
	if err != nil {
		return nil, err
	}

	if err := data.Validate(); err != nil {
		return nil, err
	}

	return u.repo.Update(data)
}

func (u *TenantSvc) Delete(id *string) error {
	if id == nil || *id == "" {
		return errors.New("id cannot be empty")
	}

	err := utils.ValidateUUID(id)
	if err != nil {
		return err
	}

	data, err := u.GetById(id)
	if err != nil {
		return err
	}

	if data != nil {
		return u.repo.Delete(id)
	}

	return errors.New("not found")
}

func (u *TenantSvc) GetByFilterMany(ctx context.Context, filter []entity.QueryDB) ([]entity.TenantResponse, error) {

	if len(filter) == 0 {
		return nil, errors.New("filter cannot be empty")
	}

	if filter[0].Key == "" || filter[0].Value == "" {
		return nil, errors.New("filter key and value cannot be empty")
	}

	return u.repo.GetByFilterMany(ctx, filter)
}

func (u *TenantSvc) GetByFilterOne(ctx context.Context, filter []entity.QueryDB) (*entity.TenantResponse, error) {
	if len(filter) == 0 {
		return nil, errors.New("filter cannot be empty")
	}

	if filter[0].Key == "" || filter[0].Value == "" {
		return nil, errors.New("filter key and value cannot be empty")
	}
	return u.repo.GetByFilterOne(ctx, filter)
}

func (u *TenantSvc) SetPlan(id *string, plan *entity.PlanResponse) error {

	if id == nil || *id == "" {
		return errors.New("id is required")
	}

	if plan.IsEmpty(plan) {
		return errors.New("plan is required")
	}

	if err := plan.Validate(); err != nil {
		return err
	}

	if err := plan.Validate(); err != nil {
		return err
	}

	data, err := u.GetById(id)
	if err != nil {
		return err
	}

	data.Plan = *plan

	_, err = u.Update(data)
	if err != nil {
		return err
	}

	return nil
}

func (u *TenantSvc) GetPlan(id *string) (*entity.PlanResponse, error) {
	if id == nil || *id == "" {
		return nil, errors.New("id is required")
	}

	if utils.ValidateUUID(id) != nil {
		return nil, errors.New("id is invalid")
	}

	data, err := u.GetById(id)
	if err != nil {
		return nil, err
	}

	return &data.Plan, nil
}
