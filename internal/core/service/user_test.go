package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/internal/core/service"
	"github.com/synera-br/financial-management/src/backend/tests/coremocks"
)

type UserServiceTestSuite struct {
	suite.Suite
	mockUserEntity        *entity.User
	mockUsersEntity       []entity.User
	mockAccountUserEntity *entity.AccountUser
	mockTenantEntity      *entity.TenantResponse
	mockTenantsEntity     []entity.TenantResponse
	mockCoreTenant        *coremocks.ITenantService
	mockCoreUser          *coremocks.IUser
	mockFilterTenant      []entity.QueryDBClause
	mockFilterUser        []entity.QueryDBClause
	ctx                   context.Context
}

func (s *UserServiceTestSuite) SetupTest() {
	s.mockUserEntity = &entity.User{
		Name:      "Teste",
		Email:     "user@domain.com",
		Provider:  "google",
		FirstName: "Nome",
		LastName:  "Sobrenome",
	}

	id := uuid.New().String()
	s.mockAccountUserEntity = &entity.AccountUser{
		ID:       id,
		TenantID: id,
		Roles:    []entity.AccountRoles{},
		User:     *s.mockUserEntity,
	}

	s.mockTenantEntity = &entity.TenantResponse{
		ID:      s.mockAccountUserEntity.TenantID,
		Name:    s.mockUserEntity.Email,
		Alias:   "",
		OwnerID: s.mockAccountUserEntity.TenantID,
		Users:   []string{s.mockUserEntity.Email},
	}

	s.mockTenantsEntity = append([]entity.TenantResponse{}, *s.mockTenantEntity)
	s.mockUsersEntity = []entity.User{*s.mockUserEntity}

	s.mockCoreTenant = new(coremocks.ITenantService)
	s.mockCoreUser = new(coremocks.IUser)
	s.ctx = context.Background()

	s.mockFilterUser = []entity.QueryDBClause{
		{
			Clause: entity.QueryClauseAnd,
			Queries: []entity.QueryDB{
				{
					Key:       "email",
					Value:     s.mockUserEntity.Email,
					Condition: string(entity.QueryFirebaseEqual),
				},
			},
		},
	}

	s.mockFilterTenant = []entity.QueryDBClause{
		{
			Clause: entity.QueryClauseAnd,
			Queries: []entity.QueryDB{
				{
					Key:       "name",
					Value:     s.mockUserEntity.Email,
					Condition: string(entity.QueryFirebaseEqual),
				},
			},
		},
	}
}

func (s *UserServiceTestSuite) TearDownTest() {
	s.mockUserEntity = nil
	s.mockAccountUserEntity = nil
	s.mockTenantEntity = nil
	s.mockCoreTenant = nil
	s.mockCoreUser = nil
	s.ctx = context.Background()
}

