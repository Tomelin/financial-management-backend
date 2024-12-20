package entity_test

import (
	"testing"
	"time"

	"github.com/Tomelin/financial-management-backend/internal/core/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TenantTestSuite struct {
	suite.Suite
	tenant   *entity.TenantResponse
	plan     *entity.PlanResponse
	tenantID string
	userID   string
}

func (s *TenantTestSuite) SetupTest() {
	s.plan = &entity.PlanResponse{
		ID:          uuid.New().String(),
		Name:        "Bronze",
		Description: "Bronze plan",
		Price:       0.0,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	s.tenantID = uuid.New().String()
	s.userID = uuid.New().String()
	s.tenant = &entity.TenantResponse{
		Name:      "user@domain.com",
		Alias:     "synera",
		OwnerID:   s.userID,
		Users:     []string{s.userID},
		ID:        s.tenantID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Plan:      *s.plan,
		Wallets:   []string{s.tenantID},
	}
}

func (s *TenantTestSuite) TearDownTest() {
	s.plan = nil
	s.tenant = nil
	s.tenantID = ""
	s.userID = ""
}

func (s *TenantTestSuite) TestNewTenant_Error_Nil() {

	tenant, err := entity.NewTenant(nil)
	s.Error(err)
	s.Nil(tenant)
	s.EqualError(err, "tenant is required")
	s.Equal(err.Error(), "tenant is required")
}

func (s *TenantTestSuite) TestNewTenant_Error_Empty() {

	tenant, err := entity.NewTenant(&entity.TenantResponse{})
	s.Error(err)
	s.Nil(tenant)
	s.EqualError(err, "tenant is required")
	s.Equal(err.Error(), "tenant is required")
}

func (s *TenantTestSuite) TestNewTenant_Error_Name() {

	s.tenant.Name = ""
	tenant, err := entity.NewTenant(s.tenant)
	s.NotNil(err)
	s.Error(err)
	s.EqualError(err, "tenant name is required")
	s.Nil(tenant)
}

func (s *TenantTestSuite) TestNewTenant_Error_InvalidName() {

	s.tenant.Name = "user"
	tenant, err := entity.NewTenant(s.tenant)
	s.NotNil(err)
	s.Error(err)
	s.EqualError(err, "tenant name is invalid. Tenant name must be a valid email")
	s.Nil(tenant)
}

func (s *TenantTestSuite) TestNewTenant_Error_InvalidOwnerID() {

	s.tenant.OwnerID = "invalid"
	tenant, err := entity.NewTenant(s.tenant)
	s.NotNil(err)
	s.Error(err)
	s.EqualError(err, "invalid owner id")
	s.Nil(tenant)
}

func (s *TenantTestSuite) TestNewTenant_Error_InvalidPlan() {

	s.tenant.Plan.ID = "invalid"
	tenant, err := entity.NewTenant(s.tenant)
	s.NotNil(err)
	s.Error(err)
	s.EqualError(err, "invalid plan id")
	s.Nil(tenant)
}

func (s *TenantTestSuite) TestNewTenant_Error_TenantID() {

	s.tenant.ID = "invalid"
	tenant, err := entity.NewTenant(s.tenant)
	s.NotNil(err)
	s.Error(err)
	s.EqualError(err, "invalid tenant id")
	s.Nil(tenant)
}

func (s *TenantTestSuite) TestNewTenant_Error_InvalidWalletID() {

	s.tenant.Wallets = []string{s.tenantID, "invalid"}
	tenant, err := entity.NewTenant(s.tenant)
	s.NotNil(err)
	s.Error(err)
	s.EqualError(err, "invalid wallet id")
	s.Nil(tenant)
}
func (s *TenantTestSuite) TestNewTenant_Error_IsEmpty() {

	b := s.tenant.IsEmpty(s.tenant)
	s.True(!b)
}

func (s *TenantTestSuite) TestNewTenant_Success_CreateTenantID() {

	s.tenant.ID = ""
	tenant, err := entity.NewTenant(s.tenant)
	s.Nil(err)
	s.NoError(err)
	s.NotNil(tenant)
	s.Equal(s.tenant.Name, tenant.Name)
}

func (s *TenantTestSuite) TestNewTenant_Success() {

	tenant, err := entity.NewTenant(s.tenant)
	s.Nil(err)
	s.NoError(err)
	s.NotNil(tenant)
	s.Equal(s.tenant.Name, tenant.Name)
}

func TestRunTenantTestSuite(t *testing.T) {
	suite.Run(t, new(TenantTestSuite))
}
