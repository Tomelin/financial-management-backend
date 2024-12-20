package repository

import (
	"context"
	"errors"

	"github.com/Tomelin/financial-management-backend/internal/core/entity"
	"github.com/Tomelin/financial-management-backend/pkg/db"
	"google.golang.org/api/iterator"
)

type ITenantRepo interface {
	entity.ITenant
	GetPlan(id *string) (*entity.PlanResponse, error)
	SetPlan(id *string, plan *entity.PlanResponse) error
}

type TenantRepo struct {
	db db.FirebaseDatabaseInterface
}

func NewTenantRepository(db db.FirebaseDatabaseInterface) (ITenantRepo, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}
	return &TenantRepo{db: db}, nil
}

func (u *TenantRepo) Create(tenant *entity.TenantResponse) (*entity.TenantResponse, error) {
	ctx := context.Background()

	_, err := u.db.Collection("tenants").Doc(tenant.ID).Set(ctx, tenant)
	if err != nil {
		return nil, err
	}

	doc, err := u.db.Collection("tenants").Doc(tenant.ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var TenantResponse entity.TenantResponse
	err = doc.DataTo(&TenantResponse)
	if err != nil {
		return nil, err
	}

	return &TenantResponse, nil
}

func (u *TenantRepo) Get() ([]entity.TenantResponse, error) {
	iter := u.db.Documents(context.Background(), "tenants")

	defer iter.Stop()

	var documents []entity.TenantResponse
	for {
		var tenant entity.TenantResponse
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&tenant)
		if err != nil {
			return nil, err
		}
		tenant.ID = doc.Ref.ID

		documents = append(documents, tenant)
	}

	return documents, nil
}

func (u *TenantRepo) GetById(id *string) (*entity.TenantResponse, error) {
	doc, err := u.db.Collection("tenants").Doc(*id).Get(context.Background())
	if err != nil {
		return nil, err
	}

	var tenant entity.TenantResponse
	err = doc.DataTo(&tenant)
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (u *TenantRepo) Update(data *entity.TenantResponse) (*entity.TenantResponse, error) {
	_, err := u.db.Collection("tenants").Doc(data.ID).Set(context.Background(), data)
	if err != nil {
		return nil, err
	}

	doc, err := u.db.Collection("tenants").Doc(data.ID).Get(context.Background())
	if err != nil {
		return nil, err
	}

	var tenant entity.TenantResponse
	err = doc.DataTo(&tenant)
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (u *TenantRepo) Delete(id *string) error {
	_, err := u.db.Collection("tenants").Doc(*id).Delete(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (u *TenantRepo) GetByFilterMany(ctx context.Context, filter []entity.QueryDB) ([]entity.TenantResponse, error) {
	// ctx, span := u.trace.Trace.Start(ctx, "TenantRepo.GetByFilterMany")
	// defer span.End()
	query := u.db.Collection("tenants").Query
	for _, f := range filter {
		condition := checkFirebaseCondition(&f.Condition)
		if f.Key != "" && f.Value != "" && condition != "" {
			query = query.Where(f.Key, condition, f.Value)
		}
	}

	iter := query.Documents(context.Background())
	defer iter.Stop()
	var tenants []entity.TenantResponse
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.New(err.Error())
		}

		var tenant entity.TenantResponse
		err = doc.DataTo(&tenant)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		tenants = append(tenants, tenant)
	}
	return tenants, nil
}

func (u *TenantRepo) GetByFilterOne(ctx context.Context, filter []entity.QueryDB) (*entity.TenantResponse, error) {

	query := u.db.Collection("tenants").Query
	for _, f := range filter {
		condition := checkFirebaseCondition(&f.Condition)
		if f.Key != "" && f.Value != "" && condition != "" {
			query = query.Where(f.Key, condition, f.Value)
		}
	}

	doc := query.Limit(1).Documents(context.Background())

	result, err := doc.Next()
	if err != nil {
		if err == iterator.Done {
			return nil, nil
		}
		return nil, errors.New(err.Error())
	}

	if result == nil {
		return nil, nil
	}

	var tenant entity.TenantResponse
	err = result.DataTo(&tenant)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &tenant, nil
}

func (u *TenantRepo) GetPlan(id *string) (*entity.PlanResponse, error) {
	return nil, nil
}

func (u *TenantRepo) SetPlan(id *string, plan *entity.PlanResponse) error { return nil }
