package application

import (
	"net/http"
	"time"

	"github.com/Gsupakin/back_end_test_challeng/internal/domain"
	"github.com/Gsupakin/back_end_test_challeng/pkg/jwt"
	"github.com/Gsupakin/back_end_test_challeng/pkg/utils"
	"github.com/Gsupakin/back_end_test_challeng/pkg/validator"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userRepo domain.UserRepository
	logRepo  domain.LogRepository
}

func NewUserHandler(userRepo domain.UserRepository, logRepo domain.LogRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		logRepo:  logRepo,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	// ตรวจสอบ Content-Type
	if c.GetHeader("Content-Type") != "application/json" {
		h.logRepo.Create(c.Request.Context(), domain.RequestLog{
			Method:    "POST",
			Path:      "/register",
			Status:    http.StatusBadRequest,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Timestamp: time.Now(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
		return
	}

	var user domain.User
	if err := c.BindJSON(&user); err != nil {
		h.logRepo.Create(c.Request.Context(), domain.RequestLog{
			Method:    "POST",
			Path:      "/register",
			Status:    http.StatusBadRequest,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Timestamp: time.Now(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบข้อมูลที่รับเข้ามา
	if err := validator.ValidateUserInput(user.Name, user.Email, user.Password); err != nil {
		h.logRepo.Create(c.Request.Context(), domain.RequestLog{
			Method:    "POST",
			Path:      "/register",
			Status:    http.StatusBadRequest,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Timestamp: time.Now(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบ email ซ้ำ
	_, err := h.userRepo.FindByEmail(c.Request.Context(), user.Email)
	if err == nil {
		h.logRepo.Create(c.Request.Context(), domain.RequestLog{
			Method:    "POST",
			Path:      "/register",
			Status:    http.StatusInternalServerError,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Timestamp: time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists"})
		return
	}

	// ตรวจสอบ name ซ้ำ
	_, err = h.userRepo.FindByName(c.Request.Context(), user.Name)
	if err == nil {
		h.logRepo.Create(c.Request.Context(), domain.RequestLog{
			Method:    "POST",
			Path:      "/register",
			Status:    http.StatusInternalServerError,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Timestamp: time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Name already exists"})
		return
	}

	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		h.logRepo.Create(c.Request.Context(), domain.RequestLog{
			Method:    "POST",
			Path:      "/register",
			Status:    http.StatusInternalServerError,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Timestamp: time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPass
	user.CreatedAt = time.Now()

	id, err := h.userRepo.Create(c.Request.Context(), user)
	if err != nil {
		h.logRepo.Create(c.Request.Context(), domain.RequestLog{
			Method:    "POST",
			Path:      "/register",
			Status:    http.StatusInternalServerError,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Timestamp: time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logRepo.Create(c.Request.Context(), domain.RequestLog{
		Method:    "POST",
		Path:      "/register",
		Status:    http.StatusCreated,
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Timestamp: time.Now(),
	})

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *UserHandler) Login(c *gin.Context) {
	// ตรวจสอบ Content-Type
	if c.GetHeader("Content-Type") != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
		return
	}

	var creds domain.User
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบข้อมูลที่จำเป็น
	if err := validator.ValidateEmail(creds.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if creds.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	user, err := h.userRepo.FindByEmail(c.Request.Context(), creds.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, _ := jwt.GenerateJWT(user.ID.Hex())
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.userRepo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ซ่อน password ของทุก user
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Password = "" // ซ่อน password
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	// ตรวจสอบ Content-Type
	if c.GetHeader("Content-Type") != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
		return
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var updateData struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	}

	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบว่ามีข้อมูลที่จะอัพเดทหรือไม่
	if updateData.Name == "" && updateData.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data to update"})
		return
	}

	update := make(map[string]interface{})
	if updateData.Name != "" {
		// ตรวจสอบ name ซ้ำ
		_, err := h.userRepo.FindByName(c.Request.Context(), updateData.Name)
		if err == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Name already exists"})
			return
		}
		update["name"] = updateData.Name
	}
	if updateData.Email != "" {
		// ตรวจสอบ email ซ้ำ
		_, err := h.userRepo.FindByEmail(c.Request.Context(), updateData.Email)
		if err == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists"})
			return
		}
		update["email"] = updateData.Email
	}

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data to update"})
		return
	}

	err = h.userRepo.Update(c.Request.Context(), objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	err = h.userRepo.Delete(c.Request.Context(), objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
