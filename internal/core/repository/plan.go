package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/db"
	"google.golang.org/api/iterator"
)

type IPlanRepo interface {
	entity.IPlan
}

type PlanRepo struct {
	db db.FirebaseDatabaseInterface
}

func NewPlanRepository(db db.FirebaseDatabaseInterface) (entity.IPlan, error) {
	if db == nil {
		return nil, fmt.Errorf("db %s", entity.ErrRequired)
	}
	return &PlanRepo{db: db}, nil
}

func (u *PlanRepo) Create(plan *entity.PlanResponse) (*entity.PlanResponse, error) {
	ctx := context.Background()

	_, err := u.db.Collection("plans").Doc(plan.ID).Set(ctx, plan)
	if err != nil {
		return nil, err
	}

	doc, err := u.db.Collection("plans").Doc(plan.ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var PlanResponse entity.PlanResponse
	err = doc.DataTo(&PlanResponse)
	if err != nil {
		return nil, err
	}

	return &PlanResponse, nil
}

func (u *PlanRepo) Get() ([]entity.PlanResponse, error) {

	iter := u.db.Documents(context.Background(), "plans")

	defer iter.Stop()

	var documents []entity.PlanResponse
	for {
		var tenant entity.PlanResponse
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

func (u *PlanRepo) GetById(id *string) (*entity.PlanResponse, error) {
	doc, err := u.db.Collection("plans").Doc(*id).Get(context.Background())
	if err != nil {
		return nil, err
	}

	var plan entity.PlanResponse
	err = doc.DataTo(&plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (u *PlanRepo) Update(data *entity.PlanResponse) (*entity.PlanResponse, error) {

	_, err := u.db.Collection("plans").Doc(data.ID).Set(context.Background(), data)
	if err != nil {
		return nil, err
	}

	doc, err := u.db.Collection("plans").Doc(data.ID).Get(context.Background())
	if err != nil {
		return nil, err
	}

	var plan entity.PlanResponse
	err = doc.DataTo(&plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (u *PlanRepo) Delete(id *string) error {
	_, err := u.db.Collection("plans").Doc(*id).Delete(context.Background())

	return err
}

func (u *PlanRepo) GetByFilterMany(key string, value *string) ([]entity.PlanResponse, error) {

	var data interface{} = *value
	resultInt, err := strconv.Atoi(*value)
	if err == nil {
		data = resultInt
	}

	resultBool, err := strconv.ParseBool(*value)
	if err == nil {
		data = resultBool
	}

	iter := u.db.Collection("plans").Where(key, "==", data).Documents(context.Background())

	defer iter.Stop()

	var documents []entity.PlanResponse
	for {
		var plan entity.PlanResponse
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&plan)
		if err != nil {
			return nil, err
		}
		plan.ID = doc.Ref.ID

		documents = append(documents, plan)
	}

	if len(documents) == 0 {
		return nil, nil
	}

	return documents, nil
}

func (u *PlanRepo) GetByFilterOne(key string, value *string) (*entity.PlanResponse, error) {
	doc := u.db.Collection("plans").Where(key, "==", *value).Limit(1).Documents(context.Background())
	defer doc.Stop()

	result, err := doc.Next()
	if err != nil {
		return nil, err

	}

	if result == nil {
		return nil, nil

	}

	var plan entity.PlanResponse
	err = result.DataTo(&plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}
