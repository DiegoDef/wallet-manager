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
	"wallet-manager/utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
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
	repo    repositories.CryptocurrencyRepository
	service services.CryptocurrencyService
	handle  *handlers.CryptocurrencyHandler
	engine  *gin.Engine
}

func beforeEach() {
	deleteAll()
}

func beforeAll() {
	tc.repo = repositories.NewCryptocurrencyRepository(testDbInstance)
	tc.service = services.NewCryptocurrencyService(tc.repo)
	tc.handle = handlers.NewCryptocurrencyHandler(tc.service)
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

func TestCryptocurrencyService(t *testing.T) {
	t.Run("Should create cryptocurrency", testCase(testCreateCryptocurrency))
	t.Run("Should get all cryptocurrency", testCase(testGetAllCryptocurrencies))
	t.Run("Should find cryptocurrency by ID", testCase(testFindCryptocurrencyById))
	t.Run("Should find all cryptocurrency", testCase(testDeleteCryptocurrency))
	t.Run("Should update cryptocurrency", testCase(testUpdateCryptocurrency))
}

func testCreateCryptocurrency(t *testing.T) {
	tc.engine.POST("/cryptocurrencies", tc.handle.Create)
	server := httptest.NewServer(tc.engine)
	defer server.Close()

	cryptoToSave := createCryptocurrency()

	request, err := http.NewRequest(http.MethodPost, server.URL+"/cryptocurrencies", createCryptocurrencyJson(cryptoToSave))
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	tc.engine.ServeHTTP(responseRecorder, request)

	var crypto models.Cryptocurrency
	err = json.NewDecoder(responseRecorder.Body).Decode(&crypto)
	savedCrypto, errGetById := tc.repo.GetByID(crypto.ID)
	require.NoError(t, err)
	require.NoError(t, errGetById)
	assert.Equal(t, http.StatusCreated, responseRecorder.Result().StatusCode)
	assert.Equal(t, strings.ToLower(cryptoToSave.Name), savedCrypto.Name)
	// assert.Equal(t, cryptoToSave.Balance, savedCrypto.Balance)
	// assert.Equal(t, cryptoToSave.CostInFiat, savedCrypto.CostInFiat)
}

func testUpdateCryptocurrency(t *testing.T) {
	tc.engine.PUT("/cryptocurrencies", tc.handle.Update)
	server := httptest.NewServer(tc.engine)
	defer server.Close()

	toSave := createCryptocurrency()
	err := tc.repo.Create(&toSave)
	require.NoError(t, err)

	cryptoToUpdate := createCryptoWithParamaters("Litecoin", decimal.NewFromInt(10), decimal.NewFromInt(70000))
	cryptoToUpdate.ID = toSave.ID

	request, err := http.NewRequest(http.MethodPut, server.URL+"/cryptocurrencies", createCryptocurrencyJson(cryptoToUpdate))
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	tc.engine.ServeHTTP(responseRecorder, request)

	var crypto models.Cryptocurrency
	err = json.NewDecoder(responseRecorder.Body).Decode(&crypto)
	updatedCrypto, errGetById := tc.repo.GetByID(crypto.ID)
	require.NoError(t, err)
	require.NoError(t, errGetById)
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, strings.ToLower(cryptoToUpdate.Name), updatedCrypto.Name)
	assert.NotEqual(t, toSave.Name, updatedCrypto.Name)
	// assert.Equal(t, toSave.Balance, updatedCrypto.Balance)
	// assert.Equal(t, toSave.CostInFiat, updatedCrypto.CostInFiat)
}

func testDeleteCryptocurrency(t *testing.T) {
	tc.engine.DELETE("/cryptocurrencies/:cryptoId", tc.handle.Delete)

	server := httptest.NewServer(tc.engine)
	defer server.Close()

	toDelete := createCryptocurrency()
	err := tc.repo.Create(&toDelete)
	require.NoError(t, err)

	idToDelete := strconv.FormatUint(uint64(toDelete.ID), 10)
	request, err := http.NewRequest(http.MethodDelete, server.URL+"/cryptocurrencies/"+idToDelete, nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	tc.engine.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusNoContent, responseRecorder.Result().StatusCode)
	exist := true
	testDbInstance.Get(&exist, "select exists(select * from cryptocurrency where cryptocurrency_id = $1)", idToDelete)
	assert.False(t, exist)
}

func testGetAllCryptocurrencies(t *testing.T) {
	crypto := models.Cryptocurrency{Name: "Bitcoin", Balance: decimal.NewFromInt(1), CostInFiat: decimal.NewFromInt(60000), CreatedDate: utils.NowFormatted()}
	tc.repo.Create(&crypto)
	tc.repo.Create(&crypto)
	tc.repo.Create(&crypto)

	tc.engine.GET("/cryptocurrencies", tc.handle.GetAll)

	server := httptest.NewServer(tc.engine)
	defer server.Close()

	request, err := http.NewRequest(http.MethodGet, server.URL+"/cryptocurrencies", nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	tc.engine.ServeHTTP(responseRecorder, request)

	var cryptos []models.Cryptocurrency
	err = json.NewDecoder(responseRecorder.Body).Decode(&cryptos)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, 3, len(cryptos))
}

func testFindCryptocurrencyById(t *testing.T) {
	tc.engine.GET("/cryptocurrencies/:cryptoId", tc.handle.GetByID)

	server := httptest.NewServer(tc.engine)
	defer server.Close()

	toFind := createCryptocurrency()
	err := tc.repo.Create(&toFind)
	require.NoError(t, err)

	idToFind := strconv.FormatUint(uint64(toFind.ID), 10)
	request, err := http.NewRequest(http.MethodGet, server.URL+"/cryptocurrencies/"+idToFind, nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	tc.engine.ServeHTTP(responseRecorder, request)

	var crypto models.Cryptocurrency
	err = json.NewDecoder(responseRecorder.Body).Decode(&crypto)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, toFind.ID, crypto.ID)
	assert.Equal(t, strings.ToLower(toFind.Name), crypto.Name, crypto.Name)
}
