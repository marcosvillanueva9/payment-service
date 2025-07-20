package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"payment-service/internal/controller"
	"payment-service/internal/service"
)

func TransferRouter(r *gin.RouterGroup, db *gorm.DB) {
	transferService := service.NewTransferService(db)
	transferController := controller.NewTransferController(transferService)

	r.POST("/", transferController.CreateTransfer)
	r.POST("/:id/webhook", transferController.UpdateStatus)
}