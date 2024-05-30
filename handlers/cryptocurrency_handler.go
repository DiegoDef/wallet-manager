package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"wallet-manager/models"
	"wallet-manager/services"
	"wallet-manager/utils"

	"github.com/gin-gonic/gin"
)

type CryptocurrencyHandler struct {
	service services.CryptocurrencyService
}

func NewCryptocurrencyHandler(service services.CryptocurrencyService) *CryptocurrencyHandler {
	return &CryptocurrencyHandler{service: service}
}

func (h *CryptocurrencyHandler) Create(c *gin.Context) {
	var crypto models.Cryptocurrency
	if err := c.ShouldBindJSON(&crypto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	crypto.CreatedDate = utils.NowFormatted()
	if err := h.service.Create(&crypto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, crypto)
}

func (h *CryptocurrencyHandler) GetAll(c *gin.Context) {
	cryptos, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cryptos)
}

func (h *CryptocurrencyHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("cryptoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	crypto, err := h.service.GetByID(uint32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if crypto == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cryptocurrency not found"})
		return
	}

	c.JSON(http.StatusOK, crypto)
}

func (h *CryptocurrencyHandler) Update(c *gin.Context) {
	var crypto models.Cryptocurrency
	if err := c.ShouldBindJSON(&crypto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(&crypto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, crypto)
}

func (h *CryptocurrencyHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("cryptoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.service.Delete(uint32(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CryptocurrencyHandler) GetMultiplePrices(c *gin.Context) {
	namesParam := c.Query("names")
	if namesParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no names provided"})
		return
	}

	names := strings.Split(namesParam, ",")
	prices, err := utils.GetCryptoPrices(names)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prices)
}
