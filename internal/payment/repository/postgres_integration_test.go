package postgres

import (
	"context"
	"infotecstechtask/internal/models"
	"infotecstechtask/pkg/database"
	"infotecstechtask/test/testutils"
	"log"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PaymentRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testutils.PGTestContainer
	repo        *PaymentRepository
	fixtures    *testutils.FixtureManager
	dataLoader  *testutils.DataLoader
	ctx         context.Context
}

func (suite *PaymentRepositoryTestSuite) SetupSuite() {
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
	suite.repo = NewPaymentRepository(client)

	suite.fixtures = testutils.NewFixtureManager(container.Pool)
	suite.dataLoader = testutils.NewDataLoader()
}

func (suite *PaymentRepositoryTestSuite) TearDownSuite() {
	if suite.pgContainer != nil {
		if err := suite.pgContainer.Close(suite.ctx); err != nil {
			log.Printf("Failed to close test container: %v", err)
		}
	}
}

func (suite *PaymentRepositoryTestSuite) BeforeTest(_, _ string) {
	_, err := suite.pgContainer.Pool.Exec(suite.ctx, `TRUNCATE TABLE wallets, transactions CASCADE`)
	if err != nil {
		log.Printf("TRUNCATE error: %v", err)
	}

	assert.NoError(suite.T(), err)
}

func TestPaymentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentRepositoryTestSuite))
}

func (suite *PaymentRepositoryTestSuite) TestCreatePaymentSuccess() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	var sender models.Wallet
	err = suite.dataLoader.LoadJSONFixture("payments/addition/sender_wallet.json", &sender)
	suite.Require().NoError(err)

	var recipient models.Wallet
	err = suite.dataLoader.LoadJSONFixture("payments/addition/recipient_wallet.json", &recipient)
	suite.Require().NoError(err)

	var request models.CreateTransactionRequest
	err = suite.dataLoader.LoadJSONFixture("payments/transaction_request.json", &request)
	suite.Require().NoError(err)

	response, err := suite.repo.CreatePayment(suite.ctx, &request)
	suite.Require().NoError(err)

	suite.Assert().Equal(request.FromAddress, response.FromAddress.String())
	suite.Assert().Equal(request.ToAddress, response.ToAddress.String())
	suite.Assert().Equal(request.Amount, response.Amount)
	suite.Assert().Equal(models.Completed, response.Status)
	suite.Assert().NotNil(response.ID)
	suite.Assert().NotNil(response.CreatedAt)

	suite.verifyWalletBalance(sender.ID, sender.Balance-response.Amount)
	suite.verifyWalletBalance(recipient.ID, recipient.Balance+response.Amount)
	suite.verifyTransactionStatus(response.ID, models.Completed)
}

func (suite *PaymentRepositoryTestSuite) TestCreatePaymentInsufficientFunds() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	var sender models.Wallet
	err = suite.dataLoader.LoadJSONFixture("payments/addition/sender_wallet.json", &sender)
	suite.Require().NoError(err)

	var recipient models.Wallet
	err = suite.dataLoader.LoadJSONFixture("payments/addition/recipient_wallet.json", &recipient)
	suite.Require().NoError(err)

	var request models.CreateTransactionRequest
	err = suite.dataLoader.LoadJSONFixture("payments/transaction_request_insufficient_funds.json", &request)
	suite.Require().NoError(err)

	response, err := suite.repo.CreatePayment(suite.ctx, &request)
	suite.Require().NoError(err)

	suite.Assert().Equal(request.FromAddress, response.FromAddress.String())
	suite.Assert().Equal(request.ToAddress, response.ToAddress.String())
	suite.Assert().Equal(request.Amount, response.Amount)
	suite.Assert().Equal(models.Failed, response.Status)
	suite.Assert().NotNil(response.ID)
	suite.Assert().NotNil(response.CreatedAt)

	suite.verifyWalletBalance(sender.ID, sender.Balance)
	suite.verifyWalletBalance(recipient.ID, recipient.Balance)
	suite.verifyTransactionStatus(response.ID, models.Failed)
}

