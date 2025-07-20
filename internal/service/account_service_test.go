package service_test

import (
	"fmt"
	"log"
	"net/http"
	"payment-service/internal/middleware/logger"
	"payment-service/internal/model"
	"payment-service/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAccountTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&model.Account{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestGetAccountBalance_Success(t *testing.T) {
	println("Running TestGetAccountBalance_Success")
	db := setupAccountTestDB()
	accountService := service.NewAccountService(db)

	testAccount := model.Account{Name: "Test Account", Balance: 100.0}
	db.Create(&testAccount)

	// create a mock context
	ctx := &gin.Context{}
	logger.Init("test")

	response, err := accountService.GetAccountBalance(fmt.Sprint(testAccount.ID), ctx)

	assert.Nil(t, err)
	assert.Equal(t, float64(100.0), response.Balance)
}

func TestGetAccountBalance_AccountNotFound(t *testing.T) {
	println("Running TestGetAccountBalance_AccountNotFound")
	db := setupAccountTestDB()
	accountService := service.NewAccountService(db)

	// create a mock context
	ctx := &gin.Context{}
	logger.Init("test")

	response, err := accountService.GetAccountBalance("nonexistent", ctx)

	assert.NotNil(t, err)
	assert.Equal(t, "Account not found", err.Message)
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Equal(t, model.AccountBalanceResponse{}, response)
}