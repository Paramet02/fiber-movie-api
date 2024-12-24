package repository

import (
	"gorm.io/gorm"
	"time"
)

// Adapter: โครงสร้างที่ทำหน้าที่ implement ฟังก์ชันใน Interface (Port) โดยเชื่อมต่อกับฐานข้อมูล
type userRepositoryDB  struct {
	db *gorm.DB
}

// ฟังก์ชันสำหรับสร้าง MovieRepository ด้วยการรวม Port (Interface) และ Adapter (Implementation)
func NewuserRepositoryDB(db *gorm.DB) UserRepository {
	return userRepositoryDB{db}
}

func (r userRepositoryDB) GetUserByEmail(Email string) (user *user , err error) {
	if result := r.db.Where("email = ?", Email).First(&user); result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r userRepositoryDB) GetUserByID(ID int) (user *user , err error) {
	if result := r.db.Where("ID = ?" , ID).First(&user) ; result.Error != nil {
		return nil , result.Error
	}
	return user , nil
}

// Create: ฟังก์ชันสำหรับสร้างผู้ใช้ใหม่ในฐานข้อมูล
func (r userRepositoryDB) Create(user *user) (*user, error) {
	// กำหนดเวลาที่สร้างและเวลาที่อัพเดท
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// สร้างผู้ใช้ในฐานข้อมูล
	if result := r.db.Create(&user); result.Error != nil {
		return nil, result.Error
	}

	// ส่งกลับผู้ใช้ที่ถูกสร้างแล้ว
	return user, nil
}
