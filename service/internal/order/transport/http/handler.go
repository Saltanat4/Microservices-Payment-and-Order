package http

import (
	"AP2_assignment1/service/internal/order/usecase"
	"database/sql"
	"log"
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
		CustomerID    string `json:"customer_id"`
		ItemName      string `json:"item_name"`
		Amount        int64  `json:"amount"`
		CustomerEmail string `json:"customer_email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.uc.CreateOrder(req.CustomerID, req.ItemName, req.Amount, req.CustomerEmail)
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "connection confused") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Payment service is down"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
	log.Printf("Received order request for amount: %d", req.Amount)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.uc.GetOrder(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	err := h.uc.CancelOrder(id)
	if err != nil {
		switch err.Error() {
		case "cannot cancel order in Paid status":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		case "sql: no rows in result set":
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
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
