package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Gsupakin/back_end_test_challeng/internal/application"
	"github.com/Gsupakin/back_end_test_challeng/internal/domain"
	"github.com/Gsupakin/back_end_test_challeng/internal/infrastructure"
	"github.com/Gsupakin/back_end_test_challeng/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	// หา path ของไฟล์ .env โดยอ้างอิงจากตำแหน่งของไฟล์ปัจจุบัน
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(b))) // ขึ้นไป 3 ระดับจากไฟล์ปัจจุบัน
	envPath := filepath.Join(projectRoot, ".env")

	// โหลดค่าจากไฟล์ .env
	if err := godotenv.Load(envPath); err != nil {
		panic("Error loading .env file: " + err.Error())
	}
}

func setupTest() (*gin.Engine, *application.UserHandler, *mongo.Client) {
	gin.SetMode(gin.TestMode)
	client := infrastructure.ConnectMongo()
	db := client.Database("Test")
	userCollection := db.Collection("users")
	logCollection := db.Collection("request_logs")

	userHandler := &application.UserHandler{
		Collection:    userCollection,
		LogCollection: logCollection,
	}

	router := gin.Default()
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	auth := router.Group("/", middleware.JWTAuth())
	{
		auth.GET("/users", userHandler.ListUsers)
		auth.GET("/users/:id", userHandler.GetUserByID)
		auth.PUT("/users/:id", userHandler.UpdateUser)
		auth.DELETE("/users/:id", userHandler.DeleteUser)
	}

	return router, userHandler, client
}

func TestRegister(t *testing.T) {
	router, _, client := setupTest()
	defer client.Disconnect(context.Background())

	t.Run("Register Success", func(t *testing.T) {
		t.Log("Testing successful user registration")
		user := domain.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonData, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %v", response)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NotEmpty(t, response["InsertedID"])
	})

	t.Run("Register Duplicate Email", func(t *testing.T) {
		t.Log("Testing registration with duplicate email")
		// ใช้ email เดียวกับกรณีแรก
		user := domain.User{
			Name:     "Test User 2",
			Email:    "test@example.com", // ใช้ email เดิม
			Password: "password123",
		}
		jsonData, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// ตรวจสอบ response
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %v", response)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, response["error"], "Email already exists")
	})

	t.Run("Register Duplicate Name", func(t *testing.T) {
		t.Log("Testing registration with duplicate name")
		// ใช้ชื่อเดียวกับกรณีแรก
		user := domain.User{
			Name:     "Test User", // ใช้ชื่อเดิม
			Email:    "test2@example.com",
			Password: "password123",
		}
		jsonData, _ := json.Marshal(user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// ตรวจสอบ response
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %v", response)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, response["error"], "Name already exists")
	})
}

func TestLogin(t *testing.T) {
	router, _, client := setupTest()
	defer client.Disconnect(context.Background())

	// Register user ก่อน
	user := domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	t.Run("Login Success", func(t *testing.T) {
		creds := domain.User{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonData, _ := json.Marshal(creds)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)

		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %v", response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, response["token"])
	})

	t.Run("Login Invalid Credentials", func(t *testing.T) {
		creds := domain.User{
			Email:    "test@example.com",
			Password: "123456",
		}
		jsonData, _ := json.Marshal(creds)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)

		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %v", response)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, response["error"], "Invalid email or password")
	})
}

func TestUserOperations(t *testing.T) {
	router, _, client := setupTest()
	defer client.Disconnect(context.Background())

	// Register user ก่อน
	user := domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	insertedID := response["InsertedID"].(string)

	// Login เพื่อรับ token
	creds := domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ = json.Marshal(creds)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	json.Unmarshal(w.Body.Bytes(), &response)
	token := response["token"].(string)

	t.Run("List Users", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		var response []domain.User
		json.Unmarshal(w.Body.Bytes(), &response)

		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %v", response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Greater(t, len(response), 0)
	})

	t.Run("Get User By ID - Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/"+insertedID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		var response domain.User
		json.Unmarshal(w.Body.Bytes(), &response)

		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %v", response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, insertedID, response.ID.Hex())
	})

	t.Run("Get User By ID - Invalid ID Format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/invalid-id", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)

		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %v", response)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, response["error"], "Invalid user ID format")
	})

	t.Run("Get User By ID - Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/507f1f77bcf86cd799439011", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Update User - Success", func(t *testing.T) {
		updateData := map[string]string{
			"name": "Updated Name",
		}
		jsonData, _ := json.Marshal(updateData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/"+insertedID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)

		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %v", response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, response["message"], "User updated successfully")
	})

	t.Run("Update User - Invalid ID Format", func(t *testing.T) {
		updateData := map[string]string{
			"name": "Updated Name",
		}
		jsonData, _ := json.Marshal(updateData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/invalid-id", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update User - No Data", func(t *testing.T) {
		updateData := map[string]string{}
		jsonData, _ := json.Marshal(updateData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/"+insertedID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update User - Duplicate Name", func(t *testing.T) {
		// สร้าง user ใหม่
		newUser := domain.User{
			Name:     "Another User",
			Email:    "another@example.com",
			Password: "password123",
		}
		jsonData, _ := json.Marshal(newUser)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// พยายามอัพเดท user แรกให้มีชื่อซ้ำกับ user ใหม่
		updateData := map[string]string{
			"name": "Another User",
		}
		jsonData, _ = json.Marshal(updateData)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("PUT", "/users/"+insertedID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Delete User - Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/"+insertedID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)

		t.Logf("Response status: %d", w.Code)
		t.Logf("Response body: %v", response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, response["message"], "User deleted successfully")
	})

	t.Run("Delete User - Invalid ID Format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/invalid-id", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Delete User - Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/507f1f77bcf86cd799439011", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
