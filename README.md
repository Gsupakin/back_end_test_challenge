# Backend Test Challenge

## การติดตั้งและรันโปรเจค

### Prerequisites
- Go 1.16 หรือสูงกว่า
- MongoDB
- Docker (optional)

### การติดตั้ง
1. Clone repository:
```bash
git clone <repository-url>
cd back_end_test_challeng
```

2. ติดตั้ง dependencies:
```bash
go mod download
```

3. สร้างไฟล์ .env ในโฟลเดอร์หลักของโปรเจค:
```env
MONGODB_URI=mongodb://localhost:27017
DB_NAME=your_database_name
JWT_SECRET=your_jwt_secret_key
```

4. รันแอพพลิเคชัน:
```bash
go run main.go
```

หรือใช้ Docker:
```bash
docker-compose up
```

## การใช้งาน JWT Token

### การรับ Token
1. สมัครสมาชิกผ่าน `/register` endpoint
2. Login ผ่าน `/login` endpoint เพื่อรับ token
3. ใช้ token ที่ได้ในการเรียกใช้ API ที่ต้องการการยืนยันตัวตน

### การใช้ Token
เพิ่ม header ในทุก request ที่ต้องการการยืนยันตัวตน:
```
Authorization: Bearer <your_token>
```

## API Endpoints และตัวอย่างการใช้งาน

### 1. Register User
```http
POST /register
Content-Type: application/json

{
    "name": "Test User",
    "email": "test@example.com",
    "password": "Password123"
}
```

Response (201 Created):
```json
{
    "InsertedID": "user_id_here"
}
```

### 2. Login
```http
POST /login
Content-Type: application/json

{
    "email": "test@example.com",
    "password": "Password123"
}
```

Response (200 OK):
```json
{
    "token": "jwt_token_here"
}
```

### 3. Get All Users (ต้องมี JWT Token)
```http
GET /users
Authorization: Bearer <your_token>
```

Response (200 OK):
```json
[
    {
        "id": "user_id",
        "name": "Test User",
        "email": "test@example.com"
    }
]
```

### 4. Get User by ID (ต้องมี JWT Token)
```http
GET /users/:id
Authorization: Bearer <your_token>
```

Response (200 OK):
```json
{
    "id": "user_id",
    "name": "Test User",
    "email": "test@example.com"
}
```

### 5. Update User (ต้องมี JWT Token)
```http
PUT /users/:id
Authorization: Bearer <your_token>
Content-Type: application/json

{
    "name": "Updated Name"
}
```

Response (200 OK):
```json
{
    "message": "User updated successfully"
}
```

### 6. Delete User (ต้องมี JWT Token)
```http
DELETE /users/:id
Authorization: Bearer <your_token>
```

Response (200 OK):
```json
{
    "message": "User deleted successfully"
}
```

## การตัดสินใจและการออกแบบ

### 1. โครงสร้างโปรเจค
- ใช้ Clean Architecture แบ่งเป็น layers: domain, application, infrastructure
- แยก business logic ออกจาก infrastructure concerns
- ใช้ dependency injection เพื่อลดการ coupling

### 2. การรักษาความปลอดภัย
- ใช้ JWT สำหรับการยืนยันตัวตน
- เข้ารหัสรหัสผ่านก่อนเก็บในฐานข้อมูล
- ตรวจสอบความถูกต้องของข้อมูลที่รับเข้ามา

### 3. การจัดการฐานข้อมูล
- ใช้ MongoDB เป็นฐานข้อมูลหลัก
- สร้าง indexes สำหรับ fields ที่ใช้ในการค้นหาบ่อย
- ใช้ transactions เมื่อจำเป็น

### 4. การจัดการ Error
- สร้าง custom error types
- ใช้ middleware สำหรับการจัดการ error
- ส่ง error messages ที่เป็นประโยชน์กลับไปให้ client

### 5. การ Logging
- บันทึก request logs ลง MongoDB
- ใช้ structured logging
- บันทึกข้อมูลสำคัญสำหรับการ debug

## การทดสอบ
รัน unit tests:
```bash
go test ./...
```

รัน integration tests:
```bash
go test ./tests/...
``` 