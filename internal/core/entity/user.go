package entity

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

// IUser interface
// Methods that must be implemented by the user
// Create, Get, GetById, Update, Delete, GetByFilterMany, GetByFilterOne, GetByEmail
type IUser interface {
	Create(ctx context.Context, user *AccountUser) (*AccountUser, error)
	Get(ctx context.Context) ([]AccountUser, error)
	GetById(ctx context.Context, id *string) (*AccountUser, error)
	Update(ctx context.Context, data *AccountUser) (*AccountUser, error)
	Delete(ctx context.Context, id *string) error
	GetByFilterMany(ctx context.Context, filter []QueryDBClause) ([]AccountUser, error)
	GetByFilterOne(ctx context.Context, filter []QueryDBClause) (*AccountUser, error)
	GetByEmail(ctx context.Context, email *string) (*AccountUser, error)
}

type AccountRoles struct {
	Key   string `json:"key"   firestore:"key"`
	Name  string `json:"name"  firestore:"name"`
	Value string `json:"value" firestore:"vale"`
}

// AccountUser struct
// User struct with tenant_id and roles
// ID, User, TenantID, Roles
type AccountUser struct {
	ID       string         `json:"id" firestore:"id"`
	TenantID string         `json:"tenant_id" firestore:"tenant_id"`
	Roles    []AccountRoles `json:"roles" firestore:"roles"`
	User
	CreatedAt time.Time `json:"created_at" firestore:"create_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"update_at"`
}

// User struct
// User struct with name, email, avatar_url, provider, first_name, last_name, nick_name, description, user_id, location
type User struct {
	Name        string `json:"name" binding:"required" firestore:"name"`
	Email       string `json:"email" binding:"required" firestore:"email"`
	AvatarURL   string `json:"avatar_url" firestore:"avatar_url"`
	Provider    string `json:"provider" binding:"required" firestore:"provider"`
	FirstName   string `json:"first_name" firestore:"first_name"`
	LastName    string `json:"last_name" firestore:"last_name"`
	NickName    string `json:"nick_name" firestore:"nick_name"`
	Description string `json:"description" firestore:"description"`
	UserID      string `json:"user_id" firestore:"user_id"`
	Location    string `json:"location" firestore:"location"`
}

// NewUser function
// Create a new user
// Return a new AccountUser and error
func NewUser(user *User) (*AccountUser, error) {

	id, _ := uuid.NewV7()

	tenant_id, _ := uuid.NewV7()

	if user == nil || *user == (User{}) {
		return nil, errors.New("account user cannot be empty")
	}

	u := &AccountUser{
		ID:        id.String(),
		TenantID:  tenant_id.String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

	if u == nil || reflect.DeepEqual(*u, User{}) {
		return errors.New("account user cannot be empty")
	}

	if u.Email == "" {
		return fmt.Errorf("invalid email")
	}

	if !utils.IsValidEmail(u.Email) {
		return errors.New("invalid email")
	}

	if u.Provider == "" || (!strings.Contains(strings.ToLower(u.Provider), "local") && !strings.Contains(strings.ToLower(u.Provider), "google")) {
		return errors.New("provider is required")
	}

	return nil
}

func (c *User) IsEmpty(data *User) bool {
	return data == nil || reflect.DeepEqual(*data, User{})
}

func (u *AccountUser) Validate() error {

	if u == nil || reflect.DeepEqual(*u, AccountUser{}) {
		return errors.New("account user cannot be empty")
	}

	_, err := uuid.Parse(u.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	_, err = uuid.Parse(u.TenantID)
	if err != nil {
		return fmt.Errorf("tenant ID is required")
	}

	return u.User.Validate()
}

func (c *AccountUser) IsEmpty(data *AccountUser) bool {
	return data == nil || reflect.DeepEqual(*data, AccountUser{})
}
