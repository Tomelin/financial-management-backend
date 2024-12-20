package service

import (
	"context"
	"errors"
	"time"

	"github.com/Tomelin/financial-management-backend/internal/core/entity"
	"github.com/Tomelin/financial-management-backend/pkg/observability"
	"github.com/Tomelin/financial-management-backend/pkg/utils"
)

type UserSvc struct {
	repo   entity.IUser
	tenant ITenantService
	tracer *observability.Tracer
}

func NewUserService(u entity.IUser, tenant ITenantService, tracer *observability.Tracer) (entity.IUser, error) {

	if u == nil || tenant == nil {
		return nil, errors.New("error creating user service")
	}

	svc := &UserSvc{
		repo:   u,
		tenant: tenant,
		tracer: tracer,
	}

	return svc, nil
}

func (u *UserSvc) Create(ctx context.Context, user *entity.AccountUser) (*entity.AccountUser, error) {

	if user == nil || user.IsEmpty(user) {
		return nil, errors.New("user is required")
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	filterUser := []entity.QueryDBClause{
		{
			Clause: entity.QueryClauseAnd,
			Queries: []entity.QueryDB{
				{
					Key:       "email",
					Value:     user.Email,
					Condition: string(entity.QueryFirebaseEqual),
				},
			},
		},
	}

	res, err := u.GetByFilterOne(ctx, filterUser)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}

	filterTenant := []entity.QueryDBClause{
		{
			Clause: entity.QueryClauseAnd,
			Queries: []entity.QueryDB{
				{
					Key:       "name",
					Value:     user.Email,
					Condition: string(entity.QueryFirebaseEqual),
				},
			},
		},
	}

	tenant, err := u.tenant.GetByFilterMany(ctx, filterTenant[0].Queries)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}

	if res.IsEmpty(res) && tenant == nil {
		response, err := u.repo.Create(ctx, user)
		if err != nil {
			return nil, err
		}

		_, err = u.tenant.Create(&entity.TenantResponse{
			Name:    user.Email,
			OwnerID: user.ID,
			ID:      user.TenantID,
			Users:   []string{user.ID},
			Alias:   user.Email,
		})
		if err != nil {
			return nil, err
		}

		return response, nil
	}

	return nil, errors.New("user already exists")
}

func (u *UserSvc) Get(ctx context.Context) ([]entity.AccountUser, error) {

	return u.repo.Get(ctx)
}

func (u *UserSvc) GetById(ctx context.Context, id *string) (*entity.AccountUser, error) {
	if id == nil || *id == "" {
		return nil, errors.New("id cannot be empty")
	}

	if err := utils.ValidateUUID(id); err != nil {
		return nil, err
	}

	return u.repo.GetById(ctx, id)
}

func (u *UserSvc) GetByEmail(ctx context.Context, email *string) (*entity.AccountUser, error) {

	if email == nil || *email == "" {
		return nil, errors.New("email cannot be empty")
	}

	if ok := utils.IsValidEmail(*email); !ok {
		return nil, errors.New("email is invalid")
	}

	user, err := u.repo.GetByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserSvc) Update(ctx context.Context, data *entity.AccountUser) (*entity.AccountUser, error) {
	if data == nil {
		return nil, errors.New("user cannot be empty")
	}

	err := utils.ValidateUUID(&data.ID)
	if err != nil {
		return nil, err
	}

	if err := data.Validate(); err != nil {
		return nil, err
	}

	data.UpdatedAt = time.Now()
	return u.repo.Update(ctx, data)
}

func (u *UserSvc) Delete(ctx context.Context, id *string) error {
	if id == nil || *id == "" {
		return errors.New("id cannot be empty")
	}

	err := utils.ValidateUUID(id)
	if err != nil {
		return err
	}

	data, err := u.GetById(ctx, id)
	if err != nil {
		return err
	}

	if data != nil {
		return u.repo.Delete(ctx, id)
	}

	return errors.New("not found")
}

func (u *UserSvc) GetByFilterMany(ctx context.Context, filter []entity.QueryDBClause) ([]entity.AccountUser, error) {
	return u.repo.GetByFilterMany(ctx, filter)
}

func (u *UserSvc) GetByFilterOne(ctx context.Context, filter []entity.QueryDBClause) (*entity.AccountUser, error) {

	return u.repo.GetByFilterOne(ctx, filter)
}
