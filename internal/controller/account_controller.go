package controller

import (
	"net/http"
	"payment-service/internal/middleware/logger"
	"payment-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AccountController interface {
	GetAccountBalance(c *gin.Context)
}

type accountController struct {
	service service.AccountService
}

func NewAccountController(service service.AccountService) AccountController {
	return &accountController{
		service: service,
	}
}

func (ctrl *accountController) GetAccountBalance(c *gin.Context) {
	log := logger.From(c)
	accountID := c.Param("id")
	if accountID == "" {
		log.Errorw("Invalid account ID", "error", "Account ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	balance, err := ctrl.service.GetAccountBalance(accountID)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"result":balance})
}