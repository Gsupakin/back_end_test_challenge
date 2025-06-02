package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User แทนข้อมูลผู้ใช้ในระบบ
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password" validate:"required"` // ไม่แสดงใน JSON
	Role      string             `json:"role" bson:"role"`                             // เพิ่ม role
	Status    string             `json:"status" bson:"status"`                         // เพิ่ม status
	LastLogin *time.Time         `json:"last_login,omitempty" bson:"last_login,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

// NewUser สร้างผู้ใช้ใหม่
func NewUser(name, email, password string) *User {
	now := time.Now()
	return &User{
		Name:      name,
		Email:     email,
		Password:  password,
		Role:      "user",   // ค่าเริ่มต้น
		Status:    "active", // ค่าเริ่มต้น
		CreatedAt: now,
		UpdatedAt: &now,
	}
}

// IsActive ตรวจสอบว่าผู้ใช้ยังใช้งานอยู่หรือไม่
func (u *User) IsActive() bool {
	return u.Status == "active" && u.DeletedAt == nil
}

// IsAdmin ตรวจสอบว่าเป็นผู้ดูแลระบบหรือไม่
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// UpdateLastLogin อัพเดทเวลาล็อกอินล่าสุด
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = &now
}

// SoftDelete ทำการลบแบบ soft delete
func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.Status = "inactive"
	u.UpdatedAt = &now
}
