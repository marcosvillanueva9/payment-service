package controller

import (
	"net/http"
	"payment-service/internal/middleware/logger"
	"payment-service/internal/model"
	"payment-service/internal/service"

	"github.com/gin-gonic/gin"
)

type TransferController interface {
	CreateTransfer(c *gin.Context)
	UpdateStatus(c *gin.Context)
}

type transferController struct {
	service service.TransferService
}

func NewTransferController(service service.TransferService) TransferController {
	return &transferController{
		service: service,
	}
}

func (ctrl *transferController) CreateTransfer(c *gin.Context) {
	logger := logger.From(c)
	logger.Infow("Creating transfer")

	var transferRequest model.TransferRequest

	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		logger.Errorw("Invalid transfer request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	transfer, err := ctrl.service.CreateTransfer(&transferRequest, c)
	if err != nil {
		if err.Code == http.StatusNotFound {
			logger.Warnw("Transfer failed due to account not found", "error", err)
			c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
			return
		} else {
			logger.Errorw("Transfer failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Transfer failed", "error": err.Message})
			return
		}
	}

	logger.Infow("Transfer created successfully", "transfer_id", transfer.ID, "amount", transfer.Amount)
	c.JSON(http.StatusCreated, gin.H{"message": "Transfer successful", "transfer": transfer})
}

func (ctrl *transferController) UpdateStatus(c *gin.Context) {
	logger := logger.From(c)
	logger.Infow("Updating transfer status")

	transferID := c.Param("id")
	if transferID == "" {
		logger.Errorw("Transfer ID is required for webhook", "error", "missing transfer ID")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Transfer ID is required"})
		return
	}

	var req model.TransferUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorw("Invalid update request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	transfer, err := ctrl.service.UpdateTransferStatus(transferID, req.Status, c)
	if err != nil {
		if err.Code == http.StatusNotFound {
			logger.Warnw("Transfer not found", "transfer_id", transferID, "error", err)
			c.JSON(http.StatusNotFound, gin.H{"message": "Transfer not found"})
			return
		} else {
			logger.Errorw("Failed to update transfer status", "transfer_id", transferID, "error", err)
			c.JSON(err.Code, gin.H{"message": "Failed to update transfer status", "error": err.Message})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully", "transfer": transfer})
}