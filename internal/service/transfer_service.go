package service

import (
	"log"
	"net/http"
	"payment-service/internal/constant"
	"payment-service/internal/middleware/logger"
	"payment-service/internal/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransferService interface {
	CreateTransfer(req *model.TransferRequest, ctx *gin.Context) (model.Transfer, *ServiceError)
	UpdateTransferStatus(transferID string, status string, ctx *gin.Context) (*model.Transfer, *ServiceError)
	CronExpireTransfers() (*ServiceError)
}

type transferService struct {
	db *gorm.DB
}

type ServiceError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   error  `json:"-"`
}

func NewTransferService(db *gorm.DB) TransferService {
	return &transferService{db: db}
}

func (s *transferService) CreateTransfer(req *model.TransferRequest, ctx *gin.Context) (model.Transfer, *ServiceError) {
	log := logger.From(ctx)

	if req.Amount <= 0 {
		log.Errorw("Transfer failed: Amount must be greater than zero")
		return model.Transfer{}, &ServiceError{Message: "Invalid transfer amount", Code: http.StatusBadRequest}
	}

	if req.OriginAccountID == req.DestinationAccountID {
		log.Errorw("Transfer failed: From and To account IDs are the same")
		return model.Transfer{}, &ServiceError{Message: "Cannot transfer to the same account", Code: http.StatusBadRequest}
	}

	var originAccount, destinationAccount model.Account
	if err := s.db.First(&originAccount, req.OriginAccountID).Error; err != nil {
		log.Errorw("Transfer failed: Origin account not found", "error", err)
		return model.Transfer{}, &ServiceError{Message: "Origin account not found", Code: http.StatusNotFound}
	}

	if err := s.db.First(&destinationAccount, req.DestinationAccountID).Error; err != nil {
		log.Errorw("Transfer failed: Destination account not found", "error", err)
		return model.Transfer{}, &ServiceError{Message: "Destination account not found", Code: http.StatusNotFound}
	}

	transfer := model.Transfer{
		OriginAccountID:      req.OriginAccountID,
		DestinationAccountID: req.DestinationAccountID,
		Amount:               req.Amount,
		Status:               constant.TransferStatusPending,
	}

	if err := s.db.Create(&transfer).Error; err != nil {
		log.Errorw("Transfer failed: Unable to create transfer", "error", err)
		return model.Transfer{}, &ServiceError{Message: "Unable to create transfer", Code: http.StatusInternalServerError}
	}

	log.Infow("Transfer created successfully", "transfer", transfer)

	return transfer, nil
}

func (s *transferService) UpdateTransferStatus(transferID string, status string, ctx *gin.Context) (*model.Transfer, *ServiceError) {
	log := logger.From(ctx)

	var transfer model.Transfer
	if err := s.db.First(&transfer, transferID).Error; err != nil {
		log.Errorw("Transfer not found", "transfer_id", transferID, "error", err)
		return nil, &ServiceError{Message: "Transfer not found", Code: http.StatusNotFound}
	}

	if status != constant.TransferStatusPending && status != constant.TransferStatusCompleted && status != constant.TransferStatusFailed {
		log.Errorw("Invalid transfer status", "status", status)
		return nil, &ServiceError{Message: "Invalid transfer status", Code: http.StatusBadRequest}
	}

	if transfer.Status != constant.TransferStatusPending {
		log.Errorw("Transfer can only be updated if it is pending", "transfer_id", transferID, "current_status", transfer.Status)
		return nil, &ServiceError{Message: "Transfer can only be updated if it is pending", Code: http.StatusBadRequest}
	}

	if status == constant.TransferStatusFailed {
		transfer.Status = constant.TransferStatusFailed
		s.db.Save(&transfer)
		log.Infow("Transfer marked as failed", "transfer_id", transfer.ID)
		return &transfer, nil
	}



	log.Infow("Transfer status updated successfully", "transfer_id", transfer.ID, "status", transfer.Status)

	return s.completeTransfer(&transfer, ctx)
}

func (s *transferService) CronExpireTransfers() (*ServiceError) {

	timeLimit := time.Now().Add(-5 * time.Minute)
	result := s.db.Model(&model.Transfer{}).
		Where("status = ? AND updated_at < ?", constant.TransferStatusPending, timeLimit).
		Update("status", constant.TransferStatusFailed)
	
	if result.Error != nil {
		return &ServiceError{Message: "Failed to expire transfers", Code: http.StatusInternalServerError, Error: result.Error}
	}

	if result.RowsAffected > 0 {
		log.Println("[Cron] Expired transfers", "count", result.RowsAffected)
	}

	return nil
}

func (s *transferService) completeTransfer(transfer *model.Transfer, ctx *gin.Context) (*model.Transfer, *ServiceError) {
	log := logger.From(ctx)

	tx := s.db.Begin()

	var originAccount, destinationAccount model.Account
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&originAccount, "id = ?", transfer.OriginAccountID).Error; err != nil {
		tx.Rollback()
		log.Errorw("Transfer failed: Origin account not found", "error", err)
		return nil, &ServiceError{Message: "Origin account not found", Code: http.StatusNotFound}
	}

	if originAccount.Balance < transfer.Amount {
		tx.Rollback()
		log.Errorw("Transfer failed: Insufficient funds", "transfer_id", transfer.ID)
		return nil, &ServiceError{Message: "Insufficient funds", Code: http.StatusBadRequest}
	}

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&destinationAccount, "id = ?", transfer.DestinationAccountID).Error; err != nil {
		tx.Rollback()
		log.Errorw("Transfer failed: Destination account not found", "error", err)
		return nil, &ServiceError{Message: "Destination account not found", Code: http.StatusNotFound}
	}

	originAccount.Balance -= transfer.Amount
	destinationAccount.Balance += transfer.Amount
	transfer.Status = constant.TransferStatusCompleted

	if err := tx.Save(&originAccount).Error; err != nil {
		tx.Rollback()
		log.Errorw("Transfer failed: Unable to update origin account", "error", err)
		return nil, &ServiceError{Message: "Unable to update origin account", Code: http.StatusInternalServerError}
	}

	if err := tx.Save(&destinationAccount).Error; err != nil {
		tx.Rollback()
		log.Errorw("Transfer failed: Unable to update destination account", "error", err)
		return nil, &ServiceError{Message: "Unable to update destination account", Code: http.StatusInternalServerError}
	}

	if err := tx.Save(&transfer).Error; err != nil {
		tx.Rollback()
		log.Errorw("Transfer failed: Unable to update transfer", "error", err)
		return nil, &ServiceError{Message: "Unable to update transfer", Code: http.StatusInternalServerError}
	}

	tx.Commit()
	log.Infow("Transfer completed successfully", "transfer_id", transfer.ID)
	return transfer, nil
}