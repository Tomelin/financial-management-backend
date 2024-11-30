package service

import (
	"errors"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
)

type IUserService interface {
	entity.IUser
	NewUser(entity.User) (*entity.UserResponse, error)
}
type UserSvc struct {
	repo entity.IUser
}

func NewUserService(u entity.IUser) (IUserService, error) {

	if u == nil {
		return nil, errors.New("repository is required")
	}

	return &UserSvc{repo: u}, nil
}

func (u *UserSvc) Create(user *entity.UserResponse) (*entity.UserResponse, error) {

	if user == nil {
		return nil, errors.New("user is required")
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	res, err := u.GetByFilterOne("email", &user.Email)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}

	if res.IsEmpty(res) {
		response, err := u.repo.Create(user)
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	return nil, errors.New("user already exists")
}

func (u *UserSvc) Get() ([]entity.UserResponse, error) {
	return nil, nil
}

func (u *UserSvc) GetById(id *string) (*entity.UserResponse, error) {
	return nil, nil
}

func (u *UserSvc) Update(data *entity.UserResponse) (*entity.UserResponse, error) {
	return nil, nil
}

func (u *UserSvc) Delete(id *string) error {
	return nil
}

func (u *UserSvc) GetByFilterMany(key string, value *string) ([]entity.UserResponse, error) {
	return nil, nil
}

func (u *UserSvc) GetByFilterOne(key string, value *string) (*entity.UserResponse, error) {
	if key == "" {
		return nil, errors.New("key cannot be empty")
	}

	if value == nil || *value == "" {
		return nil, errors.New("value cannot be empty")
	}

	return u.repo.GetByFilterOne(key, value)
}

func (u *UserSvc) NewUser(entity.User) (*entity.UserResponse, error) {
	return nil, nil
}
