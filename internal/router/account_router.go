package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AccountRouter(r *gin.RouterGroup, db *sqlx.DB) {
	r.GET("/:id/balance", func(c *gin.Context) {
		id := c.Param("id")
		var balance float64
		err := db.Get(&balance, "SELECT balance FROM accounts WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve balance", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "balance": balance})
	})
}