package entity

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type IUser interface {
	Create(*UserResponse) (*UserResponse, error)
	Get() ([]UserResponse, error)
	GetById(id *string) (*UserResponse, error)
	Update(data *UserResponse) (*UserResponse, error)
	Delete(id *string) error
	GetByFilterMany(key string, value *string) ([]UserResponse, error)
	GetByFilterOne(key string, value *string) (*UserResponse, error)
}

type User struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	AvatarURL   string `json:"avatar_url"`
	Provider    string `json:"provider" binding:"required"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	NickName    string `json:"nick_name"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	Location    string `json:"location"`
}

type UserResponse struct {
	ID string `json:"id"`
	User
}

func NewUser(user *User) (*UserResponse, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	u := &UserResponse{
		ID: id.String(),
		User: User{
			Name:        user.Name,
			Email:       user.Email,
			AvatarURL:   user.AvatarURL,
			Provider:    user.Provider,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			NickName:    user.NickName,
			Description: user.Description,
			UserID:      user.UserID,
			Location:    user.Location,
		},
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Validate() error {
	// if u.Name == "" {
	// 	return fmt.Errorf("name is required")
	// }

	if u.Email == "" {
		return fmt.Errorf("email is required")
	}

	if !utils.IsValidEmail(u.Email) {
		return fmt.Errorf("email is invalid")
	}

	if u.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	return nil
}

func (c *User) IsEmpty(data *User) bool {

	return data == nil || reflect.DeepEqual(*data, User{})
}

func (u *UserResponse) Validate() error {
	_, err := uuid.Parse(u.ID)
	if err != nil {
		return fmt.Errorf("id is required")
	}

	return u.User.Validate()
}

func (c *UserResponse) IsEmpty(data *UserResponse) bool {

	return data == nil || reflect.DeepEqual(*data, UserResponse{})
}
