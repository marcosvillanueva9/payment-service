package main

import (
	"log"
	"net/http"
	"payment-service/config"
	"payment-service/db"
	"payment-service/internal/middleware/auth"
	"payment-service/internal/middleware/logger"
	"payment-service/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	database := db.Connect(cfg.DBURL)

	logger.Init(cfg.APP_ENV)

	r := gin.Default()
	r.Use(logger.Middleware())

	// SOLO PARA PRUEBA
	r.POST("/token/:id", func(c *gin.Context) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id": c.Param("id"),
		})
		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})

	accountGroup := r.Group("/account")
	{
		accountGroup.Use(auth.Middleware(cfg.JWTSecret))
		router.AccountRouter(accountGroup, database)
	}

	transferGroup := r.Group("/transfer")
	{
		transferGroup.Use(auth.Middleware(cfg.JWTSecret))
		router.TransferRouter(transferGroup, database)
	}

	log.Fatal(r.Run(":" + cfg.PORT))
}