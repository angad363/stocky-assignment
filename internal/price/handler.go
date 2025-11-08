package price

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PriceHandler struct {
	service *PriceService
}

// NewPriceHandler initializes the handler
func NewPriceHandler(service *PriceService) *PriceHandler {
	return &PriceHandler{service: service}
}

// GetPrice handles GET /price?symbol=RELIANCE
func (h *PriceHandler) GetPrice(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}

	price, err := h.service.GetStockPrice(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get price"})
		return
	}

	c.JSON(http.StatusOK, price)
}