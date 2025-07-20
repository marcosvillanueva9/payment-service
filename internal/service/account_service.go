package service

import (
	"net/http"
	"payment-service/internal/model"

	"gorm.io/gorm"
)

import (
	"github.com/gin-gonic/gin"
)

type AccountService interface {
	GetAccountBalance(accountID string, ctx *gin.Context) (model.AccountBalanceResponse, *ServiceError)
}

type accountService struct {
	db *gorm.DB
}
func NewAccountService(db *gorm.DB) AccountService {
	return &accountService{db: db}
}

func (s *accountService) GetAccountBalance(accountID string, ctx *gin.Context) (model.AccountBalanceResponse, *ServiceError) {
	var account model.Account
	if err := s.db.First(&account, "id = ?", accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.AccountBalanceResponse{}, &ServiceError{Message: "Account not found", Code: http.StatusNotFound}
		}
		return model.AccountBalanceResponse{}, &ServiceError{Message: "Failed to retrieve account balance", Code: http.StatusInternalServerError, Error: err}
	}
	return model.AccountBalanceResponse{AccountID: account.ID, Balance: account.Balance}, nil
}