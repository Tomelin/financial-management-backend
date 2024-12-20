package entity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
)

type PlanTestSuite struct {
	suite.Suite
	features     *entity.PlanFeatures
	planResponse *entity.PlanResponse
}

func (s *PlanTestSuite) SetupTest() {
	s.features = &entity.PlanFeatures{
		Name:        "WalletShared",
		Description: "Wallet shared with other users",
		Count:       1,
		Unlimited:   false,
	}

	s.planResponse = &entity.PlanResponse{
		ID:          uuid.New().String(),
		Name:        "Bronze",
		Description: "Bronze plan",
		Price:       0.0,
		Features:    []entity.PlanFeatures{*s.features},
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func (s *PlanTestSuite) TearDownTest() {
	s.features = nil
	s.planResponse = nil
}

func (s *PlanTestSuite) TestNewPlan_Error_Nil() {

	planResponse, err := entity.NewPlan(nil)
	s.Error(err)
	s.Nil(planResponse)
	s.EqualError(err, "plan is required")
	s.Equal(err.Error(), "plan is required")
}

func (s *PlanTestSuite) TestNewPlan_Error_IsEmpty() {

	planResponse, err := entity.NewPlan(&entity.PlanResponse{})
	s.Error(err)
	s.Nil(planResponse)
	s.EqualError(err, "plan is required")
	s.Equal(err.Error(), "plan is required")
}

func (s *PlanTestSuite) TestNewPlan_Error_NameIsEmpty() {

	s.planResponse.Name = ""
	planResponse, err := entity.NewPlan(s.planResponse)
	s.Error(err)
	s.Nil(planResponse)
	s.EqualError(err, "plan name is required")
	s.Equal(err.Error(), "plan name is required")
}

func (s *PlanTestSuite) TestNewPlan_Error_Name_And_Feature() {

	s.planResponse.Features = nil
	s.planResponse.Name = "Gold"
	planResponse, err := entity.NewPlan(s.planResponse)
	s.Error(err)
	s.Nil(planResponse)
	s.EqualError(err, "features are required")
	s.Equal(err.Error(), "features are required")
}

func (s *PlanTestSuite) TestNewPlan_Error_Invalid_ID() {

	planResponse, _ := entity.NewPlan(s.planResponse)
	planResponse.ID = "invalid"
	err := planResponse.Validate()

	s.Error(err)
	s.EqualError(err, "invalid ID")
	s.Equal(err.Error(), "invalid ID")
}

func (s *PlanTestSuite) TestNewPlan_SetUpdate_Success() {

	planResponse, _ := entity.NewPlan(s.planResponse)
	planResponse.SetUpdate()
	s.NotNil(planResponse.UpdatedAt)
	s.IsType(time.Time{}, planResponse.UpdatedAt)
}

func (s *PlanTestSuite) TestNewPlan_Error_FeatureName() {

	planResponse, _ := entity.NewPlan(s.planResponse)
	planResponse.Features[0].Name = ""
	err := planResponse.Features[0].Validate()
	s.EqualError(err, "plan name is required")
	s.Equal(err.Error(), "plan name is required")
}

func (s *PlanTestSuite) TestNewPlan_Error_FeatureIsNil() {

	s.features = nil
	b := s.features.IsEmpty(s.features)
	s.True(b)
}

func (s *PlanTestSuite) TestNewPlan_Success() {

	planResponse, err := entity.NewPlan(s.planResponse)

	s.NoError(err)
	s.NotNil(planResponse)
	s.NotEmpty(planResponse.ID)
	s.Equal(s.planResponse.Name, planResponse.Name)
	s.Equal(s.planResponse.Description, planResponse.Description)

	err = planResponse.Features[0].Validate()
	s.NoError(err)
	s.Equal(s.planResponse.Features[0].Name, planResponse.Features[0].Name)
	s.Nil(err)
}

func TestRunPlanTestSuite(t *testing.T) {
	suite.Run(t, new(PlanTestSuite))
}
