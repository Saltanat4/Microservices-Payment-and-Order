package http

import (
	"AP2_assignment1/service/internal/payment/usecase"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	uc *usecase.PaymentUsecase
}

func NewPaymentHandler(uc *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id"`
		Amount  int64  `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}

	result, err := h.uc.Process(req.OrderID, req.Amount)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	orderID := c.Param("order_id")
	result, err := h.uc.GetStatus(orderID)
	if err != nil {
		c.JSON(404, gin.H{"error": "payment not found"})
		return
	}
	c.JSON(200, result)
}