func (s *UserServiceTestSuite) TestUserService_Error_NewUserService() {

	s.mockCoreUser = new(coremocks.IUser)

	svc, err := service.NewUserService(s.mockCoreUser, nil, nil)
	s.Require().Error(err)
	s.Require().Nil(svc)
	s.Require().Equal("error creating user service", err.Error())

	s.Assert().Error(err)
	s.mockCoreUser.AssertExpectations(s.T())

	svc, err = service.NewUserService(nil, s.mockCoreTenant, nil)
	s.Require().Error(err)
	s.Require().Nil(svc)
	s.Require().Equal("error creating user service", err.Error())

	s.Assert().Error(err)
	s.mockCoreTenant.AssertExpectations(s.T())
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

var (
	mockUserNil      = &entity.AccountUser{}
	mockUserEmpty    = entity.AccountUser{}
	mockUserNInvalid = &entity.AccountUser{
		User: entity.User{
			Name:        "Teste",
			Email:       "user@domain.com",
			Provider:    "google",
			FirstName:   "Nome",
			LastName:    "Sobrenome",
			NickName:    "nick",
			Description: "description",
			AvatarURL:   "avatar",
			UserID:      "user_id",
			Location:    "pt-BR",
		},
	}
)

func (s *UserServiceTestSuite) TestUserService_Error_CreateUser() {

	s.mockCoreUser.On("GetByFilterOne", s.ctx, s.mockFilterUser).Return(s.mockUserEntity, nil).Once()
	s.mockCoreTenant.On("GetByFilterMany", s.ctx, s.mockFilterTenant[0].Queries).Return(s.mockTenantsEntity, nil).Once()

	userSvc, err := service.NewUserService(s.mockCoreUser, s.mockCoreTenant, nil)
	s.Require().NoError(err)
	s.Require().NotNil(userSvc)

	// // Test with nil user
	result, err := userSvc.Create(s.ctx, mockUserNil)
	s.Require().Error(err)
	s.Require().Equal("user is required", err.Error())
	s.Require().Nil(result)

	// // Test with empty user
	result, err = userSvc.Create(s.ctx, &mockUserEmpty)
	s.Require().Error(err)
	s.Require().Equal("user is required", err.Error())
	s.Require().Nil(result)

	// // Test with invalid user
	result, err = userSvc.Create(s.ctx, mockUserNInvalid)
	s.Require().Error(err)
	s.Require().Equal("invalid user ID", err.Error())
	s.Require().Nil(result)

	// result, err = userSvc.Create(ctx, mockAccountUser)
	// s.Require().Error(err)
	// s.Require().Equal("user already exists", err.Error())
	// s.Require().Nil(result)
}

var (
	ctx      = context.Background()
	mockUser = &entity.User{
		Name:        "Teste",
		Email:       "user@domain.com",
		Provider:    "google",
		FirstName:   "Nome",
		LastName:    "Sobrenome",
		AvatarURL:   "avatar",
		UserID:      uuid.New().String(),
		Location:    "pt-BR",
		NickName:    "nick",
		Description: "description",
	}
	mockAccountUser = &entity.AccountUser{
		ID:       uuid.New().String(),
		TenantID: uuid.New().String(),
		Roles:    []entity.AccountRoles{},
		User:     *mockUser,
	}
	mockTenantData = &entity.TenantResponse{
		ID:      mockAccountUser.TenantID,
		Name:    mockUser.Email,
		Alias:   "",
		OwnerID: mockAccountUser.TenantID,
		Users:   []string{mockUser.Email},
	}
	mockTenants = []entity.TenantResponse{*mockTenantData}
)

// func TestUserService_Error_CreateUser(t *testing.T) {
// 	// Create a new mock client
// 	mockTenant := new(coremocks.ITenantService)
// 	mockCoreUser := new(coremocks.IUser)
// 	filterUser := []entity.QueryDBClause{
// 		{
// 			Clause: entity.QueryClauseAnd,
// 			Queries: []entity.QueryDB{
// 				{
// 					Key:       "email",
// 					Value:     mockUser.Email,
// 					Condition: string(entity.QueryFirebaseEqual),
// 				},
// 			},
// 		},
// 	}

// 	filterTenant := []entity.QueryDBClause{
// 		{
// 			Clause: entity.QueryClauseAnd,
// 			Queries: []entity.QueryDB{
// 				{
// 					Key:       "name",
// 					Value:     mockUser.Email,
// 					Condition: string(entity.QueryFirebaseEqual),
// 				},
// 			},
// 		},
// 	}

// 	mockCoreUser.On("GetByFilterOne", ctx, filterUser).Return(mockAccountUser, nil).Once()
// 	mockTenant.On("GetByFilterMany", ctx, filterTenant[0].Queries).Return(mockTenants, nil).Once()

// 	userSvc, err := service.NewUserService(mockCoreUser, mockTenant, nil)
// 	require.NoError(t, err)
// 	assert.NotNil(t, userSvc)

// 	// Test with nil user
// 	result, err := userSvc.Create(ctx, mockUserNil)
// 	require.Error(t, err)
// 	assert.Equal(t, "user is required", err.Error())
// 	require.Nil(t, result)

// 	// Test with empty user
// 	result, err = userSvc.Create(ctx, &mockUserEmpty)
// 	require.Error(t, err)
// 	assert.Equal(t, "user is required", err.Error())
// 	require.Nil(t, result)

// 	// Test with invalid user
// 	result, err = userSvc.Create(ctx, mockUserNInvalid)
// 	require.Error(t, err)
// 	assert.Equal(t, "invalid user ID", err.Error())
// 	require.Nil(t, result)

// 	// Test with user already exists
// 	result, err = userSvc.Create(ctx, mockAccountUser)
// 	require.Error(t, err)
// 	assert.Equal(t, "user already exists", err.Error())
// 	require.Nil(t, result)

// 	mockCoreUser.AssertExpectations(t)
// 	mockTenant.AssertExpectations(t)

// 	// Test with error to filter user
// 	mockCoreUser.On("GetByFilterOne", ctx, filterUser).Return(nil, errors.New("error to filter user")).Once()
// 	result, err = userSvc.Create(ctx, mockAccountUser)
// 	require.Error(t, err)
// 	assert.Equal(t, "error to filter user", err.Error())
// 	require.Nil(t, result)

// 	mockCoreUser.AssertExpectations(t)
// 	mockTenant.AssertExpectations(t)

// 	// Test with tenant
// 	mockTenant.On("GetByFilterMany", ctx, filterTenant[0].Queries).Return(nil, errors.New("tenant exists")).Once()
// 	mockCoreUser.On("GetByFilterOne", ctx, filterUser).Return(nil, nil).Once()

// 	userSvc2, err := service.NewUserService(mockCoreUser, mockTenant, nil)
// 	require.NoError(t, err)
// 	require.Nil(t, err)
// 	require.NotNil(t, userSvc2)
// 	assert.NotNil(t, userSvc2)

// 	// Test with tenant exists
// 	resultTenant, err2 := userSvc2.Create(ctx, mockAccountUser)
// 	require.Error(t, err2)
// 	assert.Equal(t, "tenant exists", err2.Error())
// 	require.Nil(t, resultTenant)

// 	mockCoreUser.AssertExpectations(t)
// 	mockTenant.AssertExpectations(t)

// 	// Test create user error
// 	mockTenant.On("GetByFilterMany", ctx, filterTenant[0].Queries).Return(nil, nil).Once()
// 	mockCoreUser.On("GetByFilterOne", ctx, filterUser).Return(nil, nil).Once()
// 	mockCoreUser.On("Create", ctx, mockUser).Return(nil, errors.New("error to register a new user")).Once()

// 	userSvc3, err := service.NewUserService(mockCoreUser, mockTenant, nil)
// 	require.NoError(t, err)
// 	assert.NotNil(t, userSvc3)

// 	resultTenant, err2 = userSvc3.Create(ctx, mockAccountUser)
// 	require.Error(t, err2)
// 	assert.Equal(t, "error to register a new user", err2.Error())
// 	require.Nil(t, resultTenant)

// 	mockCoreUser.AssertExpectations(t)
// 	mockTenant.AssertExpectations(t)

// 	// Test with is nil
// 	mockTenant.On("GetByFilterMany", ctx, filterTenant[0].Queries).Return(nil, nil).Once()
// 	mockCoreUser.On("GetByFilterOne", ctx, filterUser).Return(nil, nil).Once()
// 	mockCoreUser.On("Create", ctx, mockUser).Return(mockUser, nil).Once()
// 	mockTenant.On("Create", mockTenantData).Return(nil, errors.New("error to register a tenant")).Once()

// 	userSvc3, err = service.NewUserService(mockCoreUser, mockTenant, nil)
// 	require.NoError(t, err)
// 	assert.NotNil(t, userSvc3)

// 	resultTenant, err2 = userSvc3.Create(ctx, mockAccountUser)
// 	require.Error(t, err2)
// 	require.Equal(t, "error to register a tenant", err2.Error())
// 	require.Nil(t, resultTenant)

// 	mockCoreUser.AssertExpectations(t)

// 	// Test with is nil
// 	mockTenant.On("GetByFilterMany", ctx, filterTenant[0].Queries).Return(nil, nil).Once()
// 	mockCoreUser.On("GetByFilterOne", ctx, filterUser).Return(nil, nil).Once()
// 	mockCoreUser.On("Create", ctx, mockUser).Return(mockUser, nil).Once()
// 	mockTenant.On("Create", mockTenantData).Return(mockTenantData, nil).Once()

// 	userSvc3, err = service.NewUserService(mockCoreUser, mockTenant, nil)
// 	require.NoError(t, err)
// 	assert.NotNil(t, userSvc3)

// 	resultTenant, err2 = userSvc3.Create(ctx, mockAccountUser)
// 	require.NoError(t, err2)
// 	require.Nil(t, err2)
// 	require.NotNil(t, resultTenant)

// 	mockCoreUser.AssertExpectations(t)
// }
