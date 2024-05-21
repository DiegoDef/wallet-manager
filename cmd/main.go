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

	repo := repositories.NewCryptocurrencyRepository(database)
	service := services.NewCryptocurrencyService(repo)
	handler := handlers.NewCryptocurrencyHandler(service)

	r := gin.Default()

	r.POST("/cryptocurrency", handler.Create)
	r.GET("/cryptocurrency", handler.GetAll)
	r.GET("/cryptocurrency/:id", handler.GetByID)
	r.PUT("/cryptocurrency/:id", handler.Update)
	r.DELETE("/cryptocurrency/:id", handler.Delete)

	// repo = repositories.NewCryptocurrencyRepository(database)
	// service = services.NewCryptocurrencyService(repo)
	// handler = handlers.NewCryptocurrencyHandler(service)

	// r.POST("/cryptocurrency/:idCrptcy/transaction", handler.Create)
	// r.GET("/cryptocurrency/:idCrptcy/transaction", handler.GetAll)
	// r.GET("/cryptocurrency/:idCrptcy/transaction/:id", handler.GetByID)
	// r.PUT("/cryptocurrency/:idCrptcy/transaction/:id", handler.Update)
	// r.DELETE("/cryptocurrency/:idCrptcy/transaction/:id", handler.Delete)

	r.Run(":" + cfg.Port)
}
