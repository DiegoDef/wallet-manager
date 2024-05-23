package main

import (
	"wallet-manager/config"
	"wallet-manager/db"
	"wallet-manager/handlers"
	"wallet-manager/repositories"
	"wallet-manager/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	database := db.NewDB(&cfg.DB)

	cryptoRepo := repositories.NewCryptocurrencyRepository(database)
	cryptoService := services.NewCryptocurrencyService(cryptoRepo)
	cryptoHandler := handlers.NewCryptocurrencyHandler(cryptoService)

	transactionRepo := repositories.NewCryptoTransactionRepository(database)
	transationService := services.NewCryptoTransactionService(transactionRepo)
	transactionHandler := handlers.NewCryptoTransactionHandler(transationService)

	r := gin.Default()

	r.POST("/cryptocurrencies", cryptoHandler.Create)
	r.GET("/cryptocurrencies", cryptoHandler.GetAll)
	r.GET("/cryptocurrencies/:cryptoId", cryptoHandler.GetByID)
	r.PUT("/cryptocurrencies/:cryptoId", cryptoHandler.Update)
	r.DELETE("/cryptocurrencies/:cryptoId", cryptoHandler.Delete)

	r.POST("/cryptocurrencies/:cryptoId/transactions", transactionHandler.Create)
	r.GET("/cryptocurrencies/:cryptoId/transactions", transactionHandler.GetAll)
	r.GET("/cryptocurrencies/:cryptoId/transactions/:transactionId", transactionHandler.GetByID)
	r.PUT("/cryptocurrencies/:cryptoId/transactions/:transactionId", transactionHandler.Update)
	r.DELETE("/cryptocurrencies/:cryptoId/transactions/:transactionId", transactionHandler.Delete)

	r.Run(":" + cfg.Port)
}
