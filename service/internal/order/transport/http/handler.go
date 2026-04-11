package http

import (
	"AP2_assignment1/service/internal/order/usecase"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	uc *usecase.OrderUsecase
}

func NewOrderHandler(uc *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id"`
		ItemName   string `json:"item_name"`
		Amount     int64  `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.uc.CreateOrder(req.CustomerID, req.ItemName, req.Amount)

	if err != nil {
		if err.Error() == "payment service unavailable" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Payment service is down"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {

	id := c.Param("id")
	order, err := h.uc.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	if err := h.uc.CancelOrder(id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})

}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	minStr := c.Query("min_amount")
	maxStr := c.Query("max_amount")

	minVal, errMin := strconv.ParseInt(minStr, 10, 64)
	maxVal, errMax := strconv.ParseInt(maxStr, 10, 64)

	if errMin != nil || errMax != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "min_amount and max_amount must be valid numbers",
		})
		return
	}

	orders, err := h.uc.GetOrdersByRange(minVal, maxVal)

	if err != nil {
		if strings.Contains(err.Error(), "bad_request") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