func (suite *PaymentRepositoryTestSuite) TestCreatePaymentSenderNotFound() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	senderID := uuid.New()

	var recipient models.Wallet
	err = suite.dataLoader.LoadJSONFixture("payments/addition/recipient_wallet.json", &recipient)
	suite.Require().NoError(err)

	request := &models.CreateTransactionRequest{
		FromAddress: senderID.String(),
		ToAddress:   recipient.ID.String(),
		Amount:      50.0,
	}

	response, err := suite.repo.CreatePayment(suite.ctx, request)
	suite.Assert().Nil(response)
	suite.Assert().Contains(err.Error(), "Sender wallet not found")

	var count int
	err = suite.pgContainer.Pool.QueryRow(suite.ctx,
		`SELECT COUNT(*) FROM transactions WHERE from_address = $1`,
		senderID,
	).Scan(&count)
	suite.Require().NoError(err)
	suite.Assert().Equal(0, count)
}

func (suite *PaymentRepositoryTestSuite) TestCreatePaymentRecipientNotFound() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	var sender models.Wallet
	err = suite.dataLoader.LoadJSONFixture("payments/addition/sender_wallet.json", &sender)
	suite.Require().NoError(err)

	recipientID := uuid.New()

	request := &models.CreateTransactionRequest{
		FromAddress: sender.ID.String(),
		ToAddress:   recipientID.String(),
		Amount:      50.0,
	}

	response, err := suite.repo.CreatePayment(suite.ctx, request)
	suite.Require().Error(err)
	suite.Assert().Nil(response)
	suite.Assert().Contains(err.Error(), "Recipient wallet not found")

	var count int
	err = suite.pgContainer.Pool.QueryRow(suite.ctx,
		`SELECT COUNT(*) FROM transactions WHERE to_address = $1`,
		recipientID,
	).Scan(&count)
	suite.Require().NoError(err)
	suite.Assert().Equal(0, count)
}

func (suite *PaymentRepositoryTestSuite) TestCreatePaymentConcurrentTransactions() {
	err := suite.fixtures.ApplySQLFixture(suite.ctx, "wallets/wallets.sql")
	suite.Require().NoError(err)

	var sender models.Wallet
	err = suite.dataLoader.LoadJSONFixture("payments/addition/sender_wallet.json", &sender)
	suite.Require().NoError(err)

	var recipient models.Wallet
	err = suite.dataLoader.LoadJSONFixture("payments/addition/recipient_wallet.json", &recipient)
	suite.Require().NoError(err)

	var request models.CreateTransactionRequest
	err = suite.dataLoader.LoadJSONFixture("payments/transaction_request.json", &request)
	suite.Require().NoError(err)

	concurrency := 2
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			ctx := context.Background()
			_, err := suite.repo.CreatePayment(ctx, &request)
			suite.Require().NoError(err)
		}()
	}

	wg.Wait()

	var actualSenderBalance, actualRecipientBalance float32
	err = suite.pgContainer.Pool.QueryRow(suite.ctx, `SELECT balance FROM wallets WHERE id = $1`, sender.ID).Scan(&actualSenderBalance)
	suite.Require().NoError(err)
	suite.Assert().Equal(sender.Balance-float32(concurrency)*request.Amount, actualSenderBalance)

	err = suite.pgContainer.Pool.QueryRow(suite.ctx, `SELECT balance FROM wallets WHERE id = $1`, recipient.ID).Scan(&actualRecipientBalance)
	suite.Require().NoError(err)
	suite.Assert().Equal(recipient.Balance+float32(concurrency)*request.Amount, actualRecipientBalance)

	var txCount int
	err = suite.pgContainer.Pool.QueryRow(suite.ctx,
		`SELECT COUNT(*) FROM transactions WHERE from_address = $1 AND status = $2`,
		sender.ID, models.Completed,
	).Scan(&txCount)
	suite.Require().NoError(err)
	suite.Assert().Equal(concurrency, txCount)
}

func (suite *PaymentRepositoryTestSuite) verifyWalletBalance(id uuid.UUID, expected float32) {
	var balance float32
	err := suite.pgContainer.Pool.QueryRow(
		context.Background(),
		"SELECT balance FROM wallets WHERE id = $1",
		id,
	).Scan(&balance)
	suite.Require().NoError(err)
	suite.Assert().Equal(expected, balance)
}

func (suite *PaymentRepositoryTestSuite) verifyTransactionStatus(id uuid.UUID, expected models.Status) {
	var status models.Status
	err := suite.pgContainer.Pool.QueryRow(
		context.Background(),
		"SELECT status FROM transactions WHERE id = $1",
		id,
	).Scan(&status)
	suite.Require().NoError(err)
	suite.Assert().Equal(expected, status)
}
