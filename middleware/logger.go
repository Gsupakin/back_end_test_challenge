package middleware

import (
	"log"
	"time"

	"github.com/Gsupakin/back_end_test_challeng/internal/domain"

	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RequestLoggerToMongo(logCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logEntry := domain.RequestLog{
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			LatencyMS: duration.Milliseconds(),
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Timestamp: time.Now(),
		}

		go func(entry domain.RequestLog) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := logCollection.InsertOne(ctx, entry)
			if err != nil {
				log.Printf("Failed to log request: %v", err)
			}
		}(logEntry)
	}
}
