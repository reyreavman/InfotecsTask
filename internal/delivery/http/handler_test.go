package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"infotecstechtask/internal/facade"
	"infotecstechtask/internal/models"
	"infotecstechtask/internal/payment"
	"infotecstechtask/internal/wallet"
	"infotecstechtask/test/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestInfrastructure struct {
	suite.Suite
	rGroup     *gin.Engine
	dataLoader *testutils.DataLoader
	validate   *validator.Validate
}

func (tf *TestInfrastructure) SetupSuite() {
	tf.rGroup = gin.Default()
	tf.dataLoader = testutils.NewDataLoader()
	tf.validate = validator.New()
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(TestInfrastructure))
}

func (tf *TestInfrastructure) AfterTest(_, _ string) {
	tf.rGroup = gin.Default()
}

func (tf *TestInfrastructure) TestCreateTransactionSuccess() {
	var request models.CreateTransactionRequest
	err := tf.dataLoader.LoadJSONFixture("transactions/request/create_transaction_request.json", &request)
	tf.Require().NoError(err)

	var response models.TransactionResponse
	err = tf.dataLoader.LoadJSONFixture("transactions/response/transaction_response.json", &response)
	tf.Require().NoError(err)

	mockFacade := new(facade.MockFacade)
	mockFacade.On(
		"CreateTransaction",
		mock.Anything,
		&request,
	).Return(&response, nil)

	RegisterHTTPEndpoints(&tf.rGroup.RouterGroup, mockFacade, tf.validate)

	jsonRequest, err := json.Marshal(request)
	tf.Require().NoError(err)
	body := bytes.NewBuffer(jsonRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, FULL_SEND, body)
	req.Header.Add("Content-Type", "application/json")

	tf.rGroup.ServeHTTP(w, req)

	expectedResponseBody, err := json.Marshal(response)
	tf.Require().NoError(err)

	tf.Assert().Equal(200, w.Code)
	tf.Assert().Equal(string(expectedResponseBody), w.Body.String())
}

func (tf *TestInfrastructure) TestCreateTransactionErrSenderWalletNotFound() {
	var request models.CreateTransactionRequest
	err := tf.dataLoader.LoadJSONFixture("transactions/request/create_transaction_request.json", &request)
	tf.Require().NoError(err)

	var expectedErr models.Error
	err = tf.dataLoader.LoadJSONFixture("errors/sender_wallet_not_found.json", &expectedErr)
	tf.Require().NoError(err)

	mockFacade := new(facade.MockFacade)
	mockFacade.On(
		"CreateTransaction",
		mock.Anything,
		&request,
	).Return(nil, payment.ErrSenderWalletNotFound)

	RegisterHTTPEndpoints(&tf.rGroup.RouterGroup, mockFacade, tf.validate)

	jsonRequest, err := json.Marshal(request)
	tf.Require().NoError(err)
	body := bytes.NewBuffer(jsonRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, FULL_SEND, body)
	req.Header.Add("Content-Type", "application/json")

	tf.rGroup.ServeHTTP(w, req)

	expectedResponseBody, err := json.Marshal(expectedErr)
	tf.Require().NoError(err)

	tf.Assert().Equal(400, w.Code)
	tf.Assert().Equal(string(expectedResponseBody), w.Body.String())
}

func (tf *TestInfrastructure) TestCreateTransactionErrRecipientWalletNotFound() {
	var request models.CreateTransactionRequest
	err := tf.dataLoader.LoadJSONFixture("transactions/request/create_transaction_request.json", &request)
	tf.Require().NoError(err)

	var expectedErr models.Error
	err = tf.dataLoader.LoadJSONFixture("errors/recipient_wallet_not_found.json", &expectedErr)
	tf.Require().NoError(err)

	mockFacade := new(facade.MockFacade)
	mockFacade.On(
		"CreateTransaction",
		mock.Anything,
		&request,
	).Return(nil, payment.ErrRecipientWalletNotFound)

	RegisterHTTPEndpoints(&tf.rGroup.RouterGroup, mockFacade, tf.validate)

	jsonRequest, err := json.Marshal(request)
	tf.Require().NoError(err)
	body := bytes.NewBuffer(jsonRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, FULL_SEND, body)
	req.Header.Add("Content-Type", "application/json")

	tf.rGroup.ServeHTTP(w, req)

	expectedResponseBody, err := json.Marshal(expectedErr)
	tf.Require().NoError(err)

	tf.Assert().Equal(400, w.Code)
	tf.Assert().Equal(string(expectedResponseBody), w.Body.String())
}

