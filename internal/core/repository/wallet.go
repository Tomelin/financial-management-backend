package repository

import (
	"context"
	"errors"

	"github.com/Tomelin/financial-management-backend/internal/core/entity"
	"github.com/Tomelin/financial-management-backend/pkg/db"
	"google.golang.org/api/iterator"
)

type WalletRepo struct {
	db db.FirebaseDatabaseInterface
}

// NewWalletRepo
func NewWalletRepo(db db.FirebaseDatabaseInterface) (entity.IWallet, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	return &WalletRepo{db: db}, nil
}

func (w *WalletRepo) Create(ctx context.Context, userId *string, data *entity.WalletResponse) (*entity.WalletResponse, *entity.ModuleError) {

	_, err := w.db.Collection("wallets").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "Create", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	WalletResponse, Merr := w.GetWalletByIdAndUserID(ctx, userId, &data.ID)
	if Merr != nil {
		return nil, Merr
	}

	return WalletResponse, nil
}

func (w *WalletRepo) Get(ctx context.Context, userId *string) ([]entity.WalletResponse, *entity.ModuleError) {
	iter := w.db.Collection("wallets").Where("owner_id", "==", *userId).Documents(context.Background())

	defer iter.Stop()
	var wallets []entity.WalletResponse
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, entity.Error(err.Error(), "wallet", "Get", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
		}

		var wallet entity.WalletResponse
		err = doc.DataTo(&wallet)
		if err != nil {
			return nil, entity.Error(err.Error(), "wallet", "Get", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
		}
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func (w *WalletRepo) GetWalletByIdAndUserID(ctx context.Context, userId *string, id *string) (*entity.WalletResponse, *entity.ModuleError) {
	doc := w.db.Collection("wallets").Where("id", "==", *id).Where("owner_id", "==", *userId).Documents(context.Background())
	defer doc.Stop()

	data, _ := doc.Next()
	if data == nil {
		return &entity.WalletResponse{}, nil
	}

	var wallet entity.WalletResponse
	err := data.DataTo(&wallet)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "GetById", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	return &wallet, nil
}

func (w *WalletRepo) Update(ctx context.Context, userId *string, data *entity.WalletResponse) (*entity.WalletResponse, *entity.ModuleError) {
	// Primeiro, obtenha a referência ao documento
	docRef := w.db.Collection("wallets").Doc(data.ID)

	// Verifique se o documento pertence ao usuário
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	var wallet entity.WalletResponse
	err = doc.DataTo(&wallet)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	if wallet.OwnerID != *userId {
		return nil, entity.Error("unauthorized: wallet does not belong to the user", "wallet", "Update", entity.ApplicationLayerRepository, entity.ResponseCodeUnauthorized)
	}

	// Atualize o documento
	_, err = docRef.Set(ctx, data)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	// Obtenha o documento atualizado
	doc, err = docRef.Get(ctx)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	var updatedWallet entity.WalletResponse
	err = doc.DataTo(&updatedWallet)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "Update", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	return &updatedWallet, nil
}

func (w *WalletRepo) Delete(ctx context.Context, userId *string, id *string) *entity.ModuleError {

	_, err := w.db.Collection("wallets").Doc(*id).Delete(context.Background())
	if err != nil {
		return entity.Error(err.Error(), "wallet", "Delete", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	return nil
}

func (w *WalletRepo) GetByFilterMany(ctx context.Context, userId *string, filter []entity.QueryDB) ([]entity.WalletResponse, *entity.ModuleError) {
	query := w.db.Collection("wallets").Where("owner_id", "==", *userId)
	for _, f := range filter {
		condition := checkFirebaseCondition(&f.Condition)
		if f.Key != "" && f.Value != "" && condition != "" {
			query = query.Where(f.Key, condition, f.Value)
		}
	}

	iter := query.Documents(context.Background())
	defer iter.Stop()
	var wallets []entity.WalletResponse
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, entity.Error(err.Error(), "wallet", "GetByFilterMany", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
		}

		var wallet entity.WalletResponse
		err = doc.DataTo(&wallet)
		if err != nil {
			return nil, entity.Error(err.Error(), "wallet", "GetByFilterMany", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
		}
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func (w *WalletRepo) GetByFilterOne(ctx context.Context, userId *string, filter []entity.QueryDB) (*entity.WalletResponse, *entity.ModuleError) {

	query := w.db.Collection("wallets").Where("owner_id", "==", *userId)
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
		return nil, entity.Error(err.Error(), "wallet", "GetByFilterOne", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	if result == nil {
		return nil, nil
	}

	var wallet entity.WalletResponse
	err = result.DataTo(&wallet)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "GetByFilterOne", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	return &wallet, nil
}

func (w *WalletRepo) GetByID(ctx context.Context, walletId *string) (*entity.WalletResponse, *entity.ModuleError) {
	doc := w.db.Collection("wallets").Where("id", "==", *walletId).Documents(context.Background())
	defer doc.Stop()

	data, _ := doc.Next()
	if data == nil {
		return &entity.WalletResponse{}, nil
	}

	var wallet entity.WalletResponse
	err := data.DataTo(&wallet)
	if err != nil {
		return nil, entity.Error(err.Error(), "wallet", "GetById", entity.ApplicationLayerRepository, entity.ResponseCodeInternalServer)
	}

	return &wallet, nil
}

func (w *WalletRepo) UpdateBalance(ctx context.Context, walletID *string, balance *float64) (*entity.WalletResponse, *entity.ModuleError) {
	return nil, nil
}
