package validator

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyField      = errors.New("field is required")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidName     = errors.New("name must be between 2 and 50 characters")
	ErrInvalidPassword = errors.New("password must be at least 6 characters and contain at least one uppercase letter, one lowercase letter, and one number")
)

// ValidateEmail ตรวจสอบรูปแบบของ email
func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmptyField
	}

	// ตรวจสอบรูปแบบ email ด้วย regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

// ValidateName ตรวจสอบความถูกต้องของชื่อ
func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyField
	}

	// ตรวจสอบความยาวของชื่อ
	if len(strings.TrimSpace(name)) < 2 || len(name) > 50 {
		return ErrInvalidName
	}

	return nil
}

// ValidatePassword ตรวจสอบความถูกต้องของรหัสผ่าน
func ValidatePassword(password string) error {
	if password == "" {
		return ErrEmptyField
	}

	// ตรวจสอบความยาวขั้นต่ำ
	if len(password) < 6 {
		return ErrInvalidPassword
	}

	return nil
}

// ValidateUserInput ตรวจสอบข้อมูลผู้ใช้ทั้งหมด
func ValidateUserInput(name, email, password string) error {
	if err := ValidateName(name); err != nil {
		return err
	}

	if err := ValidateEmail(email); err != nil {
		return err
	}

	if err := ValidatePassword(password); err != nil {
		return err
	}

	return nil
}
