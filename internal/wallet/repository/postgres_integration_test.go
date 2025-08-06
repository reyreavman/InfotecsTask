package postgres

import (
	"context"
	"errors"
	"infotecstechtask/internal/models"
	"infotecstechtask/pkg/database"
	"infotecstechtask/test/testutils"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WalletRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testutils.PGTestContainer
	repo        *WalletRepository
	fixtures    *testutils.FixtureManager
	dataLoader  *testutils.DataLoader
	ctx         context.Context
}

func (suite *WalletRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationsPath := filepath.Join(dir, "../../../init/migrations")

	container, err := testutils.StartPGContainer(suite.ctx, migrationsPath)
	if err != nil {
		log.Fatalf("Failed to start test container: %v", err)
	}
	suite.pgContainer = container

	client := database.NewClientWithPool(container.Pool)
	suite.repo = NewWalletRepository(client)

	suite.fixtures = testutils.NewFixtureManager(container.Pool)
	suite.dataLoader = testutils.NewDataLoader()
}

func (suite *WalletRepositoryTestSuite) TearDownSuite() {
	if suite.pgContainer != nil {
		if err := suite.pgContainer.Close(suite.ctx); err != nil {
			log.Printf("Failed to close test container: %v", err)
		}
	}
}

func (suite *WalletRepositoryTestSuite) BeforeTest(_, _ string) {
	_, err := suite.pgContainer.Pool.Exec(suite.ctx, "TRUNCATE TABLE wallets, transactions CASCADE")
	if err != nil {
		log.Printf("TRUNCATE error: %v", err)
	}

	assert.NoError(suite.T(), err)
}

func TestWalletRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(WalletRepositoryTestSuite))
}

func (suite *WalletRepositoryTestSuite) TestGetWalletSuccess() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	var expected models.Wallet
	err = suite.dataLoader.LoadJSONFixture("wallets/wallet.json", &expected)
	suite.Require().NoError(err)

	actual, err := suite.repo.GetWallet(suite.ctx, expected.ID)
	suite.Require().NoError(err)

	suite.Assert().Equal(expected, *actual)
}

func (suite *WalletRepositoryTestSuite) TestGetNonExistentWallet() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	nonExistentUUID := uuid.MustParse("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a20")
	_, err = suite.repo.GetWallet(suite.ctx, nonExistentUUID)
	suite.Require().Error(err)
}

func (suite *WalletRepositoryTestSuite) TestUpdateWalletSuccess() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	var walletToUpdate models.Wallet
	err = suite.dataLoader.LoadJSONFixture("wallets/wallet.json", &walletToUpdate)
	suite.Require().NoError(err)

	var expected models.Wallet
	err = suite.dataLoader.LoadJSONFixture("wallets/updated_wallet.json", &expected)
	suite.Require().NoError(err)

	walletToUpdate.Balance = expected.Balance
	actual, err := suite.repo.UpdateWallet(suite.ctx, &walletToUpdate)
	suite.Require().NoError(err)

	suite.Assert().Equal(expected, *actual)
}

func (suite *WalletRepositoryTestSuite) TestUpdateWalletWithNegativeBalance() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	var wallet models.Wallet
	err = suite.dataLoader.LoadJSONFixture("wallets/wallet.json", &wallet)
	suite.Require().NoError(err)

	updatedWallet := models.Wallet{
		ID:      wallet.ID,
		Balance: -50,
	}
	_, err = suite.repo.UpdateWallet(suite.ctx, &updatedWallet)
	suite.Require().Error(err, "Expected error for negative balance")

	var pgErr *pgconn.PgError
	suite.Require().True(errors.As(err, &pgErr), "Error should be of type PgError")

	actualWallet, err := suite.repo.GetWallet(suite.ctx, wallet.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal(wallet.Balance, actualWallet.Balance)
}
