package router

import (
	"database/sql"
	"net/http"
	"payment-service/internal/middlewares/logger"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AccountRouter(r *gin.RouterGroup, db *sqlx.DB) {
	r.GET("/:id/balance", func(c *gin.Context) {
		log := logger.From(c)

		id := c.Param("id")
		log.Infow("Fetching account balance", "account_id", id)

		var balance float64
		err := db.Get(&balance, "SELECT balance FROM accounts WHERE id = $1", id)
		if err != nil {
			if err == sql.ErrNoRows {
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
			"balance", balance,
		)

		c.JSON(http.StatusOK, gin.H{"id": id, "balance": balance})
	})
}