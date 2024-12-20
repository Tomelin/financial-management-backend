package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
)

type UserTestSuite struct {
	suite.Suite
	user        *entity.User
	accountUser *entity.AccountUser
}

func (s *UserTestSuite) SetupTest() {
	s.user = &entity.User{
		Name:      "Teste",
		Email:     "teste@example.com",
		Provider:  "google",
		FirstName: "Nome",
		LastName:  "Sobrenome",
	}

	id := uuid.New().String()
	s.accountUser = &entity.AccountUser{
		ID:       id,
		TenantID: id,
		Roles:    []entity.AccountRoles{},
		User: entity.User{
			Name:      "Teste",
			Email:     "teste@example.com",
			Provider:  "google",
			FirstName: "Nome",
			LastName:  "Sobrenome",
		},
	}
}

func (s *UserTestSuite) TearDownTest() {
	s.user = nil
	s.accountUser = nil
}

func (s *UserTestSuite) TestNewUser() {

	accountUser, err := entity.NewUser(s.user)

	s.NoError(err)
	s.NotNil(accountUser)
	s.NotEmpty(accountUser.ID)
	s.NotEmpty(accountUser.TenantID)
	s.Equal(s.user.Name, accountUser.Name)
	s.Equal(s.user.Email, accountUser.Email)
}

func (s *UserTestSuite) TestNewUser_Error_UserNil() {

	s.user = nil
	accountUser, err := entity.NewUser(s.user)
	s.Error(err)
	s.Nil(accountUser)
	s.Equal("account user cannot be empty", err.Error())

}

func (s *UserTestSuite) TestNewUser_Error_InvalidEmail() {

	s.user.Email = ""
	accountUser, err := entity.NewUser(s.user)
	s.Error(err)
	s.Nil(accountUser)
	s.Equal("invalid email", err.Error())

	s.user.Email = "user"
	accountUser, err = entity.NewUser(s.user)
	s.Error(err)
	s.Nil(accountUser)
	s.Equal("invalid email", err.Error())

	s.user.Email = "teste@domain"
	accountUser, err = entity.NewUser(s.user)
	s.Error(err)
	s.Nil(accountUser)
	s.Equal("invalid email", err.Error())
}

func (s *UserTestSuite) TestNewUser_Error_InvalidProvider() {
	s.user.Provider = ""
	// Teste com User nulo
	accountUser, err := entity.NewUser(s.user)
	s.Error(err)
	s.Nil(accountUser)
	s.Equal("provider is required", err.Error())

	// Teste com User nulo
	s.user.Provider = "teste"
	accountUser, err = entity.NewUser(s.user)
	s.Error(err)
	s.Nil(accountUser)
	s.Equal("provider is required", err.Error())
}

func (s *UserTestSuite) TestNewUser_Error_UserValidate() {

	s.user = nil
	err := s.user.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("account user cannot be empty", err.Error())

	s.user = &entity.User{}
	err = s.user.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("account user cannot be empty", err.Error())

	emptyUser := entity.User{}
	err = emptyUser.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("account user cannot be empty", err.Error())
}

func (s *UserTestSuite) TestNewUser_Error_UserValidate_Email() {
	s.user.Email = ""
	err := s.user.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("invalid email", err.Error())

	s.user.Email = "username"
	err = s.user.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("invalid email", err.Error())

	s.user.Email = "username@domain"
	err = s.user.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("invalid email", err.Error())

	s.user.Email = "@domain"
	err = s.user.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("invalid email", err.Error())

	s.user.Email = "domain.com.br"
	err = s.user.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("invalid email", err.Error())
}

func (s *UserTestSuite) TestNewUser_Error_UserValidate_Provider() {

	s.user.Provider = ""
	err := s.user.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("provider is required", err.Error())

	s.user.Provider = "googl"
	err = s.user.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("provider is required", err.Error())
}

func (s *UserTestSuite) TestNewUser_Error_IsEmpty() {

	s.user = nil
	b := s.user.IsEmpty(s.user)
	assert.True(s.T(), b)

	s.user = &entity.User{}
	b = s.user.IsEmpty(s.user)
	assert.True(s.T(), b)

	emptyUser := entity.User{}
	b = s.user.IsEmpty(&emptyUser)
	assert.True(s.T(), b)
}

func (s *UserTestSuite) TestResponseUser_Error_UserValidate() {
	var user *entity.AccountUser
	s.accountUser = user

	// Teste com User nulo
	err := s.accountUser.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("account user cannot be empty", err.Error())

	s.accountUser = &entity.AccountUser{}
	err = s.accountUser.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("account user cannot be empty", err.Error())

	emptyUser := entity.AccountUser{}
	err = emptyUser.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("account user cannot be empty", err.Error())
}

func (s *UserTestSuite) TestResponseUser_Error_ParseUserID() {

	s.accountUser.ID = ""
	err := s.accountUser.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("invalid user ID", err.Error())

	s.accountUser.ID = "123"
	err = s.accountUser.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("invalid user ID", err.Error())
}

func (s *UserTestSuite) TestResponseUser_Error_ParseTenantID() {

	s.accountUser.TenantID = ""
	err := s.accountUser.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("tenant ID is required", err.Error())

	s.accountUser.TenantID = "123"
	err = s.accountUser.Validate()
	s.NotNil(err)
	s.Error(err)
	s.Equal("tenant ID is required", err.Error())
}

func (s *UserTestSuite) TestResponseUser_Error_IsEmpty() {

	var user *entity.AccountUser
	s.accountUser = user

	b := s.accountUser.IsEmpty(s.accountUser)
	s.True(b)

	user = &entity.AccountUser{}
	b = s.accountUser.IsEmpty(s.accountUser)
	s.True(b)

	emptyUser := entity.AccountUser{}
	b = s.accountUser.IsEmpty(&emptyUser)
	s.True(b)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
