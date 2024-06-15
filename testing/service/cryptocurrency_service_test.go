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

func TestMain(m *testing.M) {
	testDB := helper.SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func InsertCryptoPrice() {
	testDbInstance.Query("INSERT INTO crypto_price (name, price_usd) VALUES ('Bitcoin', 1.0);")
}

func TestGetAllCryptocurrencies(t *testing.T) {
	repo := repositories.NewCryptocurrencyRepository(testDbInstance)
	service := services.NewCryptocurrencyService(repo)
	h := handlers.NewCryptocurrencyHandler(service)

	InsertCryptoPrice()
	c := models.Cryptocurrency{Name: "Bitcoin", Balance: decimal.NewFromInt(1), CostInFiat: decimal.NewFromInt(60000), CreatedDate: utils.NowFormatted()}
	repo.Create(&c)
	repo.Create(&c)
	repo.Create(&c)

	engine := gin.Default()
	engine.GET("/cryptocurrencies", h.GetAll)

	server := httptest.NewServer(engine)
	defer server.Close()

	request, err := http.NewRequest(http.MethodGet, server.URL+"/cryptocurrencies", nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(responseRecorder, request)

	var cryptos []models.Cryptocurrency
	err = json.NewDecoder(responseRecorder.Body).Decode(&cryptos)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, 3, len(cryptos))
}

func TestCreateCryptocurrency(t *testing.T) {
	repo := repositories.NewCryptocurrencyRepository(testDbInstance)
	service := services.NewCryptocurrencyService(repo)
	h := handlers.NewCryptocurrencyHandler(service)

	engine := gin.Default()
	engine.POST("/cryptocurrencies", h.Create)

	server := httptest.NewServer(engine)
	defer server.Close()

	cryptoToSave := createCryptocurrency()

	request, err := http.NewRequest(http.MethodPost, server.URL+"/cryptocurrencies", createCryptocurrencyJson(cryptoToSave))
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(responseRecorder, request)

	var crypto models.Cryptocurrency
	err = json.NewDecoder(responseRecorder.Body).Decode(&crypto)
	savedCrypto, errGetById := repo.GetByID(crypto.ID)
	require.NoError(t, err)
	require.NoError(t, errGetById)
	assert.Equal(t, http.StatusCreated, responseRecorder.Result().StatusCode)
	assert.Equal(t, strings.ToLower(cryptoToSave.Name), savedCrypto.Name)
	// assert.Equal(t, cryptoToSave.Balance, savedCrypto.Balance)
	// assert.Equal(t, cryptoToSave.CostInFiat, savedCrypto.CostInFiat)
}

func TestUpdateCryptocurrency(t *testing.T) {
	repo := repositories.NewCryptocurrencyRepository(testDbInstance)
	service := services.NewCryptocurrencyService(repo)
	h := handlers.NewCryptocurrencyHandler(service)

	engine := gin.Default()
	engine.PUT("/cryptocurrencies", h.Update)

	server := httptest.NewServer(engine)
	defer server.Close()

	toSave := createCryptocurrency()
	err := repo.Create(&toSave)
	require.NoError(t, err)

	cryptoToUpdate := createCryptoWithParamaters("Litecoin", decimal.NewFromInt(10), decimal.NewFromInt(70000))
	cryptoToUpdate.ID = toSave.ID

	request, err := http.NewRequest(http.MethodPut, server.URL+"/cryptocurrencies", createCryptocurrencyJson(cryptoToUpdate))
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(responseRecorder, request)

	var crypto models.Cryptocurrency
	err = json.NewDecoder(responseRecorder.Body).Decode(&crypto)
	updatedCrypto, errGetById := repo.GetByID(crypto.ID)
	require.NoError(t, err)
	require.NoError(t, errGetById)
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, strings.ToLower(cryptoToUpdate.Name), updatedCrypto.Name)
	assert.NotEqual(t, toSave.Name, updatedCrypto.Name)
	// assert.Equal(t, toSave.Balance, updatedCrypto.Balance)
	// assert.Equal(t, toSave.CostInFiat, updatedCrypto.CostInFiat)
}

func TestDeleteCryptocurrency(t *testing.T) {
	repo := repositories.NewCryptocurrencyRepository(testDbInstance)
	service := services.NewCryptocurrencyService(repo)
	h := handlers.NewCryptocurrencyHandler(service)

	engine := gin.Default()
	engine.DELETE("/cryptocurrencies/:cryptoId", h.Delete)

	server := httptest.NewServer(engine)
	defer server.Close()

	toDelete := createCryptocurrency()
	err := repo.Create(&toDelete)
	require.NoError(t, err)

	idToDelete := strconv.FormatUint(uint64(toDelete.ID), 10)
	request, err := http.NewRequest(http.MethodDelete, server.URL+"/cryptocurrencies/"+idToDelete, nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusNoContent, responseRecorder.Result().StatusCode)
	exist := true
	testDbInstance.Get(&exist, "select exists(select * from cryptocurrency where cryptocurrency_id = $1)", idToDelete)
	assert.False(t, exist)
}

func TestFindByIdCryptocurrency(t *testing.T) {
	repo := repositories.NewCryptocurrencyRepository(testDbInstance)
	service := services.NewCryptocurrencyService(repo)
	h := handlers.NewCryptocurrencyHandler(service)

	engine := gin.Default()
	engine.GET("/cryptocurrencies/:cryptoId", h.GetByID)

	server := httptest.NewServer(engine)
	defer server.Close()

	toFind := createCryptocurrency()
	err := repo.Create(&toFind)
	require.NoError(t, err)

	idToFind := strconv.FormatUint(uint64(toFind.ID), 10)
	request, err := http.NewRequest(http.MethodGet, server.URL+"/cryptocurrencies/"+idToFind, nil)
	require.NoError(t, err)

	responseRecorder := httptest.NewRecorder()
	engine.ServeHTTP(responseRecorder, request)

	var crypto models.Cryptocurrency
	err = json.NewDecoder(responseRecorder.Body).Decode(&crypto)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
	assert.Equal(t, toFind.ID, crypto.ID)
	assert.Equal(t, strings.ToLower(toFind.Name), crypto.Name, crypto.Name)
}
