package service

import (
	"context"
	"strconv"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/observability"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type TransactionCategorySvc struct {
	repo   entity.ITransactionCategoryRepository
	tenant ITenantService
	wallet entity.IWallet
	Trace  *observability.Tracer
	user   entity.IUser
}

func NewTransactionCategorySvc(
	trace *observability.Tracer,
	repo *entity.ITransactionCategoryRepository,
	tenant *ITenantService,
	wallet *entity.IWallet,
	user *entity.IUser,
) (entity.ITransactionCategory, *entity.ModuleError) {

	if repo == nil {
		return nil, entity.Error("repo is required", "transactionCategory", "inicialization", entity.ApplicationLayerService, entity.ResponseCodeInternalServer)
	}

	if tenant == nil {
		return nil, entity.Error("tenant is required", "transactionCategory", "inicialization", entity.ApplicationLayerService, entity.ResponseCodeInternalServer)
	}

	if wallet == nil {
		return nil, entity.Error("wallet is required", "transactionCategory", "inicialization", entity.ApplicationLayerService, entity.ResponseCodeInternalServer)
	}

	return &TransactionCategorySvc{
		repo:   *repo,
		tenant: *tenant,
		wallet: *wallet,
		Trace:  trace,
		user:   *user,
	}, nil
}

func (c *TransactionCategorySvc) Create(ctx context.Context, email *string, category *entity.TransactionCategory) (*entity.TransactionCategory, *entity.ModuleError) {
	ctx, span := c.Trace.Trace.Start(ctx, "TransactionCategorySvc.Create")
	defer span.End()

	// Validate if the wallet is empty
	if category == nil || category.IsEmpty(category) {
		return nil, entity.Error("category required", "transactionCategory", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if mErr := category.Validate(); mErr != nil {
		return nil, mErr
	}

	walletDefault, _ := strconv.ParseBool(category.Default)

	if walletDefault && (category.WalletID != "" || category.TenantID != "") {
		return nil, entity.Error("wallet default cannot be wallet id or tenant id", "transactionCategory", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if !walletDefault && (category.WalletID == "" || category.TenantID == "") {
		return nil, entity.Error("wallet id and tenant id are required", "transactionCategory", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if !walletDefault && category.WalletID != "" && category.TenantID != "" {
		// Validate if the tenant ID is valid
		tenantId := category.TenantID
		if err := utils.ValidateUUID(&tenantId); err != nil {
			return nil, entity.Error(err.Error(), "transactionCategory", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
		}

		// Validate if the tenant exists
		tenant, err := c.tenant.GetById(&tenantId)
		if err != nil || tenant == nil || tenant.Name == "" {
			return nil, entity.Error(err.Error(), "transactionCategory", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
		}

		// Validate if the tenant ID is valid
		walletId := category.WalletID
		if err := utils.ValidateUUID(&walletId); err != nil {
			return nil, entity.Error(err.Error(), "transactionCategory", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
		}

		// Validate if the wallet exists
		wallet, mErr := c.wallet.GetByID(ctx, &walletId)
		if mErr != nil || wallet == nil || wallet.Name == "" {
			return nil, mErr
		}
	}

	user, err := c.user.GetByEmail(ctx, email)
	if err != nil {
		return nil, entity.Error(err.Error(), "transactionCategory", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if user == nil || user.Email == "" {
		return nil, entity.Error("user not found", "transactionCategory", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	// Validate if the category already exists
	data, _ := c.GetByFilterMany(ctx, email, []entity.QueryDB{
		{
			Key:       "name",
			Condition: string(entity.QueryFirebaseEqual),
			Value:     category.Name,
		},
	})

	if len(data) > 0 {
		return nil, entity.Error("category already exists", "transactionCategory", "Create", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	return c.repo.Create(ctx, category)
}

func (c *TransactionCategorySvc) Get(ctx context.Context, email *string, walletID *string) ([]entity.TransactionCategory, *entity.ModuleError) {
	ctx, span := c.Trace.Trace.Start(ctx, "TransactionCategorySvc.Get")
	span.SetName("TransactionCategorySvc")
	defer span.End()

	data, mErr := c.GetByFilterMany(ctx, email, []entity.QueryDB{
		{
			Key:       "default",
			Condition: string(entity.QueryFirebaseEqual),
			Value:     "true",
		},
	})
	if mErr != nil {
		return nil, mErr
	}

	if walletID != nil && *walletID != "" {
		if err := utils.ValidateUUID(walletID); err != nil {
			return nil, entity.Error(err.Error(), "transactionCategory", "Get", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
		}

		// Validate if the wallet exists
		wallet, mErr := c.wallet.GetByID(ctx, walletID)
		if mErr != nil || wallet == nil || wallet.Name == "" {
			return nil, mErr
		}

		dataTemp, mErr := c.GetByFilterMany(ctx, email, []entity.QueryDB{
			{
				Key:       "default",
				Condition: string(entity.QueryFirebaseEqual),
				Value:     "false",
			},
			{
				Key:       "wallet_id",
				Condition: string(entity.QueryFirebaseEqual),
				Value:     *walletID,
			},
		})
		if mErr != nil {
			return nil, mErr
		}

		data = append(data, dataTemp...)
	}

	return data, nil

}

func (c *TransactionCategorySvc) GetById(ctx context.Context, email *string, id *string) (*entity.TransactionCategory, *entity.ModuleError) {
	ctx, span := c.Trace.Trace.Start(ctx, "TransactionCategorySvc.GetById")
	span.SetName("TransactionCategorySvc")
	defer span.End()

	if id == nil || *id == "" {
		return nil, entity.Error("id cannot be empty", "transactionCategory", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if err := utils.ValidateUUID(id); err != nil {
		return nil, entity.Error(err.Error(), "transactionCategory", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	result, err := c.user.GetByEmail(ctx, email)
	if err != nil {
		return nil, entity.Error(err.Error(), "transactionCategory", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if result == nil || result.Email == "" {
		return nil, entity.Error("user not found", "transactionCategory", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if err != nil {
		return nil, entity.Error(err.Error(), "transactionCategory", "GetById", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	data, mErr := c.repo.GetById(ctx, id)
	if mErr != nil {
		return nil, mErr
	}

	if data == nil || data.Name == "" {
		return nil, entity.Error("category not found", "transactionCategory", "GetById", entity.ApplicationLayerService, entity.ResponseCodeNotFound)
	}

	return data, nil
}

func (c *TransactionCategorySvc) Update(ctx context.Context, email *string, category *entity.TransactionCategory) (*entity.TransactionCategory, *entity.ModuleError) {
	return nil, nil
}

func (c *TransactionCategorySvc) Delete(ctx context.Context, email *string, id *string) *entity.ModuleError {

	if email == nil || *email == "" {
		return entity.Error("email cannot be empty", "transactionCategory", "Delete", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if id == nil || *id == "" {
		return entity.Error("id cannot be empty", "transactionCategory", "Delete", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if err := utils.ValidateUUID(id); err != nil {
		return entity.Error(err.Error(), "transactionCategory", "Delete", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if ok := utils.IsValidEmail(*email); !ok {
		return entity.Error("email is invalid", "transactionCategory", "Delete", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	cat, mErr := c.GetById(ctx, email, id)
	if mErr != nil {
		return mErr
	}

	if cat == nil || cat.Name == "" {
		return entity.Error("category not found", "transactionCategory", "Delete", entity.ApplicationLayerService, entity.ResponseCodeNotFound)
	}

	if cat.Default == "true" {
		return entity.Error("category default cannot be deleted", "transactionCategory", "Delete", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	user, err := c.user.GetByEmail(ctx, email)
	if err != nil || user == nil || user.Email == "" {
		message := "user not found"
		if err != nil {
			message = err.Error()
		}
		return entity.Error(message, "transactionCategory", "Delete", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	tenant, err := c.tenant.GetById(&cat.TenantID)
	if err != nil || tenant == nil || tenant.Name == "" {
		message := "tenant not found"
		if err != nil {
			message = err.Error()
		}
		return entity.Error(message, "transactionCategory", "Delete", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	// _, mErr = c.wallet.GetByFilterOne(ctx, []entity.QueryDB{})

	mErr = c.repo.Delete(ctx, id)
	if mErr != nil {
		return mErr
	}

	return nil
}

func (c *TransactionCategorySvc) GetByFilterMany(ctx context.Context, email *string, filter []entity.QueryDB) ([]entity.TransactionCategory, *entity.ModuleError) {
	ctx, span := c.Trace.Trace.Start(ctx, "TransactionCategorySvc.GetByFilterMany")
	defer span.End()

	if len(filter) == 0 {
		return nil, entity.Error("filter cannot be empty", "transactionCategory", "GetByFilterMany", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if filter[0].Key == "" || filter[0].Value == "" {
		return nil, entity.Error("key and value cannot be empty", "transactionCategory", "GetByFilterMany", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if email == nil {
		return nil, entity.Error("user cannot be empty", "transactionCategory", "GetByFilterMany", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	user, err := c.user.GetByEmail(ctx, email)
	if err != nil {
		return nil, entity.Error(err.Error(), "transactionCategory", "GetByFilterMany", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if user == nil || user.Email == "" {
		return nil, entity.Error("user not found", "transactionCategory", "GetByFilterMany", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	tenant, err := c.tenant.GetByFilterOne(ctx, []entity.QueryDB{{
		Key:       "owner_id",
		Condition: string(entity.QueryFirebaseEqual),
		Value:     user.ID,
	}})

	if err != nil {
		return nil, entity.Error(err.Error(), "transactionCategory", "GetByFilterMany", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if tenant == nil || tenant.Name == "" {
		return nil, entity.Error("tenant not found", "transactionCategory", "GetByFilterMany", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	filters2 := []entity.QueryDBClause{
		{
			Clause: entity.QueryClauseAnd,
			Queries: append(filter,
				[]entity.QueryDB{
					{
						Key:       "default",
						Condition: string(entity.QueryFirebaseEqual),
						Value:     "true",
					},
					{
						Key:       filter[0].Key,
						Condition: string(entity.QueryFirebaseEqual),
						Value:     filter[0].Value,
					},
				}...),
		},
		{
			Clause: entity.QueryClauseOr,
			Queries: append(filter, []entity.QueryDB{
				{
					Key:       "default",
					Condition: string(entity.QueryFirebaseEqual),
					Value:     "false",
				},
				{
					Key:       "tenant_id",
					Condition: string(entity.QueryFirebaseEqual),
					Value:     tenant.ID,
				},
			}...),
		},
	}

	data, mErr := c.repo.GetByFilterMany(ctx, filters2)
	if mErr != nil {
		return nil, mErr
	}

	return data, nil
}

func (c *TransactionCategorySvc) GetByFilterOne(ctx context.Context, email *string, filter []entity.QueryDB) (*entity.TransactionCategory, *entity.ModuleError) {
	ctx, span := c.Trace.Trace.Start(ctx, "TransactionCategorySvc.GetByFilterOne")
	defer span.End()

	if len(filter) == 0 {
		return nil, entity.Error("filter cannot be empty", "transactionCategory", "GetByFilterOne", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if filter[0].Key == "" || filter[0].Value == "" {
		return nil, entity.Error("key and value cannot be empty", "transactionCategory", "GetByFilterOne", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	if email == nil {
		return nil, entity.Error("user cannot be empty", "transactionCategory", "GetByFilterOne", entity.ApplicationLayerService, entity.ResponseCodeBadRequest)
	}

	data, mErr := c.repo.GetByFilterOne(ctx, filter)
	if mErr != nil {
		return nil, mErr
	}

	return data, nil
}
