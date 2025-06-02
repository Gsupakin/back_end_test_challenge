package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Gsupakin/back_end_test_challeng/internal/application"
	"github.com/Gsupakin/back_end_test_challeng/internal/infrastructure"
	"github.com/Gsupakin/back_end_test_challeng/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	client := infrastructure.ConnectMongo()
	db := client.Database("Test")
	userCollection := db.Collection("users")
	logCollection := db.Collection("request_logs")

	// Initialize repositories
	userRepo := infrastructure.NewMongoUserRepository(userCollection)
	logRepo := infrastructure.NewMongoLogRepository(logCollection)

	// Initialize handler
	userHandler := application.NewUserHandler(userRepo, logRepo)

	router := gin.Default()
	router.Use(middleware.RequestLoggerToMongo(logCollection))

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	auth := router.Group("/", middleware.JWTAuth())
	{
		auth.GET("/users", userHandler.ListUsers)
		auth.GET("/users/:id", userHandler.GetUserByID)
		auth.PUT("/users/:id", userHandler.UpdateUser)
		auth.DELETE("/users/:id", userHandler.DeleteUser)
	}

	// Start background goroutine to log user count
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			count, err := userRepo.Count(ctx)
			cancel()

			if err != nil {
				log.Printf("❌ Failed to count users: %v", err)
			} else {
				log.Printf("👥 Total users in DB: %d", count)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	router.Run(":8080")
}
