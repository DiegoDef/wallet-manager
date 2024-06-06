package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"wallet-manager/handlers"
	"wallet-manager/models"
	"wallet-manager/repositories"
	"wallet-manager/services"
	"wallet-manager/utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDbInstance *sqlx.DB

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func InsertCryptoPrice() {
	testDbInstance.Query("INSERT INTO crypto_price (name, price_usd) VALUES ('Bitcoin', 1.0);")
}

func TestGetAllUsers(t *testing.T) {
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
	// for _, crypto := range cryptos {
	// 	assert.Contains(t, crypto.Name, "radha geethika")
	// 	// repo.GetAll()[0]
	// }
}
