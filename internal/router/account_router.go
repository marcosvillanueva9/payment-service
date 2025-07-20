package router

import (
	"payment-service/internal/controller"
	"payment-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AccountRouter(r *gin.RouterGroup, db *gorm.DB) {
	accountController := controller.NewAccountController(service.NewAccountService(db))
	r.GET("/:id/balance", accountController.GetAccountBalance)
}