package postgres

import (
	"context"
	"infotecstechtask/internal/models"
	"infotecstechtask/pkg/database"
	"infotecstechtask/test/testutils"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testutils.PGTestContainer
	repo        *TransactionRepository
	fixtures    *testutils.FixtureManager
	dataLoader  *testutils.DataLoader
	ctx         context.Context
}

func (suite *TransactionRepositoryTestSuite) SetupSuite() {
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
	suite.repo = NewTransactionRepository(client)

	suite.fixtures = testutils.NewFixtureManager(container.Pool)
	suite.dataLoader = testutils.NewDataLoader()
}

func (suite *TransactionRepositoryTestSuite) TearDownSuite() {
	if suite.pgContainer != nil {
		if err := suite.pgContainer.Close(suite.ctx); err != nil {
			log.Printf("Failed to close test container: %v", err)
		}
	}
}

func (suite *TransactionRepositoryTestSuite) BeforeTest(_, _ string) {
	_, err := suite.pgContainer.Pool.Exec(suite.ctx, "TRUNCATE TABLE wallets, transactions CASCADE")
	if err != nil {
		log.Printf("TRUNCATE error: %v", err)
	}

	assert.NoError(suite.T(), err)
}

func TestWalletRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}

func (suite *TransactionRepositoryTestSuite) TestGetAllTransactionsSuccess() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	err = suite.fixtures.ApplySQLFixture(suite.ctx, "transactions/transactions.sql")
	suite.Require().NoError(err)

	var expected []*models.Transaction
	err = suite.dataLoader.LoadJSONFixture("transactions/all_transactions.json", &expected)
	suite.Require().NoError(err)

	actual, err := suite.repo.GetAllTransactions(suite.ctx)
	suite.Require().NoError(err)

	suite.Assert().ElementsMatch(expected, actual)
}

func (suite *TransactionRepositoryTestSuite) TestGetAllTransactionsWithEmptyReturnedSliceSuccess() {
	actual, err := suite.repo.GetAllTransactions(suite.ctx)
	suite.Require().NoError(err)

	suite.Assert().Empty(actual)
}

func (suite *TransactionRepositoryTestSuite) TestGetTransactionsSuccess() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	err = suite.fixtures.ApplySQLFixture(suite.ctx, "transactions/transactions.sql")
	suite.Require().NoError(err)

	var allTransactions []*models.Transaction
	err = suite.dataLoader.LoadJSONFixture("transactions/all_transactions.json", &allTransactions)
	suite.Require().NoError(err)
	var oneTransaction []*models.Transaction
	err = suite.dataLoader.LoadJSONFixture("transactions/one_transaction.json", &oneTransaction)
	suite.Require().NoError(err)
	var twoTransactions []*models.Transaction
	err = suite.dataLoader.LoadJSONFixture("transactions/two_transactions.json", &twoTransactions)
	suite.Require().NoError(err)

	testCases := []struct {
		count    int
		expected []*models.Transaction
	}{
		{0, make([]*models.Transaction, 0)},
		{1, oneTransaction},
		{2, twoTransactions},
		{3, allTransactions},
		{5, allTransactions},
	}

	for _, tc := range testCases {
		actual, err := suite.repo.GetTransactions(suite.ctx, tc.count)
		suite.Require().NoError(err)

		suite.Assert().ElementsMatch(tc.expected, actual)
	}
}
