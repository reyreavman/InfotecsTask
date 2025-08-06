package http

import (
	"encoding/json"
	"fmt"
	"infotecstechtask/internal/facade"
	"infotecstechtask/internal/models"
	"infotecstechtask/internal/wallet"
	"infotecstechtask/test/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAllTransaction(t *testing.T) {
	r := gin.Default()
	group := r.Group("")
	dataLoader := testutils.NewDataLoader()
	validate := validator.New()

	var allTransactions []*models.Transaction
	err := dataLoader.LoadJSONFixture("transactions/all_transactions.json", &allTransactions)
	require.NoError(t, err)

	mockFacade := new(facade.MockFacade)

	mockFacade.On(
		"GetAllTransactions",
		mock.Anything,
	).Return(models.ToTransactionResponses(allTransactions), nil)

	RegisterHTTPEndpoints(group, mockFacade, validate)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, FULL_TRANSACTIONS, nil)

	r.ServeHTTP(w, req)

	expectedResponseBody, err := json.Marshal(models.ToTransactionResponses(allTransactions))
	require.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedResponseBody), w.Body.String())
}

func TestGetTransactions(t *testing.T) {
	r := gin.Default()
	group := r.Group("")
	dataLoader := testutils.NewDataLoader()
	validate := validator.New()

	var allTransactions []*models.Transaction
	err := dataLoader.LoadJSONFixture("transactions/all_transactions.json", &allTransactions)
	require.NoError(t, err)
	var oneTransaction []*models.Transaction
	err = dataLoader.LoadJSONFixture("transactions/one_transaction.json", &oneTransaction)
	require.NoError(t, err)
	var twoTransactions []*models.Transaction
	err = dataLoader.LoadJSONFixture("transactions/two_transactions.json", &twoTransactions)
	require.NoError(t, err)

	mockFacade := new(facade.MockFacade)
	RegisterHTTPEndpoints(group, mockFacade, validate)

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

		r.ServeHTTP(w, req)

		expectedResponseBody, err := json.Marshal(models.ToTransactionResponses(tc.expected))
		require.NoError(t, err)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, string(expectedResponseBody), w.Body.String())
		builder.Reset()
	}
}

func TestGetWalletSuccess(t *testing.T) {
	r := gin.Default()
	group := r.Group("")
	dataLoader := testutils.NewDataLoader()

	var expectedResp models.Wallet
	err := dataLoader.LoadJSONFixture("wallets/wallet.json", &expectedResp)
	require.NoError(t, err)

	mockFacade := new(facade.MockFacade)
	mockFacade.On(
		"GetWallet",
		mock.Anything,
		expectedResp.ID,
	).Return(&expectedResp, nil)

	validate := validator.New()

	RegisterHTTPEndpoints(group, mockFacade, validate)

	path := strings.Replace(FULL_GET_WALLET_BALANCE, ":walletId", expectedResp.ID.String(), 1)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, path, nil)

	r.ServeHTTP(w, req)

	expectedResponseBody, err := json.Marshal(expectedResp)
	require.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedResponseBody), w.Body.String())
}

func TestGetWalletNotFound(t *testing.T) {
	r := gin.Default()
	group := r.Group("")
	dataLoader := testutils.NewDataLoader()
	validate := validator.New()

	var expectedResp models.Wallet
	err := dataLoader.LoadJSONFixture("wallets/wallet.json", &expectedResp)
	require.NoError(t, err)

	mockFacade := new(facade.MockFacade)
	mockFacade.On(
		"GetWallet",
		mock.Anything,
		expectedResp.ID,
	).Return(nil, wallet.ErrWalletNotFound)

	RegisterHTTPEndpoints(group, mockFacade, validate)

	path := strings.Replace(FULL_GET_WALLET_BALANCE, ":walletId", expectedResp.ID.String(), 1)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, path, nil)

	r.ServeHTTP(w, req)

	require.NoError(t, err)

	assert.Equal(t, 404, w.Code)
}
