# Backend Test Challenge

สวัสดีครับ! นี่คือโปรเจค Backend Test Challenge ที่พัฒนาด้วย Go และ MongoDB ครับ

## วิธีการติดตั้งและรันโปรเจค

### สิ่งที่ต้องมีก่อนเริ่มต้น
- Go เวอร์ชัน 1.16 หรือใหม่กว่า
- MongoDB
- Docker (ถ้าต้องการรันผ่าน Docker)

### ขั้นตอนการติดตั้ง
1. ดาวน์โหลดโค้ดจาก repository:
```bash
git clone <repository-url>
cd back_end_test_challeng
```

2. ติดตั้ง packages ที่จำเป็น:
```bash
go mod download
```

3. สร้างไฟล์ .env ในโฟลเดอร์หลักของโปรเจค โดยใส่ค่าตามนี้:
```env
MONGODB_URI=mongodb://localhost:27017
DB_NAME=your_database_name
JWT_SECRET=your_jwt_secret_key
```

4. รันแอพพลิเคชัน:
```bash
go run main.go
```

หรือถ้าต้องการรันผ่าน Docker:
```bash
docker-compose up
```

## วิธีการใช้งาน JWT Token

### วิธีการรับ Token
1. สมัครสมาชิกผ่าน endpoint `/register`
2. เข้าสู่ระบบผ่าน endpoint `/login` เพื่อรับ token
3. นำ token ที่ได้ไปใช้ในการเรียก API ที่ต้องการการยืนยันตัวตน

### วิธีการใช้ Token
เพิ่ม header นี้ในทุก request ที่ต้องการการยืนยันตัวตน:
```
Authorization: Bearer <your_token>
```

## API Endpoints และตัวอย่างการใช้งาน

### 1. สมัครสมาชิก
```http
POST /register
Content-Type: application/json

{
    "name": "Test User",
    "email": "test@example.com",
    "password": "Password123"
}
```

คำตอบที่ได้ (201 Created):
```json
{
    "InsertedID": "user_id_here"
}
```

### 2. เข้าสู่ระบบ
```http
POST /login
Content-Type: application/json

{
    "email": "test@example.com",
    "password": "Password123"
}
```

คำตอบที่ได้ (200 OK):
```json
{
    "token": "jwt_token_here"
}
```

### 3. ดึงข้อมูลผู้ใช้ทั้งหมด (ต้องมี JWT Token)
```http
GET /users
Authorization: Bearer <your_token>
```

คำตอบที่ได้ (200 OK):
```json
[
    {
        "id": "user_id",
        "name": "Test User",
        "email": "test@example.com"
    }
]
```

### 4. ดึงข้อมูลผู้ใช้ตาม ID (ต้องมี JWT Token)
```http
GET /users/:id
Authorization: Bearer <your_token>
```

คำตอบที่ได้ (200 OK):
```json
{
    "id": "user_id",
    "name": "Test User",
    "email": "test@example.com"
}
```

### 5. อัพเดทข้อมูลผู้ใช้ (ต้องมี JWT Token)
```http
PUT /users/:id
Authorization: Bearer <your_token>
Content-Type: application/json

{
    "name": "Updated Name"
}
```

คำตอบที่ได้ (200 OK):
```json
{
    "message": "User updated successfully"
}
```

### 6. ลบผู้ใช้ (ต้องมี JWT Token)
```http
DELETE /users/:id
Authorization: Bearer <your_token>
```

คำตอบที่ได้ (200 OK):
```json
{
    "message": "User deleted successfully"
}
```

## การออกแบบ

### 1. โครงสร้างโปรเจค
- ใช้ Clean Architecture แบ่งเป็น 3 ชั้น: domain, application, infrastructure
- แยก business logic ออกจากส่วนที่เกี่ยวกับ infrastructure
- ใช้ dependency injection เพื่อให้โค้ดยืดหยุ่นและทดสอบง่าย

### 2. ความปลอดภัย
- ใช้ JWT สำหรับการยืนยันตัวตน
- เข้ารหัสรหัสผ่านก่อนเก็บในฐานข้อมูล
- มีการตรวจสอบความถูกต้องของข้อมูลที่รับเข้ามา

### 3. การจัดการฐานข้อมูล
- เลือกใช้ MongoDB เป็นฐานข้อมูลหลัก
- สร้าง indexes สำหรับ fields ที่ใช้ค้นหาบ่อยๆ
- ใช้ transactions เมื่อต้องทำการอัพเดทข้อมูลหลายส่วน

### 4. การจัดการข้อผิดพลาด
- สร้าง custom error types เพื่อจัดการข้อผิดพลาดได้ดีขึ้น
- ใช้ middleware จัดการข้อผิดพลาดแบบรวมที่เดียว
- ส่งข้อความแจ้งเตือนที่เป็นปัญหากลับไปให้ผู้ใช้

### 5. การบันทึก Log
- บันทึก request logs ลง MongoDB
- ใช้ structured logging เพื่อให้ค้นหาและวิเคราะห์ได้ง่าย
- บันทึกข้อมูลสำคัญสำหรับการแก้ไขปัญหา

## การทดสอบ
รัน unit tests ทั้งหมด:
```bash
go test ./...
```

รัน integration tests:
```bash
go test ./tests/...
```