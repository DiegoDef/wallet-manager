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
	// t.Run("Should get all cryptocurrency", testCase(testGetAllCryptocurrencies))
	// t.Run("Should find cryptocurrency by ID", testCase(testFindCryptocurrencyById))
	// t.Run("Should find all cryptocurrency", testCase(testDeleteCryptocurrency))
	// t.Run("Should update cryptocurrency", testCase(testUpdateCryptocurrency))
	// t.Run("Should find cryptocurrency when there is no crypto price for crypto name", testCase(testFindCryptocurrencyWithoutCryptoPrice))
	// t.Run("Should return cryptocurrency with profitPercentage", testCase(testFindCryptocurrencyWhihoutCryptoPrice))
}

func testCreateCryptoTransaction(t *testing.T) {
	tc.engine.POST("/cryptocurrencies/:cryptoId/transactions", tc.handle.Create)
	server := httptest.NewServer(tc.engine)
	defer server.Close()

	cryptocurrency := insertCryptocurrency(testDbInstance)
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

// func testUpdateCryptocurrency(t *testing.T) {
// 	tc.engine.PUT("/cryptocurrencies", tc.handle.Update)
// 	server := httptest.NewServer(tc.engine)
// 	defer server.Close()

// 	toSave := createCryptocurrency()
// 	err := tc.repo.Create(&toSave)
// 	require.NoError(t, err)

// 	cryptoToUpdate := createCryptoWithParamaters("Litecoin", decimal.NewFromInt(10), decimal.NewFromInt(70000))
// 	cryptoToUpdate.ID = toSave.ID

// 	request, err := http.NewRequest(http.MethodPut, server.URL+"/cryptocurrencies", createCryptocurrencyJson(cryptoToUpdate))
// 	require.NoError(t, err)

// 	responseRecorder := httptest.NewRecorder()
// 	tc.engine.ServeHTTP(responseRecorder, request)

// 	var crypto models.Cryptocurrency
// 	err = json.NewDecoder(responseRecorder.Body).Decode(&crypto)
// 	updatedCrypto, errGetById := tc.repo.GetByID(crypto.ID)
// 	require.NoError(t, err)
// 	require.NoError(t, errGetById)
// 	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
// 	assert.Equal(t, strings.ToLower(cryptoToUpdate.Name), updatedCrypto.Name)
// 	assert.NotEqual(t, toSave.Name, updatedCrypto.Name)
// 	// assert.Equal(t, toSave.Balance, updatedCrypto.Balance)
// 	// assert.Equal(t, toSave.CostInFiat, updatedCrypto.CostInFiat)
// }

// func testDeleteCryptocurrency(t *testing.T) {
// 	tc.engine.DELETE("/cryptocurrencies/:cryptoId", tc.handle.Delete)

// 	server := httptest.NewServer(tc.engine)
// 	defer server.Close()

// 	toDelete := createCryptocurrency()
// 	err := tc.repo.Create(&toDelete)
// 	require.NoError(t, err)

// 	idToDelete := strconv.FormatUint(uint64(toDelete.ID), 10)
// 	request, err := http.NewRequest(http.MethodDelete, server.URL+"/cryptocurrencies/"+idToDelete, nil)
// 	require.NoError(t, err)

// 	responseRecorder := httptest.NewRecorder()
// 	tc.engine.ServeHTTP(responseRecorder, request)

// 	assert.Equal(t, http.StatusNoContent, responseRecorder.Result().StatusCode)
// 	exist := true
// 	testDbInstance.Get(&exist, "select exists(select * from cryptocurrency where cryptocurrency_id = $1)", idToDelete)
// 	assert.False(t, exist)
// }

// func testGetAllCryptocurrencies(t *testing.T) {
// 	crypto := models.Cryptocurrency{Name: "Bitcoin", Balance: decimal.NewFromInt(1), CostInFiat: decimal.NewFromInt(60000), CreatedDate: utils.NowFormatted()}
// 	tc.repo.Create(&crypto)
// 	tc.repo.Create(&crypto)
// 	tc.repo.Create(&crypto)

// 	tc.engine.GET("/cryptocurrencies", tc.handle.GetAll)

// 	server := httptest.NewServer(tc.engine)
// 	defer server.Close()

// 	request, err := http.NewRequest(http.MethodGet, server.URL+"/cryptocurrencies", nil)
// 	require.NoError(t, err)

// 	responseRecorder := httptest.NewRecorder()
// 	tc.engine.ServeHTTP(responseRecorder, request)

// 	var cryptos []models.Cryptocurrency
// 	err = json.NewDecoder(responseRecorder.Body).Decode(&cryptos)
// 	require.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
// 	assert.Equal(t, 3, len(cryptos))
// }

// func testFindCryptocurrencyById(t *testing.T) {
// 	tc.engine.GET("/cryptocurrencies/:cryptoId", tc.handle.GetByID)

// 	server := httptest.NewServer(tc.engine)
// 	defer server.Close()

// 	toFind := createCryptocurrency()
// 	err := tc.repo.Create(&toFind)
// 	require.NoError(t, err)

// 	idToFind := strconv.FormatUint(uint64(toFind.ID), 10)
// 	request, err := http.NewRequest(http.MethodGet, server.URL+"/cryptocurrencies/"+idToFind, nil)
// 	require.NoError(t, err)

// 	responseRecorder := httptest.NewRecorder()
// 	tc.engine.ServeHTTP(responseRecorder, request)

// 	var crypto models.Cryptocurrency
// 	err = json.NewDecoder(responseRecorder.Body).Decode(&crypto)
// 	require.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
// 	assert.Equal(t, toFind.ID, crypto.ID)
// 	assert.Equal(t, strings.ToLower(toFind.Name), crypto.Name, crypto.Name)
// }

// func testFindCryptocurrencyWithoutCryptoPrice(t *testing.T) {
// 	server := httptest.NewServer(tc.engine)
// 	defer server.Close()

// 	deleteAllCryptoPrice()
// 	toFind := createCryptocurrency()
// 	toFind.Name = "CryptoWihtoutPrice"
// 	err := tc.repo.Create(&toFind)
// 	require.NoError(t, err)

// 	idToFind := strconv.FormatUint(uint64(toFind.ID), 10)
// 	request, err := http.NewRequest(http.MethodGet, server.URL+"/cryptocurrencies/"+idToFind, nil)
// 	require.NoError(t, err)

// 	responseRecorder := httptest.NewRecorder()
// 	tc.engine.ServeHTTP(responseRecorder, request)

// 	var crypto models.Cryptocurrency
// 	err = json.NewDecoder(responseRecorder.Body).Decode(&crypto)
// 	require.NoError(t, err)

// 	existCryptoPrice := true
// 	testDbInstance.Get(&existCryptoPrice, "select exists(select * from crypto_price)")

// 	assert.False(t, existCryptoPrice)
// 	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
// 	assert.Equal(t, toFind.ID, crypto.ID)
// 	assert.Equal(t, strings.ToLower(toFind.Name), crypto.Name, crypto.Name)
// }
