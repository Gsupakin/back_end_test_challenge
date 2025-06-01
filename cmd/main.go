package main

import (
	"github.com/Gsupakin/back_end_test_challeng/internal/application"
	"github.com/Gsupakin/back_end_test_challeng/internal/infrastructure"
	"github.com/Gsupakin/back_end_test_challeng/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	client := infrastructure.ConnectMongo()
	db := client.Database("yourdb")
	userCollection := db.Collection("users")
	logCollection := db.Collection("request_logs") // ✅ ใช้ log collection

	userHandler := application.UserHandler{
		Collection: userCollection,
	}

	router := gin.Default()
	router.Use(middleware.RequestLoggerToMongo(logCollection)) // ✅ ใช้งาน Middleware

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	auth := router.Group("/", middleware.JWTAuth())
	{
		auth.GET("/users", userHandler.ListUsers)
		auth.GET("/users/:id", userHandler.GetUserByID)
		auth.PUT("/users/:id", userHandler.UpdateUser)
		auth.DELETE("/users/:id", userHandler.DeleteUser)
	}
	router.Run(":8080")
}
