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
	var expectedResp models.WalletResponse
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
	var expectedResp models.WalletResponse
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

func (tf *TestInfrastructure) TestCreateTransactionWithNonValidRequest() {
	var requestWithoutFromAddress models.CreateTransactionRequest
	err := tf.dataLoader.LoadJSONFixture("transactions/request/create_transaction_request_without_from_address.json", &requestWithoutFromAddress)
	tf.Require().NoError(err)
	var requestWithNonValidFromAddress models.CreateTransactionRequest
	err = tf.dataLoader.LoadJSONFixture("transactions/request/create_transaction_request_with_non_valid_from_address.json", &requestWithNonValidFromAddress)
	tf.Require().NoError(err)
	var requestWithoutToAddress models.CreateTransactionRequest
	err = tf.dataLoader.LoadJSONFixture("transactions/request/create_transaction_request_without_to_address.json", &requestWithoutToAddress)
	tf.Require().NoError(err)
	var requestWithNonValidToAddress models.CreateTransactionRequest
	err = tf.dataLoader.LoadJSONFixture("transactions/request/create_transaction_request_with_non_valid_to_address.json", &requestWithNonValidToAddress)
	tf.Require().NoError(err)
	var requestWithoutAmount models.CreateTransactionRequest
	err = tf.dataLoader.LoadJSONFixture("transactions/request/create_transaction_request_without_amount.json", &requestWithoutAmount)
	tf.Require().NoError(err)
	var requestWithNegativeAmount models.CreateTransactionRequest
	err = tf.dataLoader.LoadJSONFixture("transactions/request/create_transaction_request_with_negative_amount.json", &requestWithNegativeAmount)
	tf.Require().NoError(err)

	var fromAddressRequired models.ValidationError
	err = tf.dataLoader.LoadJSONFixture("errors/from_address_required.json", &fromAddressRequired)
	tf.Require().NoError(err)
	var fromAddressNonValid models.ValidationError
	err = tf.dataLoader.LoadJSONFixture("errors/from_address_non_valid.json", &fromAddressNonValid)
	tf.Require().NoError(err)
	var toAddressRequired models.ValidationError
	err = tf.dataLoader.LoadJSONFixture("errors/to_address_required.json", &toAddressRequired)
	tf.Require().NoError(err)
	var toAddressNonValid models.ValidationError
	err = tf.dataLoader.LoadJSONFixture("errors/to_address_non_valid.json", &toAddressNonValid)
	tf.Require().NoError(err)
	var amountRequired models.ValidationError
	err = tf.dataLoader.LoadJSONFixture("errors/amount_required.json", &amountRequired)
	tf.Require().NoError(err)
	var amountNegative models.ValidationError
	err = tf.dataLoader.LoadJSONFixture("errors/amount_negative.json", &amountNegative)
	tf.Require().NoError(err)

	RegisterHTTPEndpoints(&tf.rGroup.RouterGroup, nil, tf.validate)

	testCases := []struct {
		request          models.CreateTransactionRequest
		expectedResponse models.ValidationError
	}{
		{
			request:          requestWithoutFromAddress,
			expectedResponse: fromAddressRequired,
		},
		{
			request:          requestWithNonValidFromAddress,
			expectedResponse: fromAddressNonValid,
		},
		{
			request:          requestWithoutToAddress,
			expectedResponse: toAddressRequired,
		},
		{
			request:          requestWithNonValidToAddress,
			expectedResponse: toAddressNonValid,
		},
		{
			request:          requestWithoutAmount,
			expectedResponse: amountRequired,
		},
		{
			request:          requestWithNegativeAmount,
			expectedResponse: amountNegative,
		},
	}

	for _, tc := range testCases {
		jsonRequest, err := json.Marshal(tc.request)
		tf.Require().NoError(err)
		body := bytes.NewBuffer(jsonRequest)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, FULL_SEND, body)
		req.Header.Add("Content-Type", "application/json")

		tf.rGroup.ServeHTTP(w, req)

		expectedResponseBody, err := json.Marshal(tc.expectedResponse)
		tf.Require().NoError(err)

		tf.Assert().Equal(400, w.Code)
		tf.Assert().Equal(string(expectedResponseBody), w.Body.String())
	}
}
