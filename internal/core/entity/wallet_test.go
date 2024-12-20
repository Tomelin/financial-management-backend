package entity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
)

type WalletTestSuite struct {
	suite.Suite
	walletResponse *entity.WalletResponse
}

func (s *WalletTestSuite) SetupTest() {
	walletID := uuid.New().String()
	ownerID := uuid.New().String()
	tenantID := uuid.New().String()
	s.walletResponse = &entity.WalletResponse{
		ID:          walletID,
		Name:        "Wallet",
		Description: "Wallet description",
		OwnerID:     ownerID,
		TenantID:    tenantID,
		Balance:     0,
		Currency:    "BRL",
	}
}

func (s *WalletTestSuite) TearDownTest() {
	s.walletResponse = nil
}

func (s *WalletTestSuite) TestNewWallet_Error_Nil() {
	walletResponse, err := entity.NewWallet(nil)
	s.Nil(walletResponse)
	s.Equal(err.Err, "wallet is required")
}
func (s *WalletTestSuite) TestNewWallet_Error_Empty() {
	walletResponse, err := entity.NewWallet(&entity.WalletResponse{})
	s.Nil(walletResponse)
	s.Equal(err.Err, "wallet is required")
}

func (s *WalletTestSuite) TestNewWallet_EmptyID() {

	s.walletResponse.ID = ""
	walletResponse, err := entity.NewWallet(s.walletResponse)
	s.NotNil(walletResponse)
	s.Nil(err)
	s.Equal(s.walletResponse.ID, walletResponse.ID)
}

func (s *WalletTestSuite) TestNewWallet_EmptyCurrency() {

	s.walletResponse.Currency = ""
	walletResponse, err := entity.NewWallet(s.walletResponse)
	s.NotNil(walletResponse)
	s.Nil(err)
	s.Equal("BRL", walletResponse.Currency)
}

func (s *WalletTestSuite) TestNewWallet_Error_OwnerID() {

	s.walletResponse.OwnerID = ""
	walletResponse, err := entity.NewWallet(s.walletResponse)
	s.Nil(walletResponse)
	s.NotNil(err)
	s.Equal(err.Err, "invalid UUID length: 0")
}

func (s *WalletTestSuite) TestNewWallet_Error_TenantID() {

	s.walletResponse.TenantID = ""
	walletResponse, err := entity.NewWallet(s.walletResponse)
	s.Nil(walletResponse)
	s.NotNil(err)
	s.Equal(err.Err, "invalid UUID length: 0")
}

func (s *WalletTestSuite) TestNewWallet_Error_VaidateEmpty() {

	s.walletResponse = nil
	s.walletResponse.Validate()
	err := s.walletResponse.Validate()
	s.NotNil(err)
	s.Equal(err.Err, "wallet cannot be empty")
}

func (s *WalletTestSuite) TestNewWallet_Error_VaidateEmptyID() {

	s.walletResponse.ID = ""
	s.walletResponse.Validate()
	err := s.walletResponse.Validate()
	s.NotNil(err)
	s.Equal(err.Err, "id is required")
}

func (s *WalletTestSuite) TestNewWallet_Error_VaidateID() {

	s.walletResponse.ID = "invalid"
	s.walletResponse.Validate()
	err := s.walletResponse.Validate()
	s.NotNil(err)
	s.Equal(err.Err, "invalid UUID length: 7")
}

func (s *WalletTestSuite) TestNewWallet_SetBalanceZero() {

	err := s.walletResponse.SetBalance(0.0)
	s.Nil(err)
	s.Equal(0.0, s.walletResponse.Balance)
}

func (s *WalletTestSuite) TestNewWallet_SetBalanceFloat100() {

	err := s.walletResponse.SetBalance(100)
	s.Nil(err)
	s.Equal(float64(100), s.walletResponse.Balance)
}

func (s *WalletTestSuite) TestNewWallet_SetUpdate() {

	s.Equal(s.walletResponse.UpdatedAt, time.Time{})
	s.Empty(s.walletResponse.UpdatedAt)
	s.walletResponse.SetUpdate()
	s.NotEqual(s.walletResponse.UpdatedAt, time.Time{})
	s.NotEmpty(s.walletResponse.UpdatedAt)
}

func (s *WalletTestSuite) TestNewWallet_WalletByPlan() {
	silverPlan := "silver"
	silver := 3
	goldPlan := "gold"
	gold := -1
	bronzelan := "bronze"
	bronze := 1
	emptyPLan := ""

	s.Equal(silver, s.walletResponse.WalletByPlan(&silverPlan))
	s.Equal(gold, s.walletResponse.WalletByPlan(&goldPlan))
	s.Equal(bronze, s.walletResponse.WalletByPlan(&bronzelan))
	s.NotEqual(silver, s.walletResponse.WalletByPlan(&goldPlan))
	s.Equal(bronze, s.walletResponse.WalletByPlan(&emptyPLan))
}

func (s *WalletTestSuite) TestNewWallet_SharedByPlan() {
	silverPlan := "silver"
	silver := 2
	goldPlan := "gold"
	gold := -1
	bronzelan := "bronze"
	bronze := 1
	emptyPLan := ""

	s.Equal(silver, s.walletResponse.SharedByPlan(&silverPlan))
	s.Equal(gold, s.walletResponse.SharedByPlan(&goldPlan))
	s.Equal(bronze, s.walletResponse.SharedByPlan(&bronzelan))
	s.NotEqual(silver, s.walletResponse.SharedByPlan(&goldPlan))
	s.Equal(bronze, s.walletResponse.SharedByPlan(&emptyPLan))
}

func (s *WalletTestSuite) TestNewWallet_Error_SetBalance100() {

	err := s.walletResponse.SetBalance(100)
	s.Nil(err)
	s.NotEqual(100, s.walletResponse.Balance)
}

func (s *WalletTestSuite) TestNewWallet_Validate() {

	s.walletResponse.Name = ""
	walletResponse, err := entity.NewWallet(s.walletResponse)
	s.NotNil(walletResponse)
	s.Nil(err)
	s.Equal("MyWallet", walletResponse.Name)
}

func (s *WalletTestSuite) TestNewWallet_Success() {

	walletResponse, err := entity.NewWallet(s.walletResponse)
	s.NotNil(walletResponse)
	s.Nil(err)
	s.Equal(s.walletResponse.Name, walletResponse.Name)
}

func TestRunWalletTestSuite(t *testing.T) {
	suite.Run(t, new(WalletTestSuite))
}
