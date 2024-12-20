package repository

import (
	"context"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/db"
	"github.com/synera-br/financial-management/src/backend/pkg/observability"
	"google.golang.org/api/iterator"
)

type TransactionCategoryRepo struct {
	db    db.FirebaseDatabaseInterface
	trace *observability.Tracer
}

func NewTransactionCategoryRepo(trace *observability.Tracer, db db.FirebaseDatabaseInterface) (entity.ITransactionCategoryRepository, *entity.ModuleError) {

	if db == nil {
		return nil, entity.Error("database is required", "transactionCategory", "inicialization", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	return &TransactionCategoryRepo{
		db:    db,
		trace: trace,
	}, nil
}

func (c *TransactionCategoryRepo) Create(ctx context.Context, category *entity.TransactionCategory) (*entity.TransactionCategory, *entity.ModuleError) {
	ctx, span := c.trace.Trace.Start(ctx, "TransactionCategoryRepo.Create")
	defer span.End()

	_, err := c.db.Collection("transaction_categories").Doc(category.ID).Set(ctx, category)
	if err != nil {
		return nil, entity.Error(err.Error(), "transactionCategory", "Create", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	doc, err := c.db.Collection("transaction_categories").Doc(category.ID).Get(ctx)
	if err != nil {
		return nil, entity.Error(err.Error(), "transactionCategory", "Create", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	var cat entity.TransactionCategory
	err = doc.DataTo(&cat)
	if err != nil {
		return nil, entity.Error(err.Error(), "transactionCategory", "Create", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	return &cat, nil
}

func (c *TransactionCategoryRepo) Get(ctx context.Context, walletID *string) ([]entity.TransactionCategory, *entity.ModuleError) {
	return nil, nil
}

func (c *TransactionCategoryRepo) GetById(ctx context.Context, id *string) (*entity.TransactionCategory, *entity.ModuleError) {
	doc := c.db.Collection("transaction_categories").Where("id", "==", *id).Documents(context.Background())
	defer doc.Stop()

	data, _ := doc.Next()
	if data == nil {
		return &entity.TransactionCategory{}, nil
	}

	var category entity.TransactionCategory
	err := data.DataTo(&category)
	if err != nil {
		return nil, entity.Error(err.Error(), "category", "GetById", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	return &category, nil
}

func (c *TransactionCategoryRepo) Update(ctx context.Context, category *entity.TransactionCategory) (*entity.TransactionCategory, *entity.ModuleError) {
	return nil, nil
}

func (c *TransactionCategoryRepo) Delete(ctx context.Context, id *string) *entity.ModuleError {
	return nil
}

func (c *TransactionCategoryRepo) GetByFilterMany(ctx context.Context, filter []entity.QueryDBClause) ([]entity.TransactionCategory, *entity.ModuleError) {
	ctx, span := c.trace.Trace.Start(ctx, "TransactionCategoryRepo.GetByFilterMany")
	defer span.End()

	hasQuery := false
	hasQueryOr := false
	query := c.db.Collection("transaction_categories").Query
	queryOr := c.db.Collection("transaction_categories").Query
	for _, f := range filter {
		if f.Clause == "" || f.Clause == entity.QueryClauseAnd {
			hasQuery = true
			for _, q := range f.Queries {
				condition := checkFirebaseCondition(&q.Condition)
				if q.Key != "" && q.Value != "" && condition != "" {
					query = query.Where(q.Key, condition, q.Value)
				}
			}
		}
		if f.Clause == entity.QueryClauseOr {
			hasQueryOr = true
			for _, q := range f.Queries {
				condition := checkFirebaseCondition(&q.Condition)
				if q.Key != "" && q.Value != "" && condition != "" {
					queryOr = queryOr.Where(q.Key, condition, q.Value)
				}
			}
		}

	}

	var categories []entity.TransactionCategory
	if hasQuery {
		iter := query.Documents(ctx)
		defer iter.Stop()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, entity.Error(err.Error(), "transactionCategory", "GetByFilterMany", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
			}

			var category entity.TransactionCategory
			err = doc.DataTo(&category)
			if err != nil {
				return nil, entity.Error(err.Error(), "transactionCategory", "GetByFilterMany", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
			}
			categories = append(categories, category)
		}
	}

	if hasQueryOr {
		iter := queryOr.Documents(ctx)
		defer iter.Stop()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, entity.Error(err.Error(), "transactionCategory", "GetByFilterMany", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
			}

			var category entity.TransactionCategory
			err = doc.DataTo(&category)
			if err != nil {
				return nil, entity.Error(err.Error(), "transactionCategory", "GetByFilterMany", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
			}
			categories = append(categories, category)
		}
	}

	return categories, nil
}

func (c *TransactionCategoryRepo) GetByFilterOne(ctx context.Context, filter []entity.QueryDB) (*entity.TransactionCategory, *entity.ModuleError) {
	ctx, span := c.trace.Trace.Start(ctx, "TransactionCategoryRepo.GetByFilterOne")
	defer span.End()

	query := c.db.Collection("transaction_categories").Query

	for _, f := range filter {
		condition := checkFirebaseCondition(&f.Condition)
		if f.Key != "" && f.Value != "" && condition != "" {
			query = query.Where(f.Key, condition, f.Value)
		}
	}

	doc := query.Limit(1).Documents(ctx)
	result, err := doc.Next()
	if err != nil {
		if err == iterator.Done {
			return nil, nil
		}
		return nil, entity.Error(err.Error(), "transactionCategory", "GetByFilterOne", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	if result == nil {
		return nil, nil
	}

	var tenant entity.TransactionCategory
	err = result.DataTo(&tenant)
	if err != nil {
		return nil, entity.Error(err.Error(), "transactionCategory", "GetByFilterOne", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	return &tenant, nil
}

func (c *TransactionCategoryRepo) filters(ctx context.Context, filter []entity.QueryDB) (*entity.TransactionCategory, *entity.ModuleError) {

	return nil, nil
}
