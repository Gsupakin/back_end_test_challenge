package domain

import "errors"

var (
	// ข้อผิดพลาดเกี่ยวกับผู้ใช้
	ErrUserNotFound      = errors.New("ไม่พบผู้ใช้ในระบบ")
	ErrUserAlreadyExists = errors.New("มีผู้ใช้นี้ในระบบแล้ว")
	ErrUserInactive      = errors.New("บัญชีผู้ใช้ถูกระงับการใช้งาน")
	ErrInvalidPassword   = errors.New("รหัสผ่านไม่ถูกต้อง")
	ErrInvalidEmail      = errors.New("อีเมลไม่ถูกต้อง")
	ErrInvalidName       = errors.New("ชื่อไม่ถูกต้อง")

	// ข้อผิดพลาดเกี่ยวกับการยืนยันตัวตน
	ErrUnauthorized     = errors.New("กรุณาเข้าสู่ระบบ")
	ErrInvalidToken     = errors.New("โทเค็นไม่ถูกต้องหรือหมดอายุ")
	ErrPermissionDenied = errors.New("ไม่มีสิทธิ์เข้าถึง")

	// ข้อผิดพลาดเกี่ยวกับฐานข้อมูล
	ErrDatabaseConnection = errors.New("ไม่สามารถเชื่อมต่อกับฐานข้อมูลได้")
	ErrDatabaseOperation  = errors.New("เกิดข้อผิดพลาดในการทำงานกับฐานข้อมูล")

	// ข้อผิดพลาดทั่วไป
	ErrInvalidInput       = errors.New("ข้อมูลไม่ถูกต้อง")
	ErrInternalServer     = errors.New("เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์")
	ErrServiceUnavailable = errors.New("บริการไม่พร้อมใช้งาน")
)
