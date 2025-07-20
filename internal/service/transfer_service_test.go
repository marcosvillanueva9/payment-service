package service_test

import (
	"log"
	"payment-service/internal/middleware/logger"
	"payment-service/internal/model"
	"payment-service/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTransferTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&model.Transfer{}, &model.Account{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	account1 := model.Account{Name: "Test Account 1", Balance: 100.0}
	account2 := model.Account{Name: "Test Account 2", Balance: 200.0}
	db.Create(&account1)
	db.Create(&account2)

	return db
}

func TestCreateTransfer_Success(t *testing.T) {
	db := setupTransferTestDB()
	transferService := service.NewTransferService(db)

	ctx := &gin.Context{}
	logger.Init("test")

	testTransfer := model.TransferRequest{OriginAccountID: 1, DestinationAccountID: 2, Amount: 50.0}
	newTransfer, err := transferService.CreateTransfer(&testTransfer, ctx)

	assert.Nil(t, err)
	assert.Equal(t, float64(50.0), newTransfer.Amount)
	assert.Equal(t, "PENDING", newTransfer.Status)
	assert.Equal(t, uint(1), newTransfer.OriginAccountID)
	assert.Equal(t, uint(2), newTransfer.DestinationAccountID)
}

func TestCreateTransfer_AmountZero(t *testing.T) {
	db := setupTransferTestDB()
	transferService := service.NewTransferService(db)

	ctx := &gin.Context{}
	logger.Init("test")

	testTransfer := model.TransferRequest{OriginAccountID: 1, DestinationAccountID: 2, Amount: 0.0}
	newTransfer, err := transferService.CreateTransfer(&testTransfer, ctx)

	assert.NotNil(t, err)
	assert.Equal(t, "Invalid transfer amount", err.Message)
	assert.Equal(t, model.Transfer{}, newTransfer)
}

func TestCreateTransfer_SameAccount(t *testing.T) {
	db := setupTransferTestDB()
	transferService := service.NewTransferService(db)

	ctx := &gin.Context{}
	logger.Init("test")

	testTransfer := model.TransferRequest{OriginAccountID: 1, DestinationAccountID: 1, Amount: 50.0}
	newTransfer, err := transferService.CreateTransfer(&testTransfer, ctx)

	assert.NotNil(t, err)
	assert.Equal(t, "Cannot transfer to the same account", err.Message)
	assert.Equal(t, model.Transfer{}, newTransfer)
}

func TestCreateTransfer_OriginAccountNotFound(t *testing.T) {
	db := setupTransferTestDB()
	transferService := service.NewTransferService(db)

	ctx := &gin.Context{}
	logger.Init("test")

	testTransfer := model.TransferRequest{OriginAccountID: 999, DestinationAccountID: 2, Amount: 50.0}
	newTransfer, err := transferService.CreateTransfer(&testTransfer, ctx)

	assert.NotNil(t, err)
	assert.Equal(t, "Origin account not found", err.Message)
	assert.Equal(t, model.Transfer{}, newTransfer)
}

func TestCreateTransfer_DestinationAccountNotFound(t *testing.T) {
	db := setupTransferTestDB()
	transferService := service.NewTransferService(db)

	ctx := &gin.Context{}
	logger.Init("test")

	testTransfer := model.TransferRequest{OriginAccountID: 1, DestinationAccountID: 999, Amount: 50.0}
	newTransfer, err := transferService.CreateTransfer(&testTransfer, ctx)

	assert.NotNil(t, err)
	assert.Equal(t, "Destination account not found", err.Message)
	assert.Equal(t, model.Transfer{}, newTransfer)
}