package handlers

import (
	"net/http"
	"strconv"
	"wallet-manager/models"
	"wallet-manager/services"
	"wallet-manager/utils"

	"github.com/gin-gonic/gin"
)

type CryptoTransactionHandler struct {
	service services.CryptoTransactionService
}

func NewCryptoTransactionHandler(service services.CryptoTransactionService) *CryptoTransactionHandler {
	return &CryptoTransactionHandler{service: service}
}

func (h *CryptoTransactionHandler) Create(c *gin.Context) {
	cryptoId, err := strconv.Atoi(c.Param("cryptoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cryptoId"})
		return
	}

	var crypto models.CryptoTransaction
	if err := c.ShouldBindJSON(&crypto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	crypto.CryptocurrencyId = uint32(cryptoId)
	crypto.CreatedDate = utils.NowFormatted()
	if err := h.service.Create(&crypto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, crypto)
}

func (h *CryptoTransactionHandler) GetAll(c *gin.Context) {
	cryptoId, err := strconv.Atoi(c.Param("cryptoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cryptoId"})
		return
	}

	cryptos, err := h.service.GetAll(uint32(cryptoId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cryptos)
}

func (h *CryptoTransactionHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("transactionId"))
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
		c.JSON(http.StatusNotFound, gin.H{"error": "cryptoTransaction not found"})
		return
	}

	c.JSON(http.StatusOK, crypto)
}

func (h *CryptoTransactionHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("transactionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var crypto models.CryptoTransaction
	if err := c.ShouldBindJSON(&crypto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	crypto.ID = uint32(id)
	if err := h.service.Update(&crypto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, crypto)
}

func (h *CryptoTransactionHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("transactionId"))
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
