package repository

import (
	"context"
	"errors"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/db"
	"google.golang.org/api/iterator"
)

type UserRepo struct {
	db db.FirebaseDatabaseInterface
}

func NewUserRepository(db db.FirebaseDatabaseInterface) (*UserRepo, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}
	return &UserRepo{db: db}, nil
}

func (u *UserRepo) Create(user *entity.UserResponse) (*entity.UserResponse, error) {
	ctx := context.Background()

	_, err := u.db.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return nil, err
	}

	doc, err := u.db.Collection("users").Doc(user.ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var userResponse entity.UserResponse
	err = doc.DataTo(&userResponse)
	if err != nil {
		return nil, err
	}

	return &userResponse, nil
}

func (u *UserRepo) Get() ([]entity.UserResponse, error) {
	return nil, nil
}

func (u *UserRepo) GetById(id *string) (*entity.UserResponse, error) {
	return nil, nil
}

func (u *UserRepo) Update(data *entity.UserResponse) (*entity.UserResponse, error) {
	return nil, nil
}

func (u *UserRepo) Delete(id *string) error {
	return nil
}

func (u *UserRepo) GetByFilterMany(key string, value *string) ([]entity.UserResponse, error) {
	return nil, nil
}

func (u *UserRepo) GetByFilterOne(key string, value *string) (*entity.UserResponse, error) {

	iter := u.db.Documents(context.Background(), "users")

	defer iter.Stop()

	for {
		var user entity.UserResponse
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&user)
		if err != nil {
			return nil, err
		}

		if user.Email == *value {
			user.ID = doc.Ref.ID
			return &user, nil
		}
	}

	return nil, nil
}