func (tf *TestInfrastructure) TestCreateTransactionSenderNotHaveEnoughBalance() {
	var request models.CreateTransactionRequest
	err := tf.dataLoader.LoadJSONFixture("transactions/request/create_transaction_request.json", &request)
	tf.Require().NoError(err)

	var expectedResp models.TransactionResponse
	err = tf.dataLoader.LoadJSONFixture("transactions/response/failed_transaction_response.json", &expectedResp)
	tf.Require().NoError(err)

	mockFacade := new(facade.MockFacade)
	mockFacade.On(
		"CreateTransaction",
		mock.Anything,
		&request,
	).Return(&expectedResp, nil)

	RegisterHTTPEndpoints(&tf.rGroup.RouterGroup, mockFacade, tf.validate)

	jsonRequest, err := json.Marshal(request)
	tf.Require().NoError(err)
	body := bytes.NewBuffer(jsonRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, FULL_SEND, body)
	req.Header.Add("Content-Type", "application/json")

	tf.rGroup.ServeHTTP(w, req)

	expectedResponseBody, err := json.Marshal(expectedResp)
	tf.Require().NoError(err)

	tf.Assert().Equal(200, w.Code)
	tf.Assert().Equal(string(expectedResponseBody), w.Body.String())
}

func (tf *TestInfrastructure) TestGetAllTransactionSuccess() {
	var allTransactions []*models.Transaction
	err := tf.dataLoader.LoadJSONFixture("transactions/all_transactions.json", &allTransactions)
	tf.Require().NoError(err)

	mockFacade := new(facade.MockFacade)

	mockFacade.On(
		"GetAllTransactions",
		mock.Anything,
	).Return(models.ToTransactionResponses(allTransactions), nil)

	RegisterHTTPEndpoints(&tf.rGroup.RouterGroup, mockFacade, tf.validate)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, FULL_TRANSACTIONS, nil)

	tf.rGroup.ServeHTTP(w, req)

	expectedResponseBody, err := json.Marshal(models.ToTransactionResponses(allTransactions))
	tf.Require().NoError(err)

	tf.Assert().Equal(200, w.Code)
	tf.Assert().Equal(string(expectedResponseBody), w.Body.String())
}

func (tf *TestInfrastructure) TestGetTransactions() {
	var allTransactions []*models.Transaction
	err := tf.dataLoader.LoadJSONFixture("transactions/all_transactions.json", &allTransactions)
	tf.Require().NoError(err)
	var oneTransaction []*models.Transaction
	err = tf.dataLoader.LoadJSONFixture("transactions/one_transaction.json", &oneTransaction)
	tf.Require().NoError(err)
	var twoTransactions []*models.Transaction
	err = tf.dataLoader.LoadJSONFixture("transactions/two_transactions.json", &twoTransactions)
	tf.Require().NoError(err)

	mockFacade := new(facade.MockFacade)
	RegisterHTTPEndpoints(&tf.rGroup.RouterGroup, mockFacade, tf.validate)

	testCases := []struct {
		count    int
		expected []*models.Transaction
	}{
		{1, oneTransaction},
		{2, twoTransactions},
		{3, allTransactions},
		{5, allTransactions},
	}

	var builder strings.Builder

	for _, tc := range testCases {
		mockFacade.On(
			"GetTransactions",
			mock.Anything,
			tc.count,
		).Return(models.ToTransactionResponses(tc.expected), nil)

		builder.Write([]byte(FULL_TRANSACTIONS))
		builder.Write([]byte(fmt.Sprintf("?count=%d", tc.count)))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, builder.String(), nil)

		tf.rGroup.ServeHTTP(w, req)

		expectedResponseBody, err := json.Marshal(models.ToTransactionResponses(tc.expected))
		tf.Require().NoError(err)

		tf.Assert().Equal(200, w.Code)
		tf.Assert().Equal(string(expectedResponseBody), w.Body.String())
		builder.Reset()
	}
}

func (tf *TestInfrastructure) TestGetWalletSuccess() {
	var expectedResp models.Wallet
	err := tf.dataLoader.LoadJSONFixture("wallets/wallet.json", &expectedResp)
	tf.Require().NoError(err)

	mockFacade := new(facade.MockFacade)
	mockFacade.On(
		"GetWallet",
		mock.Anything,
		expectedResp.ID,
	).Return(&expectedResp, nil)

	validate := validator.New()

	RegisterHTTPEndpoints(&tf.rGroup.RouterGroup, mockFacade, validate)

	path := strings.Replace(FULL_GET_WALLET_BALANCE, ":walletId", expectedResp.ID.String(), 1)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, path, nil)

	tf.rGroup.ServeHTTP(w, req)

	expectedResponseBody, err := json.Marshal(expectedResp)
	tf.Require().NoError(err)

	tf.Assert().Equal(200, w.Code)
	tf.Assert().Equal(string(expectedResponseBody), w.Body.String())
}

func (tf *TestInfrastructure) TestGetWalletNotFound() {
	var expectedResp models.Wallet
	err := tf.dataLoader.LoadJSONFixture("wallets/wallet.json", &expectedResp)
	tf.Require().NoError(err)

	mockFacade := new(facade.MockFacade)
	mockFacade.On(
		"GetWallet",
		mock.Anything,
		expectedResp.ID,
	).Return(nil, wallet.ErrWalletNotFound)

	RegisterHTTPEndpoints(&tf.rGroup.RouterGroup, mockFacade, tf.validate)

	path := strings.Replace(FULL_GET_WALLET_BALANCE, ":walletId", expectedResp.ID.String(), 1)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, path, nil)

	tf.rGroup.ServeHTTP(w, req)

	tf.Require().NoError(err)

	tf.Assert().Equal(404, w.Code)
}
