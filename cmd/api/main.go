package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gsupakin/back_end_test_challeng/internal/application"
	grpcserver "github.com/Gsupakin/back_end_test_challeng/internal/grpc"
	"github.com/Gsupakin/back_end_test_challeng/internal/infrastructure"
	"github.com/Gsupakin/back_end_test_challeng/middleware"
	pb "github.com/Gsupakin/back_end_test_challeng/proto"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
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

	// สร้าง context ที่สามารถยกเลิกได้
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// เชื่อมต่อ MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

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

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcserver.AuthInterceptor),
	)
	userServer := grpcserver.NewUserServer(userRepo)
	pb.RegisterUserServiceServer(grpcServer, userServer)

	// สร้าง HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// เริ่ม background goroutine สำหรับนับจำนวนผู้ใช้
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping user count goroutine...")
				return
			case <-ticker.C:
				countCtx, countCancel := context.WithTimeout(ctx, 5*time.Second)
				count, err := userRepo.Count(countCtx)
				countCancel()

				if err != nil {
					log.Printf("❌ Failed to count users: %v", err)
				} else {
					log.Printf("👥 Total users in DB: %d", count)
				}
			}
		}
	}()

	// เริ่ม gRPC server ใน goroutine
	go func() {
		log.Println("Starting gRPC server on :50051...")
		grpcListener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// เริ่ม HTTP server ใน goroutine
	go func() {
		log.Println("Starting server on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// รอสัญญาณการปิดโปรแกรม
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// สร้าง context สำหรับการปิด server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// ยกเลิก context หลัก
	cancel()

	// ปิด HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// ปิด gRPC server
	grpcServer.GracefulStop()

	log.Println("Server exited properly")
}
