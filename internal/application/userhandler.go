package application

import (
	"context"
	"net/http"
	"time"

	"github.com/Gsupakin/back_end_test_challeng/internal/domain"
	"github.com/Gsupakin/back_end_test_challeng/pkg/jwt"
	"github.com/Gsupakin/back_end_test_challeng/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	Collection    *mongo.Collection
	LogCollection *mongo.Collection
}

func (h *UserHandler) Register(c *gin.Context) {
	// ตรวจสอบ Content-Type
	if c.GetHeader("Content-Type") != "application/json" {
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusBadRequest,
			"error":      "Content-Type must be application/json",
			"created_at": time.Now(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
		return
	}

	var user domain.User
	if err := c.BindJSON(&user); err != nil {
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusBadRequest,
			"error":      err.Error(),
			"created_at": time.Now(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบข้อมูลที่จำเป็น
	if user.Email == "" || user.Name == "" || user.Password == "" {
		errorMsg := "Missing required fields"
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusBadRequest,
			"error":      errorMsg,
			"email":      user.Email,
			"name":       user.Name,
			"created_at": time.Now(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	// ตรวจสอบรูปแบบ email
	if !utils.IsValidEmail(user.Email) {
		errorMsg := "Invalid email format"
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusBadRequest,
			"error":      errorMsg,
			"email":      user.Email,
			"created_at": time.Now(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	// ตรวจสอบความยาวของ password
	if len(user.Password) < 6 {
		errorMsg := "Password must be at least 6 characters"
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusBadRequest,
			"error":      errorMsg,
			"email":      user.Email,
			"created_at": time.Now(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	// ตรวจสอบ email ซ้ำ
	var existingUser domain.User
	err := h.Collection.FindOne(context.Background(), bson.M{
		"email": user.Email,
	}).Decode(&existingUser)

	if err == nil {
		// กรณีพบ email ซ้ำ
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusInternalServerError,
			"error":      "Email already exists",
			"email":      user.Email,
			"created_at": time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		// กรณีเกิด error อื่นๆ
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusInternalServerError,
			"error":      err.Error(),
			"email":      user.Email,
			"created_at": time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// ตรวจสอบ name ซ้ำ
	err = h.Collection.FindOne(context.Background(), bson.M{
		"name": user.Name,
	}).Decode(&existingUser)

	if err == nil {
		// กรณีพบ name ซ้ำ
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusInternalServerError,
			"error":      "Name already exists",
			"name":       user.Name,
			"created_at": time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Name already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		// กรณีเกิด error อื่นๆ
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusInternalServerError,
			"error":      err.Error(),
			"name":       user.Name,
			"created_at": time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusInternalServerError,
			"error":      "Failed to hash password",
			"email":      user.Email,
			"created_at": time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPass
	user.CreatedAt = time.Now()

	res, err := h.Collection.InsertOne(context.Background(), user)
	if err != nil {
		// บันทึก log error
		h.LogCollection.InsertOne(context.Background(), bson.M{
			"endpoint":   "/register",
			"method":     "POST",
			"status":     http.StatusInternalServerError,
			"error":      err.Error(),
			"email":      user.Email,
			"name":       user.Name,
			"created_at": time.Now(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// บันทึก log สำเร็จ
	h.LogCollection.InsertOne(context.Background(), bson.M{
		"endpoint":   "/register",
		"method":     "POST",
		"status":     http.StatusCreated,
		"user_id":    res.InsertedID,
		"email":      user.Email,
		"name":       user.Name,
		"created_at": time.Now(),
	})

	c.JSON(http.StatusCreated, res)
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
	if creds.Email == "" || creds.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	var user domain.User
	err := h.Collection.FindOne(context.Background(), bson.M{
		"email":      creds.Email,
		"deleted_at": nil,
	}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// เพิ่ม log เพื่อตรวจสอบ password
	h.LogCollection.InsertOne(context.Background(), bson.M{
		"endpoint":        "/login",
		"method":          "POST",
		"status":          http.StatusOK,
		"email":           creds.Email,
		"hashed_password": user.Password,
		"input_password":  creds.Password,
		"password_match":  utils.CheckPasswordHash(creds.Password, user.Password),
		"created_at":      time.Now(),
	})

	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, _ := jwt.GenerateJWT(user.ID.Hex())
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ตัวอย่าง ListUsers:
func (h *UserHandler) ListUsers(c *gin.Context) {
	cursor, _ := h.Collection.Find(context.Background(), bson.M{"deleted_at": nil})
	var users []domain.User
	cursor.All(context.Background(), &users)

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

	var user domain.User
	err = h.Collection.FindOne(context.Background(), bson.M{
		"_id":        objID,
		"deleted_at": nil,
	}).Decode(&user)
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

	update := bson.M{}
	if updateData.Name != "" {
		// ตรวจสอบ name ซ้ำ
		var existingUser domain.User
		err := h.Collection.FindOne(context.Background(), bson.M{
			"name": updateData.Name,
			"_id": bson.M{
				"$ne": objID,
			},
		}).Decode(&existingUser)
		if err == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Name already exists"})
			return
		}
		update["name"] = updateData.Name
	}
	if updateData.Email != "" {
		// ตรวจสอบ email ซ้ำ
		var existingUser domain.User
		err := h.Collection.FindOne(context.Background(), bson.M{
			"email": updateData.Email,
			"_id": bson.M{
				"$ne": objID,
			},
		}).Decode(&existingUser)
		if err == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists"})
			return
		}
		update["email"] = updateData.Email
	}
	update["updated_at"] = time.Now()

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data to update"})
		return
	}

	result, err := h.Collection.UpdateOne(
		context.Background(),
		bson.M{
			"_id":        objID,
			"deleted_at": nil,
		},
		bson.M{"$set": update},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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

	result, err := h.Collection.UpdateOne(
		context.Background(),
		bson.M{
			"_id":        objID,
			"deleted_at": nil,
		},
		bson.M{"$set": bson.M{"deleted_at": time.Now()}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
