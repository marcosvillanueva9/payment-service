package router

import (
	"net/http"
	"payment-service/internal/middleware/logger"
	"payment-service/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AccountRouter(r *gin.RouterGroup, db *gorm.DB) {
	r.GET("/:id/balance", func(c *gin.Context) {
		log := logger.From(c)

		var account model.Account
		id := c.Param("id")
		log.Infow("Fetching account balance", "account_id", id)

		if err := db.First(&account, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Warnw("Account not found", "account_id", id)
				c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
				return
			} else {
				log.Errorw("Failed to retrieve balance",
					"account_id", id,
					"error", err,
				)
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve balance", "error": err.Error()})
			return
		}

		log.Infow("Account balance retrieved successfully",
			"account_id", id,
			"balance", account.Balance,
		)

		c.JSON(http.StatusOK, gin.H{"id": id, "balance": account.Balance})
	})
}