package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"wallet-manager/handlers"
	"wallet-manager/models"
	"wallet-manager/repositories"
	"wallet-manager/services"
	helper "wallet-manager/testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDbInstance *sqlx.DB
var tc testContext

func TestMain(m *testing.M) {
	testDB := helper.SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	beforeAll()
	os.Exit(m.Run())
}

type testContext struct {
	repoCrypto    repositories.CryptocurrencyRepository
	serviceCrypto services.CryptocurrencyService
	repo          repositories.CryptoTransactionRepository
	service       services.CryptoTransactionService
	handle        *handlers.CryptoTransactionHandler
	engine        *gin.Engine
}

func beforeEach() {
	deleteAll()
}

func beforeAll() {
	tc.repoCrypto = repositories.NewCryptocurrencyRepository(testDbInstance)
	tc.serviceCrypto = services.NewCryptocurrencyService(tc.repoCrypto)
	tc.repo = repositories.NewCryptoTransactionRepository(testDbInstance)
	tc.service = services.NewCryptoTransactionService(tc.repo, tc.repoCrypto)
	tc.handle = handlers.NewCryptoTransactionHandler(tc.service)
	tc.engine = gin.Default()
	insertCryptoPrice()
}

func after() {
}

func testCase(test func(t *testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		beforeEach()
		defer after()
		test(t)
	}
}

func insertCryptoPrice() {
	testDbInstance.Query("INSERT INTO crypto_price (name, price_usd) VALUES ('Bitcoin', 1.0);")
}

func deleteAll() {
	testDbInstance.Exec("DELETE FROM cryptocurrency;")
}

// func deleteAllCryptoPrice() {
// 	testDbInstance.Exec("DELETE FROM crypto_price;")
// }

func TestCryptocurrencyService(t *testing.T) {
	t.Run("Should create cryptoTransaction", testCase(testCreateCryptoTransaction))
	t.Run("Should get all cryptoTransaction", testCase(testGetAllCryptoTransaction))
	t.Run("Should find cryptoTransaction by ID", testCase(testFindCryptoTransactionById))
	t.Run("Should delete cryptoTransaction", testCase(testDeleteCryptoTransaction))
	// t.Run("Should update cryptoTransaction", testCase(testUpdatecryptoTransaction))
}

func testCreateCryptoTransaction(t *testing.T) {
	tc.engine.POST("/cryptocurrencies/:cryptoId/transactions", tc.handle.Create)
	server := httptest.NewServer(tc.engine)
	defer server.Close()

	cryptocurrency := createCryptocurrency(testDbInstance)
	transactionToInsert := createTransactionWithoutCryptocurrencyId()
	expectedPurchaseDate := transactionToInsert.PurchaseDate[:strings.LastIndex(transactionToInsert.PurchaseDate, "-")] + "Z"
	expectedCreatedDate := transactionToInsert.CreatedDate[:strings.LastIndex(transactionToInsert.CreatedDate, "-")] + "Z"
	request, err := http.NewRequest(http.MethodPost, server.URL+"/cryptocurrencies/"+strconv.FormatUint(uint64(cryptocurrency.ID), 10)+"/transactions", createCryptoTransactionJson(transactionToInsert))
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	tc.engine.ServeHTTP(responseRecorder, request)

	var transaction models.CryptoTransaction
	err = json.NewDecoder(responseRecorder.Body).Decode(&transaction)
	insertedTransaction, errGetById := tc.repo.GetByID(transaction.ID)

	require.NoError(t, err)
	require.NoError(t, errGetById)
	assert.Equal(t, http.StatusCreated, responseRecorder.Result().StatusCode)
	assert.NotNil(t, transaction.ID)
	assert.Equal(t, cryptocurrency.ID, insertedTransaction.CryptocurrencyId)
	assert.Equal(t, expectedPurchaseDate, insertedTransaction.PurchaseDate)
	assert.Equal(t, expectedCreatedDate, insertedTransaction.CreatedDate)
}

// func testUpdatecryptoTransaction(t *testing.T) {
// 	tc.engine.PUT("/cryptocurrencies/:cryptoId/transactions/:transactionId", tc.handle.Update)
// 	server := httptest.NewServer(tc.engine)
// 	defer server.Close()

// 	toSave := createTransaction(testDbInstance)
// 	toSave.CryptocurrencyId = createCryptocurrency(testDbInstance).ID
// 	err := tc.repo.Create(&toSave)
// 	require.NoError(t, err)

// 	transactionToUpdate := models.CryptoTransaction{
// 		CryptocurrencyId:     createCryptocurrency(testDbInstance).ID,
// 		CryptocurrencyAmount: decimal.NewFromInt(130),
// 		FiatAmount:           decimal.NewFromInt(150),
// 		PurchaseDate:         utils.NowFormatted(),
// 		CreatedDate:          utils.NowFormatted(),
// 	}
// 	transactionToUpdate.ID = toSave.ID

