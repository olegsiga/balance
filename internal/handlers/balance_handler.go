package handlers

import (
	"net/http"

	"balance/internal/models"
	"balance/internal/service"
	"github.com/gin-gonic/gin"
)

type BalanceHandler struct {
	balanceService *service.BalanceService
}

func NewBalanceHandler(balanceService *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{balanceService: balanceService}
}

func (h *BalanceHandler) GetBalance(c *gin.Context) {
	userIDStr := c.Param("userId")

	userID, err := h.balanceService.ValidateUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.balanceService.GetUserBalance(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (h *BalanceHandler) ProcessTransaction(c *gin.Context) {
	userIDStr := c.Param("userId")
	sourceType := c.GetHeader("Source-Type")

	if sourceType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Source-Type header is required"})
		return
	}

	if sourceType != "game" && sourceType != "server" && sourceType != "payment" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Source-Type. Must be one of: game, server, payment"})
		return
	}

	userID, err := h.balanceService.ValidateUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.balanceService.ProcessTransaction(userID, &req, sourceType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.TransactionResponse{
		Success: true,
		Message: "Transaction processed successfully",
	})
}
