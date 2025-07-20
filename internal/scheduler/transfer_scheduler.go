package scheduler

import (
	"log"
	"payment-service/internal/service"

	"github.com/robfig/cron/v3"
)

type TransferScheduler interface {
	Start()
	Stop()
}

type transferScheduler struct {
	cron    *cron.Cron
	service service.TransferService
}

func NewTransferScheduler(service service.TransferService) TransferScheduler {
	return &transferScheduler{
		cron:    cron.New(cron.WithSeconds()),
		service: service,
	}
}

func (ts *transferScheduler) Start() {
	// Schedule the CronExpireTransfers method to run every minute
	_, err := ts.cron.AddFunc("@every 1m", func() {
		if err := ts.service.CronExpireTransfers(); err != nil {
			log.Println("[CRON] Error expiring transfers:", err)
		}
	})
	if err != nil {
		log.Fatalf("[CRON] Failed to schedule transfer expiration: %v", err)
	}

	ts.cron.Start()
	log.Println("[CRON] Transfer scheduler started, will expire transfers every minute")
}

func (ts *transferScheduler) Stop() {
	ts.cron.Stop()
	log.Println("[CRON] Transfer scheduler stopped")
}