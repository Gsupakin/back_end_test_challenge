package validator

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyField      = errors.New("กรุณากรอกข้อมูล")
	ErrInvalidEmail    = errors.New("รูปแบบอีเมลไม่ถูกต้อง")
	ErrInvalidName     = errors.New("ชื่อต้องมีความยาว 2-50 ตัวอักษร")
	ErrInvalidPassword = errors.New("รหัสผ่านต้องมีความยาวอย่างน้อย 6 ตัวอักษร และต้องมีตัวพิมพ์ใหญ่ ตัวพิมพ์เล็ก และตัวเลขอย่างน้อย 1 ตัว")
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

	// ตรวจสอบความยาวของอีเมล
	if len(email) > 100 {
		return errors.New("อีเมลยาวเกินไป")
	}

	return nil
}

// ValidateName ตรวจสอบความถูกต้องของชื่อ
func ValidateName(name string) error {
	if name == "" {
		return ErrEmptyField
	}

	// ตรวจสอบความยาวของชื่อ
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 50 {
		return ErrInvalidName
	}

	// ตรวจสอบว่ามีตัวอักษรพิเศษหรือไม่
	if !regexp.MustCompile(`^[a-zA-Z0-9ก-๙\s]+$`).MatchString(name) {
		return errors.New("ชื่อต้องประกอบด้วยตัวอักษร ตัวเลข และช่องว่างเท่านั้น")
	}

	return nil
}

// ValidatePassword ตรวจสอบความถูกต้องของรหัสผ่าน
func ValidatePassword(password string) error {
	if password == "" {
		return ErrEmptyField
	}

	// ตรวจสอบความยาวของรหัสผ่าน
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
