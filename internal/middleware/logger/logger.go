package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const loggerKey = "logger"

var baseLogger *zap.SugaredLogger

func Init(env string) {
	var logger *zap.Logger
	var err error
	if env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	baseLogger = logger.Sugar()
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New().String()
		logger := baseLogger.With("request_id", requestId)

		c.Set(loggerKey, logger)
		c.Next()
	}
}

func From(c *gin.Context) *zap.SugaredLogger {
	logger, exists := c.Get(loggerKey)
	if !exists {
		return baseLogger
	}

	if sugaredLogger, ok := logger.(*zap.SugaredLogger); ok {
		return sugaredLogger
	}

	return baseLogger
}