// 	request, err := http.NewRequest(http.MethodPut, server.URL+"/cryptocurrencies"+strconv.FormatUint(uint64(toSave.CryptocurrencyId), 10)+"/transactions/"+strconv.FormatUint(uint64(toSave.ID), 10), createCryptoTransactionJson(transactionToUpdate))
// 	require.NoError(t, err)

// 	responseRecorder := httptest.NewRecorder()
// 	tc.engine.ServeHTTP(responseRecorder, request)

// 	var transaction models.CryptoTransaction
// 	err = json.NewDecoder(responseRecorder.Body).Decode(&transaction)
// 	updatedtransaction, errGetById := tc.repo.GetByID(toSave.ID)
// 	require.NoError(t, err)
// 	require.NoError(t, errGetById)
// 	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
// 	assert.Equal(t, transactionToUpdate.CryptocurrencyId, updatedtransaction.CryptocurrencyId)
// 	assert.Equal(t, transactionToUpdate.CryptocurrencyAmount, updatedtransaction.CryptocurrencyAmount)
// 	assert.Equal(t, transactionToUpdate, updatedtransaction.PurchaseDate)
// 	assert.Equal(t, transactionToUpdate, updatedtransaction.CreatedDate)
// }

func testDeleteCryptoTransaction(t *testing.T) {
	tc.engine.DELETE("/cryptocurrencies/:cryptoId/transactions/:transactionId", tc.handle.Delete)

	server := httptest.NewServer(tc.engine)
	defer server.Close()

	toDelete := createTransaction(testDbInstance)
	err := tc.repo.Create(&toDelete)
	require.NoError(t, err)

	idToDelete := strconv.FormatUint(uint64(toDelete.ID), 10)
	request, err := http.NewRequest(http.MethodDelete, server.URL+"/cryptocurrencies/"+strconv.FormatUint(uint64(toDelete.CryptocurrencyId), 10)+"/transactions/"+idToDelete, nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	tc.engine.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusNoContent, responseRecorder.Result().StatusCode)
	exist := true
	testDbInstance.Get(&exist, "select exists(select * from cryptocurrency where cryptocurrency_id = $1)", idToDelete)
	assert.False(t, exist)
}

func testGetAllCryptoTransaction(t *testing.T) {
	transaction := createTransaction(testDbInstance)
	tc.repo.Create(&transaction)
	tc.repo.Create(&transaction)
	tc.repo.Create(&transaction)

	tc.engine.GET("cryptocurrencies/:cryptoId/transactions", tc.handle.GetAll)

	server := httptest.NewServer(tc.engine)
	defer server.Close()

	request, err := http.NewRequest(http.MethodGet, server.URL+"/cryptocurrencies/"+strconv.FormatUint(uint64(transaction.CryptocurrencyId), 10)+"/transactions", nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	tc.engine.ServeHTTP(responseRecorder, request)

	var cryptos []models.CryptoTransaction
	err = json.NewDecoder(responseRecorder.Body).Decode(&cryptos)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, 3, len(cryptos))
}

func testFindCryptoTransactionById(t *testing.T) {
	tc.engine.GET("cryptocurrencies/:cryptoId/transactions/:transactionId", tc.handle.GetByID)

	server := httptest.NewServer(tc.engine)
	defer server.Close()

	toFind := createTransaction(testDbInstance)
	err := tc.repo.Create(&toFind)
	require.NoError(t, err)
	expectedPurchaseDate := toFind.PurchaseDate[:strings.LastIndex(toFind.PurchaseDate, "-")] + "Z"
	expectedCreatedDate := toFind.CreatedDate[:strings.LastIndex(toFind.CreatedDate, "-")] + "Z"

	idToFind := strconv.FormatUint(uint64(toFind.ID), 10)
	request, err := http.NewRequest(http.MethodGet, server.URL+"/cryptocurrencies/"+strconv.FormatUint(uint64(toFind.CryptocurrencyId), 10)+"/transactions/"+idToFind, nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	tc.engine.ServeHTTP(responseRecorder, request)

	var transaction models.CryptoTransaction
	err = json.NewDecoder(responseRecorder.Body).Decode(&transaction)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, toFind.ID, transaction.ID)
	assert.Equal(t, toFind.CryptocurrencyId, transaction.CryptocurrencyId)
	assert.Equal(t, expectedPurchaseDate, transaction.PurchaseDate)
	assert.Equal(t, expectedCreatedDate, transaction.CreatedDate)
}